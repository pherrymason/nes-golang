package debugger

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/nes/ppu"
)

const debuggerWidth = 350

type PPUDebugger struct {
	enabled             bool
	windowRectangle     raylib.Rectangle
	ppu                 *ppu.Ppu2c02
	dragWindow          bool
	positionOnStartDrag raylib.Vector2
}

func NewPPUDebugger(ppu *ppu.Ppu2c02) *PPUDebugger {
	return &PPUDebugger{
		enabled:         false,
		windowRectangle: raylib.Rectangle{X: 300, Width: debuggerWidth, Height: 300},
		dragWindow:      false,
		ppu:             ppu,
	}
}

func (dbg *PPUDebugger) SetEnabled(enabled bool) {
	dbg.enabled = enabled
}

func (dbg *PPUDebugger) Draw() {
	if dbg.enabled == false {
		return
	}

	dbg.updateWindowPosition()
	shouldClose := raylib.GuiWindowBox(dbg.windowRectangle, "PPU Viewer")
	if shouldClose {
		dbg.Close()
	}
	padding := float32(5)
	fullWidth := debuggerWidth - (padding * 2)

	dbg.ppuControlGroup(fullWidth, dbg.windowRectangle.X+padding, dbg.windowRectangle.Y+30+padding)
	dbg.ppuStatusGroup(fullWidth, dbg.windowRectangle.X+padding, dbg.windowRectangle.Y+30+64+padding*3)
}

func (dbg *PPUDebugger) ppuControlGroup(fullWidth float32, x float32, y float32) {
	ntXEnabled := false
	ntYEnabled := false
	incrementModeEnabled := false
	spPatternEnabled := false
	bgPatternEnabled := false
	spriteSizeEnabled := false
	masterSlaveEnabled := false
	generateNMIEnabled := false
	ppuControl := dbg.ppu.PpuControl
	if ppuControl.NameTableX == 1 {
		ntXEnabled = true
	}
	if ppuControl.NameTableY == 1 {
		ntYEnabled = true
	}
	if ppuControl.IncrementMode == 1 {
		incrementModeEnabled = true
	}
	if ppuControl.SpritePatternTableAddress == 1 {
		spPatternEnabled = true
	}
	if ppuControl.BackgroundPatternTableAddress == 1 {
		bgPatternEnabled = true
	}
	if ppuControl.SpriteSize == 1 {
		spriteSizeEnabled = true
	}
	if ppuControl.MasterSlaveSelect == 1 {
		masterSlaveEnabled = true
	}
	if ppuControl.GenerateNMIAtVBlank {
		generateNMIEnabled = true
	}

	anchor := raylib.Vector2{x, y}
	raylib.GuiGroupBox(raylib.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 64}, fmt.Sprintf("PPUControl: 0x%0X", ppuControl.Value()))

	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12}, "nt X", ntXEnabled)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 10, anchor.Y + 24, 12, 12}, "nt Y", ntYEnabled)

	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 60, anchor.Y + 10, 12, 12}, "incrementMode", incrementModeEnabled)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 60, anchor.Y + 24, 12, 12}, "sp Pattern Table", spPatternEnabled)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 60, anchor.Y + 38, 12, 12}, "bg Pattern Table", bgPatternEnabled)

	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 180, anchor.Y + 10, 12, 12}, "spriteSize 8x16", spriteSizeEnabled)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 180, anchor.Y + 10 + 14, 12, 12}, "master/Slave", masterSlaveEnabled)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 180, anchor.Y + 10 + 14 + 14, 12, 12}, "generate NMI", generateNMIEnabled)
}

func (dbg *PPUDebugger) ppuStatusGroup(fullWidth float32, x float32, y float32) {
	spriteOverflow := false
	sprite0Hit := false
	verticalBlankStarted := false
	if dbg.ppu.PpuStatus.SpriteOverflow == 1 {
		spriteOverflow = true
	}
	if dbg.ppu.PpuStatus.Sprite0Hit == 1 {
		sprite0Hit = true
	}
	if dbg.ppu.PpuStatus.VerticalBlankStarted {
		verticalBlankStarted = true
	}

	anchor := raylib.Vector2{x, y}
	raylib.GuiGroupBox(raylib.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 32}, fmt.Sprintf("PPUStatus: 0x%0X", dbg.ppu.PpuStatus.Value()))

	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12}, "sprite overflow", spriteOverflow)
	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 130, anchor.Y + 10, 12, 12}, "sprite 0 hit", sprite0Hit)

	raylib.GuiCheckBox(raylib.Rectangle{anchor.X + 220, anchor.Y + 10, 12, 12}, "VBlank", verticalBlankStarted)
}

func (dbg *PPUDebugger) updateWindowPosition() {
	mousePosition := raylib.GetMousePosition()
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {

		if raylib.CheckCollisionPointRec(mousePosition, dbg.statusBarPosition()) {
			fmt.Printf("collision\n")
			dbg.dragWindow = true
			dbg.positionOnStartDrag = raylib.Vector2{
				X: mousePosition.X - dbg.windowRectangle.X,
				Y: mousePosition.Y - dbg.windowRectangle.Y,
			}
			fmt.Printf("drag start: %d %d\n", int(dbg.windowRectangle.X), int(dbg.windowRectangle.Y))
		}
	}

	if dbg.dragWindow {
		dbg.windowRectangle.X = mousePosition.X - dbg.positionOnStartDrag.X
		dbg.windowRectangle.Y = mousePosition.Y - dbg.positionOnStartDrag.Y
		//fmt.Printf("dragging Mouse position: %d %d\n", int(dbg.windowRectangle.X), int(dbg.windowRectangle.Y))
		if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
			fmt.Printf("release\n")
			dbg.dragWindow = false
		}
	}
}

func (dbg *PPUDebugger) statusBarPosition() raylib.Rectangle {
	return raylib.Rectangle{
		X:      dbg.windowRectangle.X,
		Y:      dbg.windowRectangle.Y,
		Width:  debuggerWidth - 20,
		Height: 20,
	}
}

func (dbg *PPUDebugger) Close() {
	dbg.enabled = false
}
