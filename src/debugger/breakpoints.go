package debugger

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/types"
	"log"
)

type breakpointDebugger struct {
	emulator           *nes.Nes
	panel              *draggablePanel
	breakpointAddPanel *breakpointAdd

	dasmScroll int
	dasmActive int

	breakpoints       [4]uint16
	breakpointsCount  uint8
	breakpointEnabled bool
}

const breakpointDebuggerWidth = 500

func NewBreakpointDebugger(emulator *nes.Nes) *breakpointDebugger {
	return &breakpointDebugger{
		emulator: emulator,
		panel: NewDraggablePanel(
			"Debugger · Breakpoints",
			raylib.Vector2{300, 350},
			breakpointDebuggerWidth,
			400,
		),
		breakpointAddPanel: nil,
		breakpointsCount:   0,
	}
}

func (dbg *breakpointDebugger) Toggle() {
	dbg.panel.SetEnabled(!dbg.panel.enabled)
}

func (dbg *breakpointDebugger) Draw() {
	if !dbg.panel.Draw() {
		return
	}
	padding := float32(5)
	anchor := raylib.Vector2{dbg.panel.position.X + padding, dbg.panel.position.Y + 30}

	dbg.drawDisassembler(anchor)
	dbg.breakPointControls(anchor)
}

func (dbg *breakpointDebugger) drawDisassembler(anchor raylib.Vector2) {
	raylib.GuiLabel(
		raylib.Rectangle{anchor.X, anchor.Y, 200, 20},
		"Disassembler",
	)

	disassembled := dbg.emulator.Debugger().SortedDisassembled()
	pc := dbg.emulator.Debugger().ProgramCounter()

	// todo optimization: extract only the asm lines required, based on current dbg.dasmScroll
	var disasm []string
	disasm = make([]string, 0, len(disassembled))
	line := 0
	active := 0
	for _, asm := range disassembled {
		if asm.Address == pc {
			active = line
			//log.Printf("Current pc found: %x (idx %d)", pc, line)
		}
		disasm = append(disasm, asm.Asm)
		line++
	}

	_, _, scroll := raylib.GuiListViewEx(
		raylib.Rectangle{anchor.X, anchor.Y + 20, 200, 300},
		disasm,
		len(disasm),
		0,
		active,
		active,
	)
	dbg.dasmScroll = scroll
}

func (dbg *breakpointDebugger) breakPointControls(windowAnchor raylib.Vector2) {
	padding := float32(5)
	anchor := raylib.Vector2{windowAnchor.X + 200 + padding, windowAnchor.Y}
	width := float32(290)
	controlsWidth := 290 - padding*2

	raylib.GuiGroupBox(
		raylib.Rectangle{anchor.X, anchor.Y, width, 300},
		"Breakpoints",
	)
	y := dbg.panel.registerStackedControl(anchor.Y+20, padding)

	addBreakPointClicked := raylib.GuiButton(
		raylib.Rectangle{anchor.X, y, controlsWidth, 20},
		"Add breakpoint",
	)
	y = dbg.panel.registerStackedControl(20, padding)

	// List of breakpoints
	var breakpoints []string
	for i := uint8(0); i < dbg.breakpointsCount; i++ {
		breakpoints = append(breakpoints, fmt.Sprintf("0x%X", dbg.breakpoints[i]))
	}
	breakpointListHeight := float32(20 * 4)
	raylib.GuiListViewEx(
		raylib.Rectangle{anchor.X, y, controlsWidth, breakpointListHeight},
		breakpoints,
		len(breakpoints),
		0,
		0,
		-1,
	)
	y = dbg.panel.registerStackedControl(breakpointListHeight, padding)

	// Emulator control
	dbg.breakpointEnabled = raylib.GuiCheckBox(
		raylib.Rectangle{anchor.X, y, 20, 20},
		"Listen BP",
		dbg.breakpointEnabled,
	)
	y = dbg.panel.registerStackedControl(20, padding)

	stepClicked := raylib.GuiButton(
		raylib.Rectangle{anchor.X, y, 100, 20},
		//raylib.GuiIconText(raylib.)
		"Step",
	)

	if stepClicked {
		log.Println("click on step")
		dbg.emulator.Debugger().RunOneCPUOperationAndPause()
	}

	if dbg.breakpointsCount < 4 {
		if addBreakPointClicked {
			dbg.showBreakpointAdd()
		}
		dbg.updateBreakpointAdd()
	}
}

func (dbg *breakpointDebugger) showBreakpointAdd() {
	if dbg.breakpointAddPanel == nil {
		dbg.breakpointAddPanel = NewBreakpointAddPanel(dbg.onAddBreakpoint)
	}
	dbg.breakpointAddPanel.Open(dbg.panel.position.X+10, dbg.panel.position.Y+30)
}

func (dbg *breakpointDebugger) updateBreakpointAdd() {
	if dbg.breakpointAddPanel == nil {
		return
	}

	dbg.breakpointAddPanel.Draw()
}

func (dbg *breakpointDebugger) onAddBreakpoint(address uint16) {
	fmt.Printf("Breakpoint created: %X\n", address)
	dbg.breakpoints[dbg.breakpointsCount] = address
	dbg.breakpointsCount++
	dbg.emulator.Debugger().AddBreakPoint(types.Address(address))
}
