package debugger

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/L-P/poussin/emu/cpu"
)

func (d *Debugger) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	iW := 16
	msgW := maxX - (iW * 2) - 1
	msgH := 10
	memW := 16

	views := []struct {
		name string
		x1   int
		y1   int
		x2   int
		y2   int
	}{
		{
			"instructions",
			0,
			0,
			maxX - memW - 1,
			maxY - msgH - 1,
		},
		{
			"memory",
			maxX - memW,
			0,
			maxX - 1,
			maxY - msgH - 1,
		},
		{
			"messages",
			0,
			maxY - msgH,
			msgW,
			maxY - 1,
		},
		{
			"IO registers",
			msgW + 1,
			maxY - msgH,
			msgW + iW,
			maxY - 1,
		},
		{
			"misc",
			msgW + iW + 1,
			maxY - msgH,
			maxX - 1,
			maxY - 1,
		},
	}

	for _, v := range views {
		if view, err := g.SetView(v.name, v.x1, v.y1, v.x2, v.y2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			view.Title = v.name
			switch v.name {
			case "memory":
				view.Autoscroll = true
			case "instructions":
				view.Autoscroll = true
			case "messages":
				view.Autoscroll = true
				view.Wrap = true
			}
		}
	}

	return nil
}

func (d *Debugger) updateGUI(g *gocui.Gui) error {
	d.Lock()
	defer d.Unlock()

	if !d.hasModal.IsSet() {
		g.SetCurrentView("instructions")
	}

	if err := d.updateMiscWindow(g); err != nil {
		return err
	}

	if err := d.updateIORegistersWindow(g); err != nil {
		return err
	}

	if err := d.updateMsgWindow(g); err != nil {
		return err
	}

	if err := d.updateInsWindow(g); err != nil {
		return err
	}

	if err := d.updateMemoryWindow(g); err != nil {
		return err
	}

	return nil
}

func (d *Debugger) updateMsgWindow(g *gocui.Gui) error {
	if d.msgBuffer.Len() <= 0 {
		return nil
	}

	msgView, err := g.View("messages")
	if err != nil {
		return err
	}

	msgView.Clear()
	fmt.Fprintln(msgView, d.msgBuffer.String())

	return nil
}

func (d *Debugger) updateIORegistersWindow(g *gocui.Gui) error {
	v, err := g.View("IO registers")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintf(v, "IF:      %s\n", ieString(d.ioIF))
	fmt.Fprintf(v, "IE:      %s\n", ieString(d.ioIE))
	fmt.Fprintf(v, "IME:     %t\n", d.ioIME)
	fmt.Fprintf(v, "DIV:     %02X\n", d.ioDIV)
	fmt.Fprintf(v, "IDIV:    %04X\n", d.ioInternalDIV)
	fmt.Fprintf(v, "TMA:     %02X\n", d.ioTMA)
	fmt.Fprintf(v, "TAC:     %02X\n", d.ioTAC)
	fmt.Fprintf(v, "TIMA:    %02X\n", d.ioTIMA)

	return nil
}

var ieFlagNames = [5]byte{'V', 'L', 'T', 'S', 'J'}

// Returns a human-readable version of the IE flag
func ieString(v byte) string {
	ret := [5]byte{'-', '-', '-', '-', '-'}
	for i := range ieFlagNames {
		if (v & (1 << byte(i))) > 0 {
			ret[i] = ieFlagNames[i]
		}
	}

	return fmt.Sprintf("%s", ret)
}

func (d *Debugger) updateMiscWindow(g *gocui.Gui) error {
	v, err := g.View("misc")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
        "OPS: %d\nFPS: %d\nDepth: %d\nLCDC: %02X\nSTAT: %02X\nP1: %02X",
		d.opPerSecond,
		d.framePerSecond,
		d.callDepth,
		d.ioLCDC,
		d.ioSTAT,
		d.ioP1,
	)

	return nil
}

