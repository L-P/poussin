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
