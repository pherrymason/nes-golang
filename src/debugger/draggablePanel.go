package debugger

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type draggablePanel struct {
	title               string
	enabled             bool
	position            rl.Vector2
	width               float32
	height              float32
	dragWindow          bool
	positionOnStartDrag rl.Vector2
	layoutYPositions    []float32
}

func NewDraggablePanel(title string, position rl.Vector2, width int, height int) *draggablePanel {
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

	panel.layoutYPositions = nil
	panel.updateWindowPosition()
	/*
		shouldClose := raygui.WindowBox(
			rl.Rectangle{
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
	*/
	return true
}

func (panel *draggablePanel) Close() {
	panel.enabled = false
}

func (panel *draggablePanel) updateWindowPosition() {
	mousePosition := rl.GetMousePosition()
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		if rl.CheckCollisionPointRec(mousePosition, panel.statusBarPosition()) {
			panel.dragWindow = true
			panel.positionOnStartDrag = rl.Vector2{
				X: mousePosition.X - panel.position.X,
				Y: mousePosition.Y - panel.position.Y,
			}
		}
	}

	if panel.dragWindow {
		panel.position.X = mousePosition.X - panel.positionOnStartDrag.X
		panel.position.Y = mousePosition.Y - panel.positionOnStartDrag.Y
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			panel.dragWindow = false
		}
	}
}

func (panel *draggablePanel) statusBarPosition() rl.Rectangle {
	return rl.Rectangle{
		X:      panel.position.X,
		Y:      panel.position.Y,
		Width:  panel.width - 20,
		Height: 20,
	}
}

// registerStackedControl registers the height of a gui control rendered, and returns the Y position for the next element
func (panel *draggablePanel) registerStackedControl(height float32, padding float32) float32 {
	panel.layoutYPositions = append(panel.layoutYPositions, height)
	sum := float32(0)
	for i := 0; i < len(panel.layoutYPositions); i++ {
		sum += panel.layoutYPositions[i] + padding
	}
	return sum
}
