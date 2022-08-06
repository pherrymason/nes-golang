package debugger

import "github.com/lachee/raylib-goplus/raylib"

type breakpointDebugger struct {
	panel *draggablePanel
}

const width = 300

func NewBreakpointDebugger() *breakpointDebugger {
	return &breakpointDebugger{
		panel: NewDraggablePanel(
			"Debugger · Breakpoints",
			raylib.Vector2{300, 350},
			width,
			400,
		),
	}
}

func (bp *breakpointDebugger) Toggle() {
	bp.panel.SetEnabled(!bp.panel.enabled)
}

func (dbg *breakpointDebugger) Draw() {
	if !dbg.panel.Draw() {
		return
	}
}
