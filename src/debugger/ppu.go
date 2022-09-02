package debugger

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/raulferras/nes-golang/src/nes/ppu"
)

const ppuPanelWidth = 350

type PPUDebugger struct {
	panel *draggablePanel
	ppu   *ppu.P2c02
}

func NewPPUDebugger(ppu *ppu.P2c02) *PPUDebugger {
	return &PPUDebugger{
		panel: NewDraggablePanel(
			"PPU Registers",
			rl.Vector2{X: 300},
			ppuPanelWidth,
			450,
		),
		ppu: ppu,
	}
}

func (dbg *PPUDebugger) Toggle() {
	dbg.panel.SetEnabled(!dbg.panel.enabled)
}

func (dbg *PPUDebugger) Draw() {
	if !dbg.panel.Draw() {
		return
	}
	padding := float32(5)
	fullWidth := ppuPanelWidth - (padding * 2)

	y := dbg.panel.position.Y + 30 + padding
	dbg.ppuControlGroup(fullWidth, dbg.panel.position.X+padding, y)

	y += 64 + padding*2
	dbg.ppuStatusGroup(fullWidth, dbg.panel.position.X+padding, y)

	y += 32 + padding*2
	dbg.ppuMaskGroup(fullWidth, dbg.panel.position.X+padding, y)

	y += 64 + padding*2
	dbg.loopyRegister(fullWidth, dbg.panel.position.X+padding, y, dbg.ppu.VRam(), "V")

	y += 64 + padding*2
	dbg.loopyRegister(fullWidth, dbg.panel.position.X+padding, y, dbg.ppu.TRam(), "T")

	y += 64 + padding*2
	dbg.renderingInfo(fullWidth, dbg.panel.position.X+padding, y)
}

func (dbg *PPUDebugger) ppuControlGroup(fullWidth float32, x float32, y float32) {
	/*
		ntXEnabled := false
		ntYEnabled := false
		incrementModeEnabled := false
		spPatternEnabled := false
		bgPatternEnabled := false
		spriteSizeEnabled := false
		masterSlaveEnabled := false
		generateNMIEnabled := false
		if dbg.ppu.PpuControl.NameTableX == 1 {
			ntXEnabled = true
		}
		if dbg.ppu.PpuControl.NameTableY == 1 {
			ntYEnabled = true
		}
		if dbg.ppu.PpuControl.IncrementMode == 1 {
			incrementModeEnabled = true
		}
		if dbg.ppu.PpuControl.SpritePatternTableAddress == 1 {
			spPatternEnabled = true
		}
		if dbg.ppu.PpuControl.BackgroundPatternTableAddress == 1 {
			bgPatternEnabled = true
		}
		if dbg.ppu.PpuControl.SpriteSize == 1 {
			spriteSizeEnabled = true
		}
		if dbg.ppu.PpuControl.MasterSlaveSelect == 1 {
			masterSlaveEnabled = true
		}
		if dbg.ppu.PpuControl.GenerateNMIAtVBlank {
			generateNMIEnabled = true
		}

		anchor := rl.Vector2{x, y}
		raygui.GroupBox(rl.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 64}, fmt.Sprintf("PPUControl: 0x%0X", dbg.ppu.PpuControl.Value()))

		raygui.CheckBox(rl.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12}, "nt X", ntXEnabled)
		raygui.CheckBox(rl.Rectangle{anchor.X + 10, anchor.Y + 24, 12, 12}, "nt Y", ntYEnabled)

		raygui.CheckBox(rl.Rectangle{anchor.X + 60, anchor.Y + 10, 12, 12}, "incrementMode", incrementModeEnabled)
		raygui.CheckBox(rl.Rectangle{anchor.X + 60, anchor.Y + 24, 12, 12}, "sp Pattern Table", spPatternEnabled)
		raygui.CheckBox(rl.Rectangle{anchor.X + 60, anchor.Y + 38, 12, 12}, "bg Pattern Table", bgPatternEnabled)

		raygui.CheckBox(rl.Rectangle{anchor.X + 180, anchor.Y + 10, 12, 12}, "spriteSize 8x16", spriteSizeEnabled)
		raygui.CheckBox(rl.Rectangle{anchor.X + 180, anchor.Y + 10 + 14, 12, 12}, "master/Slave", masterSlaveEnabled)
		raygui.CheckBox(rl.Rectangle{anchor.X + 180, anchor.Y + 10 + 14 + 14, 12, 12}, "generate NMI", generateNMIEnabled)
	*/
}

func (dbg *PPUDebugger) ppuStatusGroup(fullWidth float32, x float32, y float32) {
	/*
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

		anchor := rl.Vector2{x, y}
		raygui.GroupBox(rl.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 32}, fmt.Sprintf("PPUStatus: 0x%0X", dbg.ppu.PpuStatus.Value()))

		raygui.CheckBox(rl.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12}, "sprite overflow", spriteOverflow)
		raygui.CheckBox(rl.Rectangle{anchor.X + 130, anchor.Y + 10, 12, 12}, "sprite 0 hit", sprite0Hit)

		raygui.CheckBox(rl.Rectangle{anchor.X + 220, anchor.Y + 10, 12, 12}, "VBlank", verticalBlankStarted)
	*/
}

