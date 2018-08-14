package cpu

import (
	"bytes"
	"errors"
	"fmt"

	"home.leo-peltier.fr/poussin/emu/ppu"
	"home.leo-peltier.fr/poussin/emu/rom"
)

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

	// Stopped is set by the STOP instruction, it can only be reset by interrupts.
	Stopped bool

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

	// LastOpcodeWasCB indicates if LastOpcode was from the CB opcode table
	LastOpcodeWasCB bool

	// LastLowArg contains the LSB of the 16b argument given to the last
	// executed opcode if applicable.
	LastLowArg byte

	// LastHighArg contains the MSB of the 16b argument given to the last
	// executed opcode if applicable.
	LastHighArg byte

	// SBBuffer contains the bytes written to the serial I/O, it's the debugger
	// responsibility to clear this buffer (not concurenttly with Execute of course).
	SBBuffer bytes.Buffer // bytes written to IOSB
	// }}} Debug

	// InterruptEnable contains the IE flag, as we only keep 0xFFFF worth of
	// memory we can't address it otherwise.
	InterruptEnable byte // Adressable at 0xFFFF

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
}

func New(ppu *ppu.PPU) CPU {
	c := CPU{
		PPU:         ppu,
		EnableDebug: true,
	}

	c.Reset()

	return c
}

func (c *CPU) Reset() {
	c.InternalDIV = 0
	c.WriteIF(0)
	c.WriteTAC(0)
}

func (c *CPU) Step() (int, error) {
	defer c.UpdateTimers()

	if cycles := c.CheckInterrupts(); cycles > 0 {
		c.InterruptMaster = false
		c.Cycle += cycles
		c.Stopped = false
		return cycles, nil
	}

	if c.Stopped {
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

func (c *CPU) LoadBootROM(data []byte) error {
	if count := copy(c.Boot[:], data); count != 256 {
		return fmt.Errorf("did not copy 256 bytes: %d", count)
	}

	return nil
}

func (c *CPU) LoadROM(data []byte) error {
	copy(c.ROM[:], data)

	h := rom.NewHeader(data)
	if h.CGBOnly {
		return errors.New("only DMG games are supported")
	}

	return nil
}

func (c *CPU) DoInterrupt(addr uint16) int {
	c.InCycle = true
	c.StackPush16b(c.PC)
	c.PC = addr
	c.InCycle = false
	return 20
}

func (c *CPU) CheckInterrupts() int {
	if !c.InterruptMaster {
		return 0
	}

	if c.IFIsSet(IETimer) && c.IEEnabled(IETimer) {
		c.UnSetIF(IETimer)
		c.Mem[IOTIMA] = c.Mem[IOTMA]
		return c.DoInterrupt(0x0050)
	}

	if c.PPU.InterruptVBlank && c.IEEnabled(IEVBlank) {
		c.PPU.InterruptVBlank = false
		return c.DoInterrupt(0x0040)
	}

	return 0
}

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
