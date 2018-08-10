package cpu

import "fmt"

type Registers struct {
	// Registers
	A             uint8  // Accumulator
	FlagZero      bool   // True if last result was 0
	FlagSubstract bool   // True if last operation was a substraction / decrement
	FlagHalfCarry bool   // True if the last operation had a half-carry
	FlagCarry     bool   // True if the last operation had a carry
	BC            uint16 // General-purpose register
	DE            uint16 // General-purpose register
	HL            uint16 // General-purpose register that doubles as a faster memory pointer
	SP            uint16 // Stack pointer
	PC            uint16 // Program counter
}

// Low/High byte setters for 16b registers
func (r *Registers) SetB(b byte) {
	r.BC = (r.BC & 0x00FF) | (uint16(b) << 8)
}
func (r *Registers) SetC(b byte) {
	r.BC = (r.BC & 0xFF00) | uint16(b)
}
func (r *Registers) SetD(b byte) {
	r.DE = (r.DE & 0x00FF) | (uint16(b) << 8)
}
func (r *Registers) SetE(b byte) {
	r.DE = (r.DE & 0xFF00) | uint16(b)
}
func (r *Registers) SetH(b byte) {
	r.HL = (r.HL & 0x00FF) | (uint16(b) << 8)
}
func (r *Registers) SetL(b byte) {
	r.HL = (r.HL & 0xFF00) | uint16(b)
}

func (r *Registers) GetA() byte {
	return r.A
}
func (r *Registers) SetA(b byte) {
	r.A = b
}

// Low/High byte getters for 16b registers
func (r *Registers) GetB() byte {
	return byte((r.BC & 0xFF00) >> 8)
}
func (r *Registers) GetC() byte {
	return byte(r.BC & 0x00FF)
}
func (r *Registers) GetD() byte {
	return byte((r.DE & 0xFF00) >> 8)
}
func (r *Registers) GetE() byte {
	return byte(r.DE & 0x00FF)
}
func (r *Registers) GetH() byte {
	return byte((r.HL & 0xFF00) >> 8)
}
func (r *Registers) GetL() byte {
	return byte(r.HL & 0x00FF)
}

func (r *Registers) ClearFlags() {
	r.FlagZero = false
	r.FlagSubstract = false
	r.FlagHalfCarry = false
	r.FlagCarry = false
}

// Returns the flags as a byte (ie. the F register)
func (r *Registers) GetF() byte {
	f := byte(0x00)

	if r.FlagZero {
		f |= 1 << 7
	}
	if r.FlagSubstract {
		f |= 1 << 6
	}
	if r.FlagHalfCarry {
		f |= 1 << 5
	}
	if r.FlagCarry {
		f |= 1 << 4
	}

	return f
}

// Sets the flags from a byte (ie. write to F register)
func (r *Registers) SetF(b byte) {
	r.FlagZero = (b & (1 << 7)) > 0
	r.FlagSubstract = (b & (1 << 6)) > 0
	r.FlagHalfCarry = (b & (1 << 5)) > 0
	r.FlagCarry = (b & (1 << 4)) > 0
}

// Returns the get/set function for a register given by name (eg. 'H')
func (r *Registers) GetRegisterCallbacks(name byte) (get func() byte, set func(byte)) {
	switch name {
	case 'A':
		return r.GetA, r.SetA
	case 'B':
		return r.GetB, r.SetB
	case 'C':
		return r.GetC, r.SetC
	case 'D':
		return r.GetD, r.SetD
	case 'E':
		return r.GetE, r.SetE
	case 'H':
		return r.GetH, r.SetH
	case 'L':
		return r.GetL, r.SetL
	}

	panic("unreachable")
}

func (r *Registers) GetRegisterAddress(name string) *uint16 {
	switch name {
	case "BC":
		return &r.BC
	case "DE":
		return &r.DE
	case "HL":
		return &r.HL
	case "SP":
		return &r.SP
	}

	panic("unreachable")
}

func (r *Registers) String() string {
	flags := [4]byte{'-', '-', '-', '-'}
	if r.FlagZero {
		flags[0] = 'Z'
	}
	if r.FlagSubstract {
		flags[1] = 'N'
	}
	if r.FlagHalfCarry {
		flags[2] = 'H'
	}
	if r.FlagCarry {
		flags[3] = 'C'
	}

	return fmt.Sprintf(
		"A:%02X BC:%04X DE:%04X HL:%04X SP:%04X PC:%04X Flags:%s",
		r.A, r.BC, r.DE, r.HL, r.SP, r.PC, flags,
	)
}

func (r *Registers) MarshalBinary() ([]byte, error) {
	return []byte{
			r.GetA(),
			r.GetF(),
			r.GetB(),
			r.GetC(),
			r.GetD(),
			r.GetE(),
			r.GetH(),
			r.GetL(),
			byte((r.PC & 0x0F)),
			byte((r.PC & 0xF0) >> 8),
			byte((r.SP & 0x0F)),
			byte((r.SP & 0xF0) >> 8),
		},
		nil
}

func (r *Registers) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return fmt.Errorf("data less than 12 bytes: %d:", len(data))
	}

	r.SetA(data[0])
	r.SetF(data[1])
	r.SetB(data[2])
	r.SetC(data[3])
	r.SetD(data[4])
	r.SetE(data[5])
	r.SetH(data[6])
	r.SetL(data[7])
	r.PC = uint16(data[8]) | (uint16(data[9]) << 8)
	r.SP = uint16(data[10]) | (uint16(data[11]) << 8)

	return nil
}
