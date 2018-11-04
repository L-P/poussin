package cpu

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/L-P/poussin/emu/ppu"
	"github.com/L-P/poussin/emu/rom"
)

// CPU is a Sharp LR35902 emulator.
type CPU struct {
	Registers

	// PPU is our pixel processing unit, the CPU access it for VRAM access,
	// VBlank interrupts, some I/O registers, and more.
	PPU *ppu.PPU

	// Mem holds the unmapped memory for the whole addressable space, thus it
	// lacks VRAM, ROM0/ROMX, mirrored memory, etc.
	Mem [0xFFFF]byte

	// Boot holds the bootstrap ROM mapped to 0x0000-0x00FF on DMG
	Boot [256]byte

	// ROM holds the ROM mapped to ROM0/ROMX
	ROM [1024 * 1024 * 8]byte // Per Wikipedia, a GB ROM is 8 MiB max

	// Halted is set by the HALT instruction, it can only be reset by interrupts.
	Halted bool

	// {{{ Debug
	// EnableDebug is the master flag for enabling debug values.
	EnableDebug bool

	// InCycle is true when the CPU is executing to allow tracking memory use
	// without polluting the buffers when the debugger reads memory outside of
	// a cycle.
	InCycle bool

	// MemIOBuffer contains all memory reads/writes (except for instruction
	// fetching) since it was last cleared. It holds six bytes per fetch/write,
	// two for the current PC, one for read (0x01) or write (0x02), two for the
	// address of the read (little-endian), and one for the value.
	MemIOBuffer bytes.Buffer

	// LastOpcode is the last opcode executed by the CPU
	LastOpcode byte

	// LastCycleWasInterrupt is set to true right after an interrupt was triggered
	LastCycleWasInterrupt bool

	// LastOpcodeWasCB indicates if LastOpcode was from the CB opcode table
	LastOpcodeWasCB bool

	// LastLowArg contains the LSB of the 16b argument given to the last
	// executed opcode if applicable.
	LastLowArg byte

	// LastHighArg contains the MSB of the 16b argument given to the last
	// executed opcode if applicable.
	LastHighArg byte

	// SBBuffer contains the bytes written to the serial I/O, it's the debugger
	// responsibility to clear this buffer (not concurrently with Execute of course).
	SBBuffer bytes.Buffer // bytes written to IOSB

	// Jumped indicates if the last conditional call or return did jump
	Jumped bool
	// }}} Debug

	// InterruptEnable contains the IE flag, as we only keep 0xFFFF worth of
	// memory we can't address it otherwise.
	InterruptEnable byte // Addressable at 0xFFFF

	// InterruptMaster is IME, the global hidden flag set by EI and DI.
	InterruptMaster bool

	// Cycle holds the current clock cycle number, it is never reset, which may
	// bite me later, I don't know. Why would you let the emulator run for a
	// month anyway?
	Cycle int

	// LastTimerUpdateCycle is set with the current Cycle when we update the
	// timer registers (DIV, TMA, TIMA)
	LastTimerUpdateCycle int

	// On DMG DIV is implemented as a 16b value that increases on every cycle,
	// the actual DIV register value is the upper 8 bits of this internal
	// register.
	InternalDIV uint16

	Joypad      JoypadState
	JoypadInput <-chan JoypadState
}

// New creates a new CPU object. It is safe to pass a nil input channel if you
// don't want to send inputs.
func New(ppu *ppu.PPU, input <-chan JoypadState, debug bool) CPU {
	c := CPU{
		PPU:         ppu,
		EnableDebug: debug,
		JoypadInput: input,
	}

	c.Reset()

	return c
}

// Reset resets the CPU internal state.
func (c *CPU) Reset() {
	c.InternalDIV = 0
	c.WriteIF(0)
	c.WriteIE(0)
	c.WriteTAC(0)
}

// SimulateBoot puts the CPU in the same state it would be after running the Nintendo boot ROM.
func (c *CPU) SimulateBoot() {
	c.A = 0x01
	c.SetF(0xB0)
	c.BC = 0x0013
	c.DE = 0x00D8
	c.HL = 0x014D
	c.PC = 0x0100
	c.SP = 0xFFFE

	c.WriteIF(0)
	c.WriteIE(0)
	c.WriteTAC(0)
	c.InterruptMaster = false
	c.WriteIO(IODisableBootROM, 0x01)
}

