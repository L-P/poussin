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

	InterruptEnable byte // Adressable at 0xFFFF
	InterruptTimer  bool

	// InterruptMaster is IME, the global hidden flag set by EI and DI.
	InterruptMaster bool

	// NextOpcodeIsCB is true when the next opcode has to be decoded from the
	// CB opcode table, it's set by the PREFIX CB instruction.
	// from the CB opcodes
	NextOpcodeIsCB bool

	// Cycle holds the current clock cycle number, it is never reset, which may
	// bite me later, I don't know. Why would you let the emulator run for a
	// month anyway?
	Cycle int

	// LastTimerUpdateCycle is set with the current Cycle when we update the
	// timer registers (DIV, TMA, TIMA, TAC)
	LastTimerUpdateCycle int

	// TimerOverflow is true when TIMA overflowed on the previous cycle
	TimerOverflow bool
}

func New(ppu *ppu.PPU) CPU {
	return CPU{
		PPU:         ppu,
		EnableDebug: true,
	}
}

func (c *CPU) Step() (int, error) {
	if cycles := c.CheckInterrupts(); cycles > 0 {
		c.InterruptMaster = false
		c.Cycle += cycles
		c.Stopped = false
		c.UpdateTimer()
		return cycles, nil
	}

	if c.Stopped {
		c.Cycle += 4
		c.UpdateTimer()
		return 4, nil
	}

	c.UpdateTimer()

	opcode := c.Fetch(c.PC)
	cb := c.NextOpcodeIsCB
	ins, err := c.Decode(opcode)
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

	if c.EnableDebug {
		c.LastOpcode = opcode
		c.LastOpcodeWasCB = cb // CB is reset after Decode, hence cb
		c.LastLowArg = l
		c.LastHighArg = h
	}

	return c.Execute(ins, l, h)
}

func (c *CPU) Decode(opcode byte) (Instruction, error) {
	ins := Decode(opcode, c.NextOpcodeIsCB)
	if !ins.Valid() {
		return Instruction{}, fmt.Errorf(
			"opcode not found: 0x%02X (CB: %t)",
			opcode,
			c.NextOpcodeIsCB,
		)
	}
	c.NextOpcodeIsCB = false

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
	c.StackPush16b(c.PC)
	c.PC = addr
	return 20
}

func (c *CPU) CheckInterrupts() int {
	if !c.InterruptMaster {
		return 0
	}

	if c.InterruptTimer && c.IEEnabled(IETimer) {
		c.InterruptTimer = false
		return c.DoInterrupt(0x0050)
	}

	if c.PPU.InterruptVBlank && c.IEEnabled(IEVBlank) {
		c.PPU.InterruptVBlank = false
		return c.DoInterrupt(0x0040)
	}

	return 0
}

func (c *CPU) UpdateTimer() {
	delta := c.Cycle - c.LastTimerUpdateCycle
	if delta <= 0 {
		return
	}

	c.Mem[IODIV] += byte(delta / 4)
	c.UpdateTimerInterrupt(delta)
}

func (c *CPU) UpdateTimerInterrupt(delta int) {
	if (c.Mem[IOTAC] & (1 << 2)) == 0x00 {
		return
	}

	if c.TimerOverflow {
		c.Mem[IOTIMA] = c.Mem[IOTMA]
		c.InterruptTimer = true
		c.TimerOverflow = false
		return
	}

	c.TimerOverflow = false

	for i := 0; i < delta; i += 4 {
		c.Mem[IOTIMA]++

		// TIMA overflow
		if c.Mem[IOTIMA] == 0x00 {
			c.TimerOverflow = true
		}
	}
}
