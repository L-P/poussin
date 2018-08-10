package debugger

import (
	"container/ring"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/tevino/abool"
	"home.leo-peltier.fr/poussin/emu/cpu"
	"home.leo-peltier.fr/poussin/emu/ppu"
)

type Debugger struct {
	cpu *cpu.CPU
	ppu *ppu.PPU

	gui *gocui.Gui

	lastCPUError error
	quit         *abool.AtomicBool
	pause        *abool.AtomicBool

	out *os.File

	insBuffer              *ring.Ring
	msgBuffer              string
	insLastRenderedOpCount int

	opCount         int
	lastPerfOpCount int
	frameCount      int
	lastPerfDisplay time.Time

	opPerSecond    int
	framePerSecond int
}

type instruction struct {
	opcode    byte
	cb        bool
	h         byte
	l         byte
	registers []byte
}

func New(c *cpu.CPU, p *ppu.PPU) Debugger {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}

	out, err := ioutil.TempFile("/dev/shm", "poussin.out.")
	if err != nil {
		panic(err)
	}

	d := Debugger{
		cpu: c,
		ppu: p,
		gui: gui,
		out: out,

		insBuffer: ring.New(2048),
		pause:     abool.New(),
		quit:      abool.New(),
	}

	d.gui.SetManagerFunc(d.layout)
	d.initKeybinds()

	return d
}

func (d *Debugger) initKeybinds() {
	if err := d.gui.SetKeybinding("", 'p', gocui.ModNone, d.cbPause); err != nil {
		panic(err)
	}
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

func (d *Debugger) Close() {
	d.quit.Set()
	d.pause.Set()
	d.gui.Close()
	fmt.Println(d.out.Name())
	d.out.Close()
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

func (d *Debugger) Run(quit chan int) {
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
}

func (d *Debugger) Panic(err error) {
	d.lastCPUError = err
	d.pause.Set()
}

func (d *Debugger) Update() {
	for d.pause.IsSet() {
		time.Sleep(time.Duration(100 * time.Millisecond))
	}

	d.updateInstructions()
	d.updateMessages()
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
	for d.pause.IsSet() {
		time.Sleep(time.Duration(100 * time.Millisecond))
	}

	registers, _ := d.cpu.MarshalBinary()
	d.insBuffer.Value = instruction{
		opcode:    d.cpu.LastOpcode,
		cb:        d.cpu.LastOpcodeWasCB,
		l:         d.cpu.LastLowArg,
		h:         d.cpu.LastHighArg,
		registers: registers,
	}
	d.insBuffer = d.insBuffer.Next()
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

	d.insBuffer.Do(func(v interface{}) {
		if v == nil {
			return
		}

		b := v.(instruction)

		ins := cpu.Decode(b.opcode, b.cb)
		if !ins.Valid() {
			return
		}

		var registers cpu.Registers
		if err := registers.UnmarshalBinary(b.registers); err != nil {
			return
		}

		fmt.Fprintf(
			view,
			"%-22s %s\n",
			ins.String(b.l, b.h),
			registers.String(),
		)
	})

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
