package debugger

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/tevino/abool"
	"home.leo-peltier.fr/poussin/emu/cpu"
	"home.leo-peltier.fr/poussin/emu/ppu"
)

// registers + opcode + CB + low arg + high arg
const insBufferStride = 12 + 4

type Debugger struct {
	cpu *cpu.CPU
	ppu *ppu.PPU

	gui *gocui.Gui

	lastCPUError error
	quit         *abool.AtomicBool
	pause        *abool.AtomicBool
	stepOnce     *abool.AtomicBool

	// Keep 1280 Kio worth of history, which is quite a lot of instructions
	insBuffer              [insBufferStride * 40960]byte
	curInsBufferWriteIndex int
	msgBuffer              string
	insLastRenderedOpCount int

	opCount         int
	lastPerfOpCount int
	frameCount      int
	lastPerfDisplay time.Time

	opPerSecond    int
	framePerSecond int
}

func New(c *cpu.CPU, p *ppu.PPU) Debugger {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}

	d := Debugger{
		cpu: c,
		ppu: p,
		gui: gui,

		pause:    abool.New(),
		quit:     abool.New(),
		stepOnce: abool.New(),
	}

	d.gui.SetManagerFunc(d.layout)
	d.initKeybinds()

	return d
}

func (d *Debugger) initKeybinds() {
	if err := d.gui.SetKeybinding("", 'p', gocui.ModNone, d.cbPause); err != nil {
		panic(err)
	}

	if err := d.gui.SetKeybinding("", 'q', gocui.ModNone, d.cbQuit); err != nil {
		panic(err)
	}

	if err := d.gui.SetKeybinding("", 'l', gocui.ModNone, d.cbStep); err != nil {
		panic(err)
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
	d.quit.Set()
	return nil
}

func (d *Debugger) Close() {
	d.quit.Set()
	d.pause.Set()
	d.gui.Close()
}

func (d *Debugger) layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("instructions", 0, 0, 80, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
	}

	if _, err := g.SetView("messages", 0, maxY-4, 64, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView("perf", 65, maxY-4, 80, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}

func (d *Debugger) Run(quit chan bool) {
	go func() {
		if err := d.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
			// HACK: When using pprof gocui throws this, this is weird and should be investigated.
			if err.Error() != "invalid dimensions" {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for !d.quit.IsSet() {
		select {
		case <-ticker.C:
			d.gui.Update(d.updateGUI)
		case <-quit:
			d.quit.Set()
			return
		}
	}

	quit <- true
}

func (d *Debugger) Panic(err error) {
	d.lastCPUError = err
	d.updateMessages()
	d.pause.Set()
}

func (d *Debugger) Update() {
	d.updateInstructions()
	d.updateMessages()

	for d.pause.IsSet() {
		time.Sleep(time.Duration(100 * time.Millisecond))
		if d.stepOnce.IsSet() {
			d.stepOnce.UnSet()
			break
		}
	}
}

func (d *Debugger) updateGUI(g *gocui.Gui) error {
	if err := d.updatePerfWindow(g); err != nil {
		return err
	}

	if err := d.updateInsWindow(g); err != nil {
		return err
	}

	return nil
}

func (d *Debugger) updateInstructions() {
	var cb byte
	if d.cpu.LastOpcodeWasCB {
		cb = 0x01
	}

	d.cpu.WriteToArray(d.insBuffer[:], d.curInsBufferWriteIndex)
	d.insBuffer[d.curInsBufferWriteIndex+12+0] = d.cpu.LastOpcode
	d.insBuffer[d.curInsBufferWriteIndex+12+1] = cb
	d.insBuffer[d.curInsBufferWriteIndex+12+2] = d.cpu.LastLowArg
	d.insBuffer[d.curInsBufferWriteIndex+12+3] = d.cpu.LastHighArg
	d.curInsBufferWriteIndex = (d.curInsBufferWriteIndex + insBufferStride) % len(d.insBuffer)
}

func (d *Debugger) updateMessages() {
	if d.lastCPUError == nil {
		d.updatePerfCounters()
		return
	}

	d.msgBuffer = d.lastCPUError.Error()
}

func (d *Debugger) updatePerfCounters() {
	d.opCount++
	now := time.Now()
	if now.Sub(d.lastPerfDisplay) >= time.Duration(1*time.Second) {
		d.lastPerfDisplay = now

		d.opPerSecond = d.opCount - d.lastPerfOpCount
		d.lastPerfOpCount = d.opCount

		d.framePerSecond = d.ppu.PushedFrames - d.frameCount
		d.frameCount = d.ppu.PushedFrames
	}
}

func (d *Debugger) updatePerfWindow(g *gocui.Gui) error {
	v, err := g.View("perf")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(v, "OPS: %d\nFPS: %d", d.opPerSecond, d.framePerSecond)

	if d.msgBuffer != "" {
		msgView, err := g.View("messages")
		if err != nil {
			return err
		}

		msgView.Clear()
		fmt.Fprintln(msgView, d.msgBuffer)
	}

	return nil
}

func (d *Debugger) updateInsWindow(g *gocui.Gui) error {
	if !d.pause.IsSet() {
		return nil
	}

	if d.insLastRenderedOpCount >= d.opCount {
		return nil
	}

	view, err := g.View("instructions")
	if err != nil {
		return err
	}

	for i := d.curInsBufferWriteIndex; i != d.curInsBufferWriteIndex-insBufferStride; i = (i + insBufferStride) % len(d.insBuffer) {
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

		fmt.Fprintf(
			view,
			"%-22s %s\n",
			ins.String(l, h),
			registers.String(),
		)
	}

	d.insLastRenderedOpCount = d.opCount

	return nil
}

func (d *Debugger) Pause() {
	d.pause.Set()
}

func (d *Debugger) clearInstructionsView() {
	v, err := d.gui.View("instructions")
	if err != nil {
		return
	}
	v.Clear()
}
