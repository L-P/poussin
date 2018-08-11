package debugger

import (
	"github.com/jroimartin/gocui"
)

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
			"performance",
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
				view.Wrap = true
			}
		}
	}

	return nil
}

func (d *Debugger) initKeybinds() error {
	binds := []struct {
		handler keybindHandler
		key     rune
	}{
		{d.cbPause, 'p'},
		{d.cbQuit, 'q'},
		{d.cbStep, 'j'},
	}

	for _, v := range binds {
		if err := d.gui.SetKeybinding(
			"",
			v.key,
			gocui.ModNone,
			d.proxyCallback(v.handler),
		); err != nil {
			return err
		}
	}

	return nil
}

type keybindHandler func(g *gocui.Gui, v *gocui.View) error

func (d *Debugger) proxyCallback(cb keybindHandler) keybindHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		err := cb(g, v)
		d.updateGUI(g)
		return err
	}
}

func (d *Debugger) cbStep(g *gocui.Gui, v *gocui.View) error {
	if !d.pause.IsSet() {
		return nil
	}

	d.stepOnce.Set()

	return nil
}

func (d *Debugger) cbPause(g *gocui.Gui, v *gocui.View) error {
	if d.pause.IsSet() {
		d.clearInstructionsView()
		d.pause.UnSet()
	} else {
		d.pause.Set()
	}

	return nil
}

func (d *Debugger) cbQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