// Step runs the next CPU instruction.
func (c *CPU) Step() (int, error) {
	c.updateJoypad()
	c.Jumped = false
	defer c.UpdateTimers()

	c.LastCycleWasInterrupt = false
	if cycles := c.CheckInterrupts(); cycles > 0 {
		c.LastCycleWasInterrupt = true
		c.InterruptMaster = false
		c.Cycle += cycles
		c.Halted = false
		return cycles, nil
	}

	if c.Halted {
		c.Cycle += 4
		return 4, nil
	}

	opcode := c.Fetch(c.PC)
	cb := opcode == 0xCB
	if cb {
		c.PC++
		opcode = c.Fetch(c.PC)
	}
	c.LastOpcodeWasCB = cb

	ins, err := c.Decode(opcode, cb)
	c.Cycle += 4
	if err != nil {
		return 0, err
	}

	var l, h byte
	if ins.Length > 1 {
		l = c.Fetch(c.PC + 1)
	}
	if ins.Length > 2 {
		h = c.Fetch(c.PC + 2)
	}

	c.LastOpcode = opcode
	c.LastLowArg = l
	c.LastHighArg = h

	return c.Execute(ins, l, h)
}

// Decode decodes an opcode into an Instruction.
func (c *CPU) Decode(opcode byte, cb bool) (Instruction, error) {
	ins := Decode(opcode, cb)
	if !ins.Valid() {
		return Instruction{}, fmt.Errorf(
			"opcode not found: 0x%02X (CB: %t)",
			opcode,
			cb,
		)
	}

	return ins, nil
}

// Execute runs a single instruction and returns the number of cycles it took to run.
func (c *CPU) Execute(ins Instruction, l, h byte) (int, error) {
	c.InCycle = true
	defer func() { c.InCycle = false }()
	if ins.Func == nil {
		return 0, fmt.Errorf("no function defined for %s", ins.Name)
	}

	c.PC += uint16(ins.Length)
	c.Cycle += int(ins.Cycles)
	ins.Func(c, l, h)

	return int(ins.Cycles), nil
}

// LoadBootROM puts a boot rom in the 256 first bytes or RAM.
func (c *CPU) LoadBootROM(data []byte) error {
	if count := copy(c.Boot[:], data); count != 256 {
		return fmt.Errorf("did not copy 256 bytes: %d", count)
	}

	return nil
}

// LoadROM loads a ROM in RAM.
func (c *CPU) LoadROM(data []byte) error {
	copy(c.ROM[:], data)

	h := rom.NewHeader(data)
	if h.CGBOnly {
		return errors.New("only DMG games are supported")
	}

	return nil
}

// DoInterrupt jumps to the given interrupt handler and pushes the previous
// program counter to the stack.
func (c *CPU) DoInterrupt(addr uint16) int {
	c.InCycle = true
	c.StackPush16b(c.PC)
	c.PC = addr
	c.InCycle = false
	return 20
}

// CheckInterrupts checks if any interrupt should run this step instead of
// continuing execution.
func (c *CPU) CheckInterrupts() int {
	// When HALTed, interrupts are active even when IME is off.
	// IF seems to be reset only if IME is on though.
	if !c.InterruptMaster && !c.Halted {
		return 0
	}

	if c.PPU.InterruptVBlank {
		c.SetIF(IEVBlank)
		c.PPU.InterruptVBlank = false
	}

	if c.IFIsSet(IETimer) && c.IEEnabled(IETimer) {
		c.Mem[IOTIMA] = c.Mem[IOTMA]
		if c.InterruptMaster {
			c.UnSetIF(IETimer)
		}

		return c.DoInterrupt(0x0050)
	}

	if c.IFIsSet(IEVBlank) && c.IEEnabled(IEVBlank) {
		if c.InterruptMaster {
			c.UnSetIF(IEVBlank)
		}

		c.UnSetIF(IEVBlank)
		return c.DoInterrupt(0x0040)
	}

	if c.IFIsSet(IEJoypad) {
		c.UnSetIF(IEJoypad)
		return c.DoInterrupt(0x0060)
	}

	return 0
}

// UpdateTimers updates the TIMA and DIV registers based on cycle count.
func (c *CPU) UpdateTimers() {
	delta := c.Cycle - c.LastTimerUpdateCycle
	if delta <= 0 {
		return
	}
	c.LastTimerUpdateCycle = c.Cycle

	oldTIMA := c.Mem[IOTIMA]
	oldIDIV := c.InternalDIV
	c.InternalDIV += uint16(delta)

	var valMask uint16
	var owMask uint16
	if c.IsTACEnabled() {
		switch c.FetchTAC() & TACSpeedMask {
		case TACSpeed262:
			valMask = 0x0F
			owMask = 0x10
		case TACSpeed65:
			valMask = 0x3F
			owMask = 0x40
		case TACSpeed16:
			valMask = 0xFF
			owMask = 0x0100
		case TACSpeed4:
			valMask = 0x03FF
			owMask = 0x0400
		}

		if (((oldIDIV & valMask) + (uint16(delta) & valMask)) & owMask) == owMask {
			c.Mem[IOTIMA]++
		}
	}

	if oldTIMA > c.Mem[IOTIMA] {
		c.SetIF(IETimer)
	}
}