func (d *Debugger) updateMemoryWindow(g *gocui.Gui) error {
	v, err := g.View("memory")
	if err != nil {
		return err
	}

	v.Clear()
	for j := 0; j < len(d.memBuffer)/memBufferStride; j++ {
		i := (d.curMemBufferWriteIndex + (j * memBufferStride)) % len(d.memBuffer)

		rw := d.memBuffer[i+2]
		if rw != 0x01 && rw != 0x02 {
			return nil
		}

		rwFlag := 'R'
		if rw == 0x02 {
			rwFlag = 'W'
		}

		pc := uint16(d.memBuffer[i]) | (uint16(d.memBuffer[i+1]) << 8)
		addr := uint16(d.memBuffer[i+3]) | (uint16(d.memBuffer[i+4]) << 8)
		val := d.memBuffer[i+5]

		fmt.Fprintf(v, "%04X %c %04X %02X\n", pc, rwFlag, addr, val)
	}

	return nil
}

func (d *Debugger) updateInsWindow(g *gocui.Gui) error {
	view, err := g.View("instructions")
	if err != nil {
		return err
	}

	view.Clear()
	var prevRegisters cpu.Registers
	for j := 0; j < len(d.insBuffer)/insBufferStride; j++ {
		i := (d.curInsBufferWriteIndex + (j * insBufferStride)) % len(d.insBuffer)

		registers := cpu.ReadFromArray(d.insBuffer[:], i)
		opcode := d.insBuffer[i+12+0]
		l := d.insBuffer[i+12+2]
		h := d.insBuffer[i+12+3]

		var cb bool
		if d.insBuffer[i+12+1] != 0x00 {
			cb = true
		}

		ins := cpu.Decode(opcode, cb)
		if !ins.Valid() {
			continue
		}

		d.printInstruction(view, ins, l, h, registers, prevRegisters)
		prevRegisters = registers
	}

	return nil
}

func (d *Debugger) printInstruction(
	view *gocui.View,
	ins cpu.Instruction,
	l byte,
	h byte,
	regs cpu.Registers,
	prev cpu.Registers,
) {
	prevStr := prev.String()
	curStr := regs.String()

	var final bytes.Buffer

	if len(prevStr) != len(curStr) {
		panic("len(prevStr) != len(curStr)")
	}

	var diff bool
	for i, _ := range prevStr {
		if prevStr[i] != curStr[i] {
			if !diff {
				diff = true
				final.WriteString("\x1b[1;31m")
			}
		} else {
			diff = false
			final.WriteString("\x1b[0m")
		}

		final.WriteByte(curStr[i])
	}

	final.WriteString("\x1b[0m")

	fmt.Fprintf(
		view,
		"%-18s %s\n",
		ins.String(l, h),
		final.String(),
	)
}

// Creates a small centered editable window
func inputModalView(g *gocui.Gui, title string) (*gocui.View, error) {
	maxX, maxY := g.Size()
	w := 32
	h := 2
	x := maxX/2 - w/2
	y := maxY/2 - h/2

	v, err := g.SetView(title, x, y, x+w, y+h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
	}

	v.Editable = true
	v.Title = title
	g.SetViewOnTop(title)
	g.SetCurrentView(title)
	v.Clear()

	return v, nil
}

func (d *Debugger) inputUIntModal(g *gocui.Gui, title string, intWidth int, cb func(int64)) error {
	if d.hasModal.IsSet() {
		return nil
	}
	d.hasModal.Set()

	if _, err := inputModalView(g, title); err != nil {
		return err
	}

	if err := g.SetKeybinding(
		title,
		gocui.KeyEnter,
		gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			buf := strings.Trim(v.Buffer(), "\n ")

			if buf != "" {
				s, err := strconv.ParseUint(buf, 16, intWidth)
				if err != nil {
					d.msgBuffer.WriteString(err.Error() + "\n")
					return nil
				}
				cb(int64(s))
			}

			d.hasModal.UnSet()
			g.DeleteKeybindings(title)
			return g.DeleteView(title)
		},
	); err != nil {
		return err
	}

	return nil
}

func (d *Debugger) inputUInt8Modal(g *gocui.Gui, title string, cb func(byte)) error {
	return d.inputUIntModal(g, title, 8, func(v int64) { cb(byte(v)) })
}

func (d *Debugger) inputUInt16Modal(g *gocui.Gui, title string, cb func(uint16)) error {
	return d.inputUIntModal(g, title, 16, func(v int64) { cb(uint16(v)) })
}
