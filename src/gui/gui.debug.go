package gui

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/types"
)

const DEBUG_X_OFFSET = 380

func colorFlag(flag bool) raylib.Color {
	if flag {
		return raylib.Green
	}

	return raylib.RayWhite
}

func DrawDebug(console nes.Nes) {
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
	drawASM(console)
	drawCHR(console, font)
}

func drawASM(console nes.Nes) {
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

func drawCHR(console nes.Nes, font *raylib.Font) {
	x := DEBUG_X_OFFSET
	y := 40 + 15*20 + 20
	posX := 0
	posY := 0
	scale := 3

	// Draw defined palettes (8)
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

	// CHR Left Container

	y = 40 + 15*20 + 50
	raylib.DrawRectangle(x, y, (16*8)*scale+10, (16*8)*scale+10, raylib.RayWhite)

	decodedPatternTable := console.Debugger().PatternTable(0)
	for i := 0; i < 128*128; i++ {
		pixelValue := decodedPatternTable[i]

		color := pixelColor2RaylibColor(pixelValue.Color)
		posX = pixelValue.X + (pixelValue.X * (scale - 1)) + DEBUG_X_OFFSET + 5
		posY = pixelValue.Y + (pixelValue.Y * (scale - 1)) + y + 5
		raylib.DrawRectangle(posX, posY, scale, scale, color)
	}
	//fmt.Printf("%d - %d\n", posX, posY)
	// CHR Right Container
	//r.DrawRectangle(x, y, 16*8, 16*8, r.RayWhite)
}

func pixelColor2RaylibColor(pixelColor types.Color) raylib.Color {
	return raylib.NewColor(
		pixelColor.R,
		pixelColor.G,
		pixelColor.B,
		255,
	)
}