func (dbg *PPUDebugger) ppuMaskGroup(fullWidth float32, x float32, y float32) {
	/*
		greyScale := false
		showBgLeftMost := false
		showSpLeftMost := false
		showBG := false
		showSP := false
		emphasizeRed := false
		emphasizeGreen := false
		emphasizeBlue := false
		if dbg.ppu.PpuMask.GreyScale == 1 {
			greyScale = true
		}
		if dbg.ppu.PpuMask.ShowBackgroundLeftMost == 1 {
			showBgLeftMost = true
		}
		if dbg.ppu.PpuMask.ShowSpritesLeftMost == 1 {
			showSpLeftMost = true
		}
		if dbg.ppu.PpuMask.ShowBackground == 1 {
			showBG = true
		}
		if dbg.ppu.PpuMask.ShowSprites == 1 {
			showSP = true
		}
		if dbg.ppu.PpuMask.EmphasizeRed == 1 {
			emphasizeRed = true
		}
		if dbg.ppu.PpuMask.EmphasizeGreen == 1 {
			emphasizeGreen = true
		}
		if dbg.ppu.PpuMask.EmphasizeBlue == 1 {
			emphasizeBlue = true
		}

		anchor := rl.Vector2{x, y}
		raygui.GroupBox(rl.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 64}, fmt.Sprintf("PPUMask: 0x%0X", dbg.ppu.PpuMask.Value()))

		raygui.CheckBox(rl.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12}, "Grey", greyScale)
		raygui.CheckBox(rl.Rectangle{anchor.X + 10, anchor.Y + 24, 12, 12}, "BG Left Most", showBgLeftMost)

		raygui.CheckBox(rl.Rectangle{anchor.X + 100, anchor.Y + 10, 12, 12}, "SP Left Most", showSpLeftMost)
		raygui.CheckBox(rl.Rectangle{anchor.X + 100, anchor.Y + 24, 12, 12}, "Show Background", showBG)
		raygui.CheckBox(rl.Rectangle{anchor.X + 100, anchor.Y + 38, 12, 12}, "Show Sprites", showSP)

		raygui.CheckBox(rl.Rectangle{anchor.X + 250, anchor.Y + 10, 12, 12}, "Emphasize R", emphasizeRed)
		raygui.CheckBox(rl.Rectangle{anchor.X + 250, anchor.Y + 10 + 14, 12, 12}, "Emphasize G", emphasizeGreen)
		raygui.CheckBox(rl.Rectangle{anchor.X + 250, anchor.Y + 10 + 14 + 14, 12, 12}, "Emphasize B", emphasizeBlue)
	*/
}

func (dbg *PPUDebugger) loopyRegister(fullWidth float32, x float32, y float32, register ppu.LoopyRegister, title string) {
	/*
		anchor := rl.Vector2{x, y}

		raygui.GroupBox(
			rl.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 64},
			fmt.Sprintf("%s: 0x%0X", title, register.Value()),
		)

		raygui.Label(
			rl.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12},
			fmt.Sprintf("Coarse X: 0x%0X (%d)", register.CoarseX(), register.CoarseX()),
		)
		raygui.Label(
			rl.Rectangle{anchor.X + 10, anchor.Y + 24, 12, 12},
			fmt.Sprintf("Coarse Y: 0x%0X (%d)", register.CoarseY(), register.CoarseY()),
		)

		raygui.Label(
			rl.Rectangle{anchor.X + 130, anchor.Y + 10, 12, 12},
			fmt.Sprintf("NX: %d", register.NameTableX()),
		)
		raygui.Label(
			rl.Rectangle{anchor.X + 130, anchor.Y + 24, 12, 12},
			fmt.Sprintf("NY: %d", register.NameTableY()),
		)

		raygui.Label(
			rl.Rectangle{anchor.X + 200, anchor.Y + 10, 12, 12},
			fmt.Sprintf("Fine Y: %d", register.FineY()),
		)
		raygui.Label(
			rl.Rectangle{anchor.X + 200, anchor.Y + 24, 12, 12},
			fmt.Sprintf("Fine X: %d", dbg.ppu.FineX()),
		)
	*/
}
func (dbg *PPUDebugger) renderingInfo(fullWidth float32, x float32, y float32) {
	anchor := rl.Vector2{x, y}
	//raygui.GroupBox(rl.Rectangle{anchor.X + 0, anchor.Y + 0, fullWidth, 64}, "Rendering")

	raygui.Label(
		rl.Rectangle{anchor.X + 10, anchor.Y + 10, 12, 12},
		fmt.Sprintf("Scanline: %d", dbg.ppu.Scanline()),
	)
	raygui.Label(
		rl.Rectangle{anchor.X + 10, anchor.Y + 24, 12, 12},
		fmt.Sprintf("Render Cycle: %d", dbg.ppu.RenderCycle()),
	)
}
