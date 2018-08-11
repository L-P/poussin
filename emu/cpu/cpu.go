package cpu

import (
	"errors"
	"fmt"

	"home.leo-peltier.fr/poussin/emu/ppu"
	"home.leo-peltier.fr/poussin/emu/rom"
)

type CPU struct {
	Registers

	Mem  [0xFFFF]byte
	Boot [256]byte
	ROM  [1024 * 1024 * 8]byte // Per Wikipedia, a GB ROM is 8 MiB max

	PPU *ppu.PPU

	// Set by the STOP instruction, reset by interrupts
	Stopped bool

	// For debugging purposes
	LastOpcode      byte
	LastOpcodeWasCB bool
	LastLowArg      byte
	LastHighArg     byte

	// Adressable at 0xFFFF
	InterruptEnable byte
	InterruptTimer  bool

	// Interrupt Master Flag, not addressable
	InterruptMaster bool

	// Switching to CB opcode (0xCB for code bank?) is an instruction in
	// itself, when this flag is set it means the next opcode we read will be
	// from the CB opcodes
	NextOpcodeIsCB bool

	// Current clock cycle number
	Cycle int

	// On which cycle we last updated the timer registers
	LastTimerUpdateCycle int

	// TIMA overflowed on last cycle
	TimerOverflow bool
}

func New(ppu *ppu.PPU) CPU {
	return CPU{
		PPU: ppu,
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

	c.LastOpcode = opcode
	c.LastOpcodeWasCB = cb // CB is reset after Decode, hence cb
	c.LastLowArg = l
	c.LastHighArg = h

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
