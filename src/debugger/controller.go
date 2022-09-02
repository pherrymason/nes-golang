package debugger

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/raulferras/nes-golang/src/audio"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image/color"
)

const DEBUG_X_OFFSET = 300

type GuiDebugger struct {
	chrPaletteSelector    uint8
	overlayTileIdx        bool
	overlayAttributeTable bool
	font                  *rl.Font

	emulator           *nes.Nes
	ppuDebugger        *PPUDebugger
	breakpointDebugger *breakpointDebugger
	audioDebugger      *audioDebugger
}

type Panel interface {
	Draw()
}

func NewDebugger(emulator *nes.Nes, audio *audio.Audio) GuiDebugger {
	font := rl.LoadFont("./assets/Pixel_NES.otf")
	//rl.GuiLoadStyle("./assets/style.rgs")

	return GuiDebugger{
		chrPaletteSelector:    0,
		overlayTileIdx:        false,
		overlayAttributeTable: false,
		font:                  &font,
		emulator:              emulator,
		ppuDebugger:           NewPPUDebugger(emulator.PPU()),
		breakpointDebugger:    NewBreakpointDebugger(emulator),
		audioDebugger:         NewAudioDebugger(audio),
	}
}

func (dbg *GuiDebugger) Close() {
	defer rl.UnloadFont(*dbg.font)
}

func (dbg *GuiDebugger) Tick() {
	dbg.listenKeyboard()
	rl.DrawFPS(0, 0)

	dbg.ppuDebugger.Draw()
	dbg.breakpointDebugger.Draw()
	//dbg.DrawDebugger(dbg.emulator)
	dbg.audioDebugger.Draw()
}

func (dbg *GuiDebugger) DrawDebugger(emulator *nes.Nes) {
	x := DEBUG_X_OFFSET
	y := 10

	dbg.listenKeyboard()

	textColor := rl.RayWhite
	fontSize := float32(16)

	rl.SetTextureFilter(dbg.font.Texture, rl.FilterPoint)

	// Status Register
	graphics.DrawText("STATUS:", x, y, textColor, fontSize)

	graphics.DrawText("N", x+70, y, colorFlag(emulator.Debugger().N()), fontSize)
	graphics.DrawText("O", x+90, y, colorFlag(emulator.Debugger().O()), fontSize)
	graphics.DrawText("-", x+110, y, rl.RayWhite, fontSize)
	graphics.DrawText("B", x+130, y, colorFlag(emulator.Debugger().B()), fontSize)
	graphics.DrawText("D", x+150, y, colorFlag(emulator.Debugger().D()), fontSize)
	graphics.DrawText("I", x+170, y, colorFlag(emulator.Debugger().I()), fontSize)
	graphics.DrawText("Z", x+190, y, colorFlag(emulator.Debugger().Z()), fontSize)
	graphics.DrawText("C", x+210, y, colorFlag(emulator.Debugger().C()), fontSize)

	// Program counter
	msg := fmt.Sprintf("PC: 0x%X", emulator.Debugger().ProgramCounter())
	graphics.DrawText(msg, x, y+15, textColor, fontSize)

	// A, X, Y Registers
	msg = fmt.Sprintf("A:0x%X", emulator.Debugger().ARegister())
	graphics.DrawText(msg, x+130, y+15, textColor, fontSize)

	msg = fmt.Sprintf("X:0x%X", emulator.Debugger().XRegister())
	graphics.DrawText(msg, x+200, y+15, textColor, fontSize)

	msg = fmt.Sprintf("Y: 0x%X", emulator.Debugger().YRegister())
	graphics.DrawText(msg, x+270, y+15, textColor, fontSize)

	//registers := fmt.Sprintf("A:0x%0X X:0x%X Y:0x%X P:0x%X SP:0x%X", 0, 0, 0, 0, 0)
	//position := rl.Vector2{X: 380, Y: 20}
	//rl.DrawTextEx(*font, registers, position, 16, 0, rl.RayWhite)
	//
	//scale := 3
	drawASM(emulator)
	//drawPalettes(emulator, scale, DEBUG_X_OFFSET, 40+15*20, debuggerGUI)
	dbg.drawCHR(emulator, 2, DEBUG_X_OFFSET, 40+15*20+50)

	if emulator.Debugger().DebugPPU {
		drawPPUDebugger(emulator)
	}

	//drawObjectAttributeEntries(emulator)
}

func (dbg *GuiDebugger) listenKeyboard() {
	if rl.IsKeyPressed(rl.KeyP) {
		//dbg.chrPaletteSelector += 1
		//if dbg.chrPaletteSelector > (8 - 1) {
		//	dbg.chrPaletteSelector = 0
		//}
		dbg.ppuDebugger.Toggle()
	}

	if rl.IsKeyPressed(rl.KeyO) {
		dbg.breakpointDebugger.Toggle()
	}
}

func colorFlag(flag bool) rl.Color {
	if flag {
		return rl.Green
	}

	return rl.RayWhite
}

func drawASM(console *nes.Nes) {
	textColor := rl.RayWhite
	yOffset := 60
	yIteration := 0
	ySeparation := 15
	disassembled := console.Debugger().Disassembled()

	for i := 0; i < 20; i++ {
		currentAddress := console.Debugger().ProgramCounter() - 10 + types.Address(i)
		if currentAddress == console.Debugger().ProgramCounter() {
			textColor = rl.Blue
		} else {
			textColor = rl.White
		}

		code := disassembled[currentAddress]
		if len(code) > 0 {
			graphics.DrawText(code, 380, yOffset+(yIteration*ySeparation), textColor, 16)
			yIteration++
		}
	}
}

