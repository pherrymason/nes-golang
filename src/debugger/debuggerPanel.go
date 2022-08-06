package debugger

import "github.com/lachee/raylib-goplus/raylib"

type draggablePanel struct {
	title               string
	enabled             bool
	position            raylib.Vector2
	width               float32
	height              float32
	dragWindow          bool
	positionOnStartDrag raylib.Vector2
}

func NewDraggablePanel(title string, position raylib.Vector2, width int, height int) *draggablePanel {
	return &draggablePanel{
		title:    title,
		enabled:  false,
		position: position,
		width:    float32(width),
		height:   float32(height),
	}
}

func (panel *draggablePanel) SetEnabled(enabled bool) {
	panel.enabled = enabled
}

// Draw returns true if panel is active
func (panel *draggablePanel) Draw() bool {
	if panel.enabled == false {
		return panel.enabled
	}

	panel.updateWindowPosition()
	shouldClose := raylib.GuiWindowBox(
		raylib.Rectangle{
			X:      panel.position.X,
			Y:      panel.position.Y,
			Width:  panel.width,
			Height: panel.height,
		},
		panel.title,
	)
	if shouldClose {
		panel.Close()
	}

	return true
}

func (panel *draggablePanel) Close() {
	panel.enabled = false
}

func (panel *draggablePanel) updateWindowPosition() {
	mousePosition := raylib.GetMousePosition()
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {

		if raylib.CheckCollisionPointRec(mousePosition, panel.statusBarPosition()) {
			panel.dragWindow = true
			panel.positionOnStartDrag = raylib.Vector2{
				X: mousePosition.X - panel.position.X,
				Y: mousePosition.Y - panel.position.Y,
			}
		}
	}

	if panel.dragWindow {
		panel.position.X = mousePosition.X - panel.positionOnStartDrag.X
		panel.position.Y = mousePosition.Y - panel.positionOnStartDrag.Y
		if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
			panel.dragWindow = false
		}
	}
}

func (panel *draggablePanel) statusBarPosition() raylib.Rectangle {
	return raylib.Rectangle{
		X:      panel.position.X,
		Y:      panel.position.Y,
		Width:  ppuPanelWidth - 20,
		Height: 20,
	}
}
