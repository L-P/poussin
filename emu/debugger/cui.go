package debugger

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"github.com/jroimartin/gocui"
	"home.leo-peltier.fr/poussin/emu/cpu"
)

func (d *Debugger) updateGUI(g *gocui.Gui) error {
	d.Lock()
	defer d.Unlock()

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
	v, err := g.View("I/O registers")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintf(v, "IF:      %02X\n", d.ioIF)
	fmt.Fprintf(v, "IE:      %02X\n", d.ioIE)
	fmt.Fprintf(v, "IMaster: %t\n", d.ioIMaster)
	fmt.Fprintf(v, "DIV:     %02X\n", d.ioDIV)
	fmt.Fprintf(v, "TMA:     %02X\n", d.ioTMA)
	fmt.Fprintf(v, "TAC:     %02X\n", d.ioTAC)
	fmt.Fprintf(v, "TIMA:    %02X\n", d.ioTIMA)

	return nil
}

func (d *Debugger) updateMiscWindow(g *gocui.Gui) error {
	v, err := g.View("misc")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
		"OPS: %d\nFPS: %d\nDepth: %d\n",
		d.opPerSecond,
		d.framePerSecond,
		d.callDepth,
	)

	return nil
}

func (d *Debugger) updateInsWindow(g *gocui.Gui) error {
	view, err := g.View("instructions")
	if err != nil {
		return err
	}

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
		"%-22s %s\n",
		ins.String(l, h),
		final.String(),
	)
}

func (d *Debugger) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	iW := 17
	msgW := maxX - (iW * 2)
	msgH := 9

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
			"I/O registers",
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

type keybindHandler func(g *gocui.Gui, v *gocui.View) error

func (d *Debugger) initKeybinds() error {
	binds := []struct {
		handler keybindHandler
		key     rune
	}{
		{d.cbPause, 'p'},
		{d.cbQuit, 'q'},
		{d.cbStepOut, 'h'},
		{d.cbStepOver, 'j'},
		{d.cbStepIn, 'l'},
	}

	for _, v := range binds {
		if err := d.gui.SetKeybinding(
			"",
			v.key,
			gocui.ModNone,
			d.keybindWrapper(v.handler),
		); err != nil {
			return err
		}
	}

	return nil
}

func (d *Debugger) keybindWrapper(h keybindHandler) keybindHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		d.Lock()
		defer d.Unlock()
		return h(g, v)
	}
}

func (d *Debugger) cbStepOut(g *gocui.Gui, v *gocui.View) error {
	if atomic.LoadInt32(&d.flowState) == FlowPause {
		atomic.StoreInt32(&d.flowState, FlowStepOut)
	}
	d.requestedDepth = d.callDepth - 1

	return nil
}

func (d *Debugger) cbStepOver(g *gocui.Gui, v *gocui.View) error {
	if atomic.LoadInt32(&d.flowState) == FlowPause {
		atomic.StoreInt32(&d.flowState, FlowStepOver)
	}

	return nil
}

func (d *Debugger) cbStepIn(g *gocui.Gui, v *gocui.View) error {
	if atomic.LoadInt32(&d.flowState) == FlowPause {
		atomic.StoreInt32(&d.flowState, FlowStepIn)
		d.requestedDepth = d.callDepth + 1
	}

	return nil
}

func (d *Debugger) cbPause(g *gocui.Gui, v *gocui.View) error {
	if atomic.LoadInt32(&d.flowState) == FlowPause {
		atomic.StoreInt32(&d.flowState, FlowRun)
	} else {
		atomic.StoreInt32(&d.flowState, FlowPause)
	}

	return nil
}

func (d *Debugger) cbQuit(g *gocui.Gui, v *gocui.View) error {
	atomic.StoreInt32(&d.flowState, FlowQuit)
	return gocui.ErrQuit
}
