package debugger

import (
	"bytes"
	"fmt"
	"os"
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

	gui    *gocui.Gui
	closed *abool.AtomicBool

	// View buffers
	insBuffer              [insBufferStride * 256]byte
	curInsBufferWriteIndex int
	msgBuffer              bytes.Buffer
	lastCPUError           error

	// I/O registers
	ioIF      byte
	ioIE      byte
	ioIMaster bool
	ioDIV     byte
	ioTMA     byte
	ioTAC     byte
	ioTIMA    byte

	// Performance counters
	opCount         int
	lastPerfOpCount int
	frameCount      int
	lastPerfDisplay time.Time
	opPerSecond     int
	framePerSecond  int
}

func New(c *cpu.CPU, p *ppu.PPU) (*Debugger, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	d := Debugger{
		cpu:    c,
		ppu:    p,
		gui:    gui,
		closed: abool.New(),
	}

	d.gui.SetManagerFunc(d.layout)
	if err := d.initKeybinds(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Debugger) Close() {
	// If gocui/termbox was closed once alreay it will block on reading a channel
	if !d.closed.IsSet() {
		d.gui.Close()
	}
}

func (d *Debugger) RunGUI(shouldClose <-chan bool) {
	guiClosed := make(chan bool)
	go func() {
		if err := d.gui.MainLoop(); err != nil {
			if err != gocui.ErrQuit && err.Error() != "invalid dimensions" {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		guiClosed <- true
	}()

	select {
	case <-shouldClose:
		d.gui.Update(func(*gocui.Gui) error { return gocui.ErrQuit })
		d.gui.Close()
		<-guiClosed // wait for GUI to exit
	case <-guiClosed:
	}

	d.closed.Set()
	d.gui.Close()
}

func (d *Debugger) Closed() bool {
	return d.closed.IsSet()
}

func (d *Debugger) Panic(err error) {
	d.lastCPUError = err
	d.updateMessages()
}

// Update updates the debugger internal state from the current CPU/PPU state
// and will block if the user requested a pause or breakpoint.
func (d *Debugger) Update() {
	d.updateInstructions()
	d.updateMessages()
	d.updateIORegisters()
	d.updatePerfCounters()
}

func (d *Debugger) updateGUI(g *gocui.Gui) error {
	if err := d.updatePerfWindow(g); err != nil {
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

func (d *Debugger) updateIORegisters() {
	d.ioIF = d.cpu.Mem[cpu.IOIF]
	d.ioIE = d.cpu.InterruptEnable
	d.ioIMaster = d.cpu.InterruptMaster
	d.ioDIV = d.cpu.Mem[cpu.IODIV]
	d.ioTMA = d.cpu.Mem[cpu.IOTMA]
	d.ioTAC = d.cpu.Mem[cpu.IOTAC]
	d.ioTIMA = d.cpu.Mem[cpu.IOTIMA]
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
	if d.lastCPUError != nil {
		d.msgBuffer.WriteString(d.lastCPUError.Error())
		return
	}

	if d.cpu.SBBuffer.Len() > 0 {
		d.msgBuffer.Write(d.cpu.SBBuffer.Bytes())
		d.cpu.SBBuffer.Reset()
	}
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
