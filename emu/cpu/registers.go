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

// Returns the address of a 16b register by its name.
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

// Returns a human-readable version of the registers state.
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
		"A:%02X BC:%04X DE:%04X HL:%04X SP:%04X PC:%04X %s",
		r.A, r.BC, r.DE, r.HL, r.SP, r.PC, flags,
	)
}

// Writes 12 bytes of register data to the given array starting at the given offset.
func (r *Registers) WriteToArray(a []byte, i int) {
	a[i+0] = r.GetA()
	a[i+1] = r.GetF()
	a[i+2] = r.GetB()
	a[i+3] = r.GetC()
	a[i+4] = r.GetD()
	a[i+5] = r.GetE()
	a[i+6] = r.GetH()
	a[i+7] = r.GetL()
	a[i+8] = byte((r.PC & 0xFF00) >> 8)
	a[i+9] = byte((r.PC & 0x00FF))
	a[i+10] = byte((r.SP & 0xFF00) >> 8)
	a[i+11] = byte((r.SP & 0x00FF))
}

// Creates a new Registers object from an array and a base offset.
func ReadFromArray(a []byte, i int) Registers {
	r := Registers{
		A:  a[i+0],
		BC: (uint16(a[i+2]) << 8) | uint16(a[i+3]),
		DE: (uint16(a[i+4]) << 8) | uint16(a[i+5]),
		HL: (uint16(a[i+6]) << 8) | uint16(a[i+7]),
		PC: (uint16(a[i+8]) << 8) | uint16(a[i+9]),
		SP: (uint16(a[i+10]) << 8) | uint16(a[i+11]),
	}
	r.SetF(a[i+1])

	return r
}
