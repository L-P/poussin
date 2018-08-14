package debugger

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/tevino/abool"
	"home.leo-peltier.fr/poussin/emu/cpu"
	"home.leo-peltier.fr/poussin/emu/ppu"
)

// registers + opcode + CB + low arg + high arg
const insBufferStride = 12 + 4
const insBufferCount = 128

// PC + r/w + addr + value
const memBufferStride = 2 + 1 + 2 + 1
const memBufferCount = 128

// Debugger is a CLI for debugging a GameBoy ROM execution.
type Debugger struct {
	// We have two routines, one for writing the CPU state and one for
	// displaying it, we embed a mutex to safely run the routines concurrently.
	sync.Mutex

	cpu *cpu.CPU
	ppu *ppu.PPU

	// Set to true when the debugger has quit for any reason.
	closed   *abool.AtomicBool
	gui      *gocui.Gui
	hasModal *abool.AtomicBool

	// Flow control
	flowState      int32
	stepToPC       uint16
	stopWhenSB     byte
	stepToOpcode   byte
	requestedDepth int
	// _will_ be negative, nothing prevents you from pushing the stack and RET without having a CALL
	callDepth int

	// View buffers
	insBuffer              [insBufferStride * insBufferCount]byte
	curInsBufferWriteIndex int
	memBuffer              [memBufferStride * memBufferCount]byte
	curMemBufferWriteIndex int
	msgBuffer              bytes.Buffer
	lastCPUError           error

	// I/O registers
	ioIF          byte
	ioIE          byte
	ioIME         bool
	ioDIV         byte
	ioTMA         byte
	ioTAC         byte
	ioTIMA        byte
	ioInternalDIV uint16

	// Performance counters
	opCount         int
	lastPerfOpCount int
	frameCount      int
	lastPerfDisplay time.Time
	opPerSecond     int
	framePerSecond  int
}

const (
	FlowRun = int32(iota)
	FlowQuit
	FlowPause
	FlowStepIn
	FlowStepOut
	FlowStepToPC
	FlowStepOver
	FlowStopWhenSB
	FlowStepToOpcode
)

// New creates a new debugger instance.
func New(c *cpu.CPU, p *ppu.PPU) (*Debugger, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	gui.InputEsc = true

	d := Debugger{
		cpu:       c,
		ppu:       p,
		gui:       gui,
		closed:    abool.New(),
		hasModal:  abool.New(),
		flowState: FlowPause,
	}

	d.gui.SetManagerFunc(d.layout)
	if err := d.initKeybinds(); err != nil {
		return nil, err
	}

	return &d, nil
}

// Close cleans up the terminal and should be called when the debugger is done running.
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

	// Update GUI at 15 fps, no slowdown on my if > ~6ms ticker
	ticker := time.NewTicker(66 * time.Millisecond)
	defer ticker.Stop()

loop:
	for true {
		select {
		case <-ticker.C:
			d.gui.Update(d.updateGUI)
		case <-shouldClose:
			d.gui.Update(func(*gocui.Gui) error { return gocui.ErrQuit })
			<-guiClosed // wait for GUI to exit
			break loop
		case <-guiClosed:
			break loop
		}
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
	d.Lock()

	d.updateInstructions()
	d.updateMemory()
	d.updateMessages()
	d.updateIORegisters()
	d.updateMiscCounters()

	d.Unlock()

	d.flowControl()
}

func (d *Debugger) flowControl() {
	switch atomic.LoadInt32(&d.flowState) {
	case FlowStepIn:
		fallthrough
	case FlowStepOut:
		if d.callDepth == d.requestedDepth {
			atomic.StoreInt32(&d.flowState, FlowPause)
		}
	case FlowStepOver:
		atomic.StoreInt32(&d.flowState, FlowPause)
	case FlowStopWhenSB:
		if d.cpu.Mem[cpu.IOSB] == d.stopWhenSB {
			atomic.StoreInt32(&d.flowState, FlowPause)
		}
	case FlowStepToPC:
		if d.cpu.PC == d.stepToPC {
			atomic.StoreInt32(&d.flowState, FlowPause)
		}
	case FlowStepToOpcode:
		if d.cpu.LastOpcode == d.stepToOpcode && !d.cpu.LastOpcodeWasCB {
			atomic.StoreInt32(&d.flowState, FlowPause)
		}
	}

	for atomic.LoadInt32(&d.flowState) == FlowPause {
		time.Sleep(50 * time.Millisecond)
	}
}

func (d *Debugger) updateIORegisters() {
	d.ioIF = d.cpu.FetchIF()
	d.ioIE = d.cpu.FetchIE()
	d.ioIME = d.cpu.InterruptMaster
	d.ioDIV = d.cpu.FetchIO(cpu.IODIV)
	d.ioInternalDIV = d.cpu.InternalDIV
	d.ioTMA = d.cpu.FetchIO(cpu.IOTMA)
	d.ioTAC = d.cpu.FetchIO(cpu.IOTAC)
	d.ioTIMA = d.cpu.FetchIO(cpu.IOTIMA)
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

	d.updateCallDepth()
}

func (d *Debugger) updateCallDepth() {
	if !d.cpu.LastOpcodeWasCB {
		switch d.cpu.LastOpcode {
		// RST
		case 0xC7:
			fallthrough
		case 0xCF:
			fallthrough
		case 0xD7:
			fallthrough
		case 0xDF:
			fallthrough
		case 0xE7:
			fallthrough
		case 0xEF:
			fallthrough
		case 0xF7:
			fallthrough
		case 0xFF:
			fallthrough
		// CALL
		case 0xCD:
			fallthrough
		case 0xC4:
			fallthrough
		case 0xCC:
			fallthrough
		case 0xD4:
			fallthrough
		case 0xDC:
			d.callDepth++
		case 0xC9:
			fallthrough
		case 0xD9:
			fallthrough
		case 0xC0:
			fallthrough
		case 0xC8:
			fallthrough
		case 0xD0:
			fallthrough
		case 0xD8:
			d.callDepth--
		}
	}

	// Interrupts count as one depth
	if d.cpu.Mem[cpu.IODisableBootROM] == 0x01 {
		switch d.cpu.PC {
		case 0x0040:
			fallthrough
		case 0x0048:
			fallthrough
		case 0x0050:
			fallthrough
		case 0x0058:
			fallthrough
		case 0x0060:
			fallthrough
		case 0x0080:
			d.callDepth++
		}
	}
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

func (d *Debugger) updateMiscCounters() {
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

func (d *Debugger) updateMemory() {
	if d.cpu.MemIOBuffer.Len() <= 0 {
		return
	}

	if (d.cpu.MemIOBuffer.Len() % memBufferStride) != 0 {
		panic("garbage in cpu.MemIOBuffer")
	}

	buf := make([]byte, 6)
	for d.cpu.MemIOBuffer.Len() > 0 {
		n, err := d.cpu.MemIOBuffer.Read(buf)
		if err != nil {
			panic(err)
		}
		if n != memBufferStride {
			panic("n != memBufferStride")
		}

		for i, v := range buf {
			d.memBuffer[d.curMemBufferWriteIndex+i] = v
		}

		d.curMemBufferWriteIndex = (d.curMemBufferWriteIndex + memBufferStride) % len(d.memBuffer)
	}

	d.cpu.MemIOBuffer.Reset()
}
