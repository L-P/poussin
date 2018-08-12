package debugger

import (
	"sync/atomic"

	"github.com/jroimartin/gocui"
)

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
		{d.cbStepToPC, 'i'},
		{d.cbStepWhenSB, 'o'},
	}

	for _, v := range binds {
		if err := d.gui.SetKeybinding(
			"instructions",
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

func (d *Debugger) cbStepToPC(g *gocui.Gui, v *gocui.View) error {
	if d.hasModal.IsSet() {
		return nil
	}

	cb := func(addr uint16) {
		d.stepToPC = addr
		atomic.StoreInt32(&d.flowState, FlowStepToPC)
	}

	atomic.StoreInt32(&d.flowState, FlowPause)
	if err := d.inputUInt16Modal(g, "Jump to PC", cb); err != nil {
		return err
	}

	return nil
}

func (d *Debugger) cbStepWhenSB(g *gocui.Gui, v *gocui.View) error {
	if d.hasModal.IsSet() {
		return nil
	}

	cb := func(v byte) {
		d.stopWhenSB = v
		atomic.StoreInt32(&d.flowState, FlowStopWhenSB)
	}

	atomic.StoreInt32(&d.flowState, FlowPause)
	if err := d.inputUInt8Modal(g, "Stop when SB=", cb); err != nil {
		return err
	}

	return nil
}
