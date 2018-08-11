package cpu

// Rotates a byte to the left using a carry as the 9th bit, returns the new
// carry bit state.
func rotateLeftWithCarry(v byte, carry bool) (byte, bool) {
	oldCarry := uint8(0)
	if carry {
		oldCarry = uint8(1)
	}

	newCarry := (v & (1 << 7)) > 0

	return (v << 1) | oldCarry, newCarry
}

// Rotates a byte to the right using a carry as the 9th bit, returns the new
// carry bit state.
func rotateRightWithCarry(v byte, carry bool) (byte, bool) {
	oldCarry := uint8(0)
	if carry {
		oldCarry = uint8(1)
	}

	newCarry := v&0x01 == 0x01

	return (v >> 1) | (oldCarry << 7), newCarry
}

// Returns decremented byte and half carry flag
func decrement(v byte) (byte, bool) {
	return v - 1, (v & 0xF) == 0x00
}

// Returns incremented byte and half carry flag
func increment(v byte) (byte, bool) {
	return v + 1, (((v & 0xF) + 1) & 0x10) > 0
}

// Offset an address interpreting the given offset as a signed byte
func signedOffset(base uint16, offset byte) uint16 {
	return uint16(int16(base) + int16(int8(offset)))
}