func drawPalettes(console *nes.Nes, scale int32, xOffset int32, yOffset int32, debuggerGUI *GuiDebugger) {
	// Draw defined palettes (8)
	x := int32(xOffset)
	y := int32(yOffset)
	posX := int32(0)
	posY := int32(0)
	border := int32(1)
	oneColorWidth := int32(5 * scale)
	paletteWidth := oneColorWidth * 4
	height := int32(2 * scale)
	paddingX := int32(5)
	paddingY := int32(15)

	for i := int32(0); i < 8; i++ {
		m := int32(i % 4)
		posX = x + (m * paletteWidth) + paddingX*m
		posY = y + (height+border*2+paddingY)*(i/4)
		colors := console.Debugger().GetPaletteFromRam(uint8(i))

		rl.DrawRectangle(int32(posX-border), int32(posY-border), int32(oneColorWidth*4+border*2), int32(height+border*2), rl.White)

		drawColorWatch(posX, posY, oneColorWidth, height, colors[0], uint8(i), 0, console)
		drawColorWatch(posX+oneColorWidth, posY, oneColorWidth, height, colors[1], uint8(i), 1, console)
		drawColorWatch(posX+oneColorWidth*2, posY, oneColorWidth, height, colors[2], uint8(i), 2, console)
		drawColorWatch(posX+oneColorWidth*3, posY, oneColorWidth, height, colors[3], uint8(i), 3, console)

		if int32(debuggerGUI.chrPaletteSelector) == i {
			graphics.DrawArrow(int32(posX+oneColorWidth+3), int32(posY-height-2), int32(scale))
		}
	}
}

func drawColorWatch(coordX int32, coordY int32, width int32, height int32, color color.Color, paletteIndex byte, colorIndex byte, console *nes.Nes) {
	rl.DrawText(
		fmt.Sprintf("%0X", console.Debugger().GetPaletteColorFromPaletteRam(paletteIndex, colorIndex)),
		coordX,
		coordY-10,
		10,
		rl.White,
	)
	rl.DrawRectangle(coordX, coordY, width, height, pixelColor2rlColor(color))
}

func (dbg *GuiDebugger) drawCHR(console *nes.Nes, scale int32, xOffset int32, yOffset int32) {
	drawIndexes := false

	if rl.IsKeyDown(rl.KeyZero) {
		drawIndexes = true
	}

	// CHR Left Container
	rl.DrawRectangle(xOffset, yOffset, (16*8)*scale+10, (16*8)*scale+10, rl.RayWhite)

	const BorderWidth = int32(5)
	nextOffsetY := yOffset
	for patternTableIdx := 0; patternTableIdx < 2; patternTableIdx++ {
		decodedPatternTable := console.Debugger().PatternTable(byte(patternTableIdx), dbg.chrPaletteSelector)
		yOffset = nextOffsetY

		for x := 0; x < decodedPatternTable.Bounds().Max.X; x++ {
			for y := 0; y < decodedPatternTable.Bounds().Max.Y; y++ {
				screenX := int32(DEBUG_X_OFFSET + BorderWidth)
				screenX += int32(x) * scale

				screenY := int32(yOffset + BorderWidth)
				screenY += int32(y) * scale

				rl.DrawRectangle(
					screenX,
					screenY,
					scale,
					scale,
					pixelColor2rlColor(decodedPatternTable.At(x, y)),
				)

				nextOffsetY = screenY
			}
		}

		if drawIndexes {
			for i := 0; i < 16*8; i++ {
				screenX := DEBUG_X_OFFSET + BorderWidth +
					types.LinearToXCoordinate(i, 16)*8*scale
				screenY := yOffset + BorderWidth + types.LinearToYCoordinate(i, 16)*8*scale

				rl.DrawRectangle(
					screenX-1,
					screenY-1,
					20,
					10,
					rl.RayWhite,
				)
				rl.DrawText(
					fmt.Sprintf("%d", i),
					screenX,
					screenY,
					10,
					rl.Black,
				)
			}
		}
	}
	//fmt.Printf("%d - %d\n", posX, posY)
	// CHR Right Container
	//r.DrawRectangle(x, y, 16*8, 16*8, r.RayWhite)
}

func drawPPUDebugger(console *nes.Nes) {
	drawBackgroundTileIDs(console, 600, 10)
}

func drawObjectAttributeEntries(console *nes.Nes) {
	for i := 0; i < 20; i++ {
		oae := console.Debugger().OAM(byte(i))
		graphics.DrawText(
			fmt.Sprintf("[%d] x:%d y:%d tileId: %x",
				i,
				oae[3],
				oae[0],
				oae[1],
			),
			50,
			300+i*16,
			rl.White,
			10,
		)
	}
}

func drawBackgroundTileIDs(console *nes.Nes, xOffset int, yOffset int) {
	padding := 10 + xOffset
	//paddingY := 100 + yOffset
	// Debug background tiles IDS
	//offsetY := paddingY + types.SCREEN_HEIGHT + 10
	offsetY := yOffset
	framePattern := console.FramePattern()
	tilesWidth := 32
	//tilesHeight := 30
	for i := 0; i < tilesWidth*32; i++ {
		x := i % tilesWidth * 17
		y := (i / tilesWidth) * 17
		graphics.DrawText(
			fmt.Sprintf("%X", framePattern[i]),
			padding+x,
			offsetY+y,
			rl.White,
			10,
		)
	}
}

func pixelColor2rlColor(pixelColor color.Color) rl.Color {
	r, g, b, a := pixelColor.RGBA()

	return rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
}
