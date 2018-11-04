package cpu

// JoypadState represents the state of every button on the DMG joypad as should
// be fed to the CPU by whoever is in charge of listening to user input.
type JoypadState struct {
	Up     bool
	Down   bool
	Left   bool
	Right  bool
	A      bool
	B      bool
	Start  bool
	Select bool
}

func (c *CPU) updateJoypad() {
	select {
	case curState := <-c.JoypadInput:
		if curState != c.Joypad {
			c.Joypad = curState
			p1 := c.FetchIOP1()
			if p1&(1<<5|1<<4) == 0 && c.IEEnabled(IEJoypad) {
				c.SetIF(IEJoypad)
			}
		}
	default:
	}
}
