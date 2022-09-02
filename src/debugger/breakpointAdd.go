package debugger

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type breakpointAdd struct {
	panel                   *draggablePanel
	breakPointAddress       uint16
	inputAddress            string
	onAddBreakPointCallback func(addressBreakPoint uint16)
}

func NewBreakpointAddPanel(onOk func(addressBreakPoint uint16)) *breakpointAdd {
	return &breakpointAdd{
		breakPointAddress:       0,
		inputAddress:            "",
		onAddBreakPointCallback: onOk,
		panel: NewDraggablePanel(
			"Add breakpoint",
			rl.Vector2{
				X: 0,
				Y: 0,
			},
			300,
			300,
		),
	}
}

func (dbg *breakpointAdd) Open(x float32, y float32) {
	dbg.panel.position.X = x
	dbg.panel.position.Y = y
	dbg.panel.SetEnabled(true)
}

func (dbg *breakpointAdd) Draw() {
	/*
		if !dbg.panel.Draw() {
			return
		}

		rl.GuiLabel(
			rl.Rectangle{
				X:      dbg.panel.position.X + 5,
				Y:      dbg.panel.position.Y + 30 + 5,
				Height: 20,
			},
			"Address",
		)
		pressed, value := rl.GuiTextBox(
			rl.Rectangle{
				X:      dbg.panel.position.X + 40 + 5 + 5,
				Y:      dbg.panel.position.Y + 30 + 5,
				Width:  70,
				Height: 20,
			},
			dbg.inputAddress,
			5,
			true,
		)
		valueValid := false
		dbg.inputAddress = value
		formattedAddress := fmt.Sprintf("%04s", value)
		if len(formattedAddress) == 4 {
			valueValid = true
			dbg.breakPointAddress = hexit.UnhexUint16Str(formattedAddress)
		}

		if valueValid {
			pressed = rl.GuiButton(
				rl.Rectangle{
					X:      dbg.panel.position.X + 5,
					Y:      dbg.panel.position.Y + dbg.panel.height - 5 - 20,
					Width:  dbg.panel.width - 5 - 5,
					Height: 20,
				},
				"Ok",
			)
		}

		if pressed {
			dbg.panel.Close()
			dbg.onAddBreakPointCallback(dbg.breakPointAddress)
		}
	*/
}
