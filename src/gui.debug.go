package main

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image/color"
)

const DEBUG_X_OFFSET = 380

func colorFlag(flag bool) raylib.Color {
	if flag {
		return raylib.Green
	}

	return raylib.RayWhite
}

func printRomInfo(cartridge *gamePak.GamePak) {
	inesHeader := cartridge.Header()

	if inesHeader.HasTrainer() {
		fmt.Println("Rom has trainer")
	} else {
		fmt.Println("Rom has no trainer")
	}

	if inesHeader.Mirroring() == gamePak.VerticalMirroring {
		fmt.Println("Vertical Mirroring")
	} else {
		fmt.Println("Horizontal Mirroring")
	}

	fmt.Println("PRG:", inesHeader.ProgramSize(), "x 16KB Banks")
	fmt.Println("CHR:", inesHeader.CHRSize(), "x 8KB Banks")
	fmt.Println("Mapper:", inesHeader.MapperNumber())
	fmt.Println("Tv System:", inesHeader.TvSystem())
}

func DrawDebug(console *nes.Nes) {
	x := DEBUG_X_OFFSET
	y := 10
	raylib.DrawFPS(0, 0)

	textColor := raylib.RayWhite

	raylib.SetTextureFilter(font.Texture, raylib.FilterPoint)

	// Status Register
	graphics.DrawText("STATUS:", x, y, textColor)

	graphics.DrawText("N", x+70, y, colorFlag(console.Debugger().N()))
	graphics.DrawText("O", x+90, y, colorFlag(console.Debugger().O()))
	graphics.DrawText("-", x+110, y, raylib.RayWhite)
	graphics.DrawText("B", x+130, y, colorFlag(console.Debugger().B()))
	graphics.DrawText("D", x+150, y, colorFlag(console.Debugger().D()))
	graphics.DrawText("I", x+170, y, colorFlag(console.Debugger().I()))
	graphics.DrawText("Z", x+190, y, colorFlag(console.Debugger().Z()))
	graphics.DrawText("C", x+210, y, colorFlag(console.Debugger().C()))

	// Program counter
	msg := fmt.Sprintf("PC: 0x%X", console.Debugger().ProgramCounter())
	graphics.DrawText(msg, x, y+15, textColor)

	// A, X, Y Registers
	msg = fmt.Sprintf("A:0x%X", console.Debugger().ARegister())
	graphics.DrawText(msg, x+130, y+15, textColor)

	msg = fmt.Sprintf("X:0x%X", console.Debugger().XRegister())
	graphics.DrawText(msg, x+200, y+15, textColor)

	msg = fmt.Sprintf("Y: 0x%X", console.Debugger().YRegister())
	graphics.DrawText(msg, x+270, y+15, textColor)

	//registers := fmt.Sprintf("A:0x%0X X:0x%X Y:0x%X P:0x%X SP:0x%X", 0, 0, 0, 0, 0)
	//position := raylib.Vector2{X: 380, Y: 20}
	//raylib.DrawTextEx(*font, registers, position, 16, 0, raylib.RayWhite)
	//
	scale := 3
	drawASM(console)
	drawPalettes(console, scale, DEBUG_X_OFFSET, 40+15*20)
	drawCHR(console, 2, DEBUG_X_OFFSET, 40+15*20+50, font)
	//drawBackgroundTileIDs(console, DEBUG_X_OFFSET+356, 0)
}

func drawASM(console *nes.Nes) {
	textColor := raylib.RayWhite
	yOffset := 60
	yIteration := 0
	ySeparation := 15
	disassembled := console.Debugger().Disassembled()

	for i := 0; i < 20; i++ {
		currentAddress := console.Debugger().ProgramCounter() - 10 + types.Address(i)
		if currentAddress == console.Debugger().ProgramCounter() {
			textColor = raylib.GopherBlue
		} else {
			textColor = raylib.White
		}

		code := disassembled[currentAddress]
		if len(code) > 0 {
			graphics.DrawText(code, 380, yOffset+(yIteration*ySeparation), textColor)
			yIteration++
		}
	}
}

func drawPalettes(console *nes.Nes, scale int, xOffset int, yOffset int) {
	// Draw defined palettes (8)
	x := xOffset
	y := yOffset
	posX := 0
	posY := 0
	selectedPalette := 5
	for i := 0; i < 8; i++ {
		width := 5 * scale
		height := 2 * scale
		posX = x + (i*width*3 + (i * 5))
		posY = y
		colors := console.Debugger().GetPaletteFromRam(uint8(i))

		raylib.DrawRectangle(posX, posY, width, height, pixelColor2RaylibColor(colors[0]))
		raylib.DrawRectangle(posX+width, posY, width, height, pixelColor2RaylibColor(colors[1]))
		raylib.DrawRectangle(posX+width*2, posY, width, height, pixelColor2RaylibColor(colors[2]))

		if selectedPalette == i {
			graphics.DrawArrow(posX+width+3, posY-height-2, scale)
		}
	}
}

func drawCHR(console *nes.Nes, scale int, xOffset int, yOffset int, font *raylib.Font) {
	drawIndexes := false

	if raylib.IsKeyDown(raylib.KeyZero) {
		drawIndexes = true
	}

	// CHR Left Container
	raylib.DrawRectangle(xOffset, yOffset, (16*8)*scale+10, (16*8)*scale+10, raylib.RayWhite)

	const BorderWidth = 5
	decodedPatternTable := console.Debugger().PatternTable(0)

	for x := 0; x < decodedPatternTable.Bounds().Max.X; x++ {
		for y := 0; y < decodedPatternTable.Bounds().Max.Y; y++ {
			screenX := DEBUG_X_OFFSET + BorderWidth
			screenX += x * scale

			screenY := yOffset + BorderWidth
			screenY += y * scale

			raylib.DrawRectangle(
				screenX,
				screenY,
				scale,
				scale,
				pixelColor2RaylibColor(decodedPatternTable.At(x, y)),
			)
		}
	}

	if drawIndexes {
		for i := 0; i < 16*8; i++ {
			screenX := DEBUG_X_OFFSET + BorderWidth +
				types.LinearToXCoordinate(i, 16)*8*scale
			screenY := yOffset + BorderWidth + types.LinearToYCoordinate(i, 16)*8*scale

			raylib.DrawRectangle(
				screenX-1,
				screenY-1,
				20,
				10,
				raylib.RayWhite,
			)
			raylib.DrawText(
				fmt.Sprintf("%d", i),
				screenX,
				screenY,
				10,
				raylib.Black,
			)
		}
	}
	//fmt.Printf("%d - %d\n", posX, posY)
	// CHR Right Container
	//r.DrawRectangle(x, y, 16*8, 16*8, r.RayWhite)
}

func drawBackgroundTileIDs(console nes.Nes, xOffset int, yOffset int) {
	padding := 20 + xOffset
	//paddingY := 100 + yOffset
	// Debug background tiles IDS
	//offsetY := paddingY + types.SCREEN_HEIGHT + 10
	offsetY := yOffset
	framePattern := console.FramePattern()
	tilesWidth := 32
	//tilesHeight := 30
	for i := 0; i < tilesWidth*30; i++ {
		x := i % tilesWidth * 17
		y := (i / tilesWidth) * 17
		raylib.DrawText(
			fmt.Sprintf("%X", framePattern[i]),
			padding+x,
			offsetY+y,
			8,
			raylib.Violet,
		)
	}
}

func pixelColor2RaylibColor(pixelColor color.Color) raylib.Color {
	r, g, b, a := pixelColor.RGBA()

	return raylib.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
}