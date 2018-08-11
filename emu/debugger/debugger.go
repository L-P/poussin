package debugger

import (
	"bytes"
	"fmt"
	"log"
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

	gui  *gocui.Gui
	quit chan bool

	lastCPUError error
	pause        *abool.AtomicBool
	stepOnce     *abool.AtomicBool

	insBuffer              [insBufferStride * 256]byte
	curInsBufferWriteIndex int
	msgBuffer              string

	opCount         int
	lastPerfOpCount int
	frameCount      int
	lastPerfDisplay time.Time

	opPerSecond    int
	framePerSecond int
}

func New(c *cpu.CPU, p *ppu.PPU) (*Debugger, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	d := Debugger{
		cpu: c,
		ppu: p,
		gui: gui,

		pause:    abool.New(),
		stepOnce: abool.New(),
	}

	d.gui.SetManagerFunc(d.layout)
	if err := d.initKeybinds(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Debugger) Close() {
	d.pause.Set()
	d.gui.Close()
}

func (d *Debugger) RunGUI(quit chan bool) {
	d.quit = quit

	if err := d.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		// HACK: When using pprof gocui throws this, this is weird and should be investigated.
		log.Printf("gocui.Gui.MainLoop stopped: %s", err)
	}

	quit <- true
}

func (d *Debugger) Panic(err error) {
	d.lastCPUError = err
	d.updateMessages()
	d.pause.Set()
	d.gui.Update(d.updateGUI)
}

// Update updates the debugger internal state from the current CPU/PPU state
// and will block if the user requested a pause or breakpoint.
func (d *Debugger) Update() {
	d.updateInstructions()
	d.updateMessages()
	d.updatePerfCounters()

	if d.stepOnce.IsSet() {
		d.stepOnce.UnSet()
		d.pause.Set()
		d.gui.Update(d.updateGUI)
	}

	for d.pause.IsSet() {
		time.Sleep(time.Duration(16 * time.Millisecond))

		if d.stepOnce.IsSet() {
			break
		}
	}
}

func (d *Debugger) updateGUI(g *gocui.Gui) error {
	// Don't update unpaused, we don't want our data to get written as we read it.
	if !d.pause.IsSet() {
		return nil
	}

	if err := d.updatePerfWindow(g); err != nil {
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

		d.gui.Update(d.updatePerfWindow)
	}
}

func (d *Debugger) updateMsgWindow(g *gocui.Gui) error {
	if d.msgBuffer == "" {
		return nil
	}

	msgView, err := g.View("messages")
	if err != nil {
		return err
	}

	msgView.Clear()
	fmt.Fprintln(msgView, d.msgBuffer)

	return nil
}

func (d *Debugger) updatePerfWindow(g *gocui.Gui) error {
	v, err := g.View("performance")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(v, "OPS: %d\nFPS: %d", d.opPerSecond, d.framePerSecond)

	return nil
}

func (d *Debugger) updateInsWindow(g *gocui.Gui) error {
	if !d.pause.IsSet() {
		return nil
	}

	view, err := g.View("instructions")
	if err != nil {
		return err
	}

	var prevRegisters cpu.Registers
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
				final.WriteByte(0x1B)
				final.WriteString("[1;31m")
			}
		} else {
			diff = false
			final.WriteByte(0x1B)
			final.WriteString("[0m")
		}

		final.WriteByte(curStr[i])
	}

	final.WriteByte(0x1B)
	final.WriteString("[0m")

	fmt.Fprintf(
		view,
		"%-22s %s\n",
		ins.String(l, h),
		final.String(),
	)
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
