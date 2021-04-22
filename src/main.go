package main

import (
	"fmt"
	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

func main() {
	// Init systems
	r.InitWindow(800, 800, "NES golang")
	graphics.InitDrawer()

	fmt.Printf("Nes Emulator\n")
	console := nes.CreateNes()
	//gamePak := nes.ReadRom("./roms/nestest/nestest.nes")
	gamePak := nes.ReadRom("./roms/Donkey Kong (World) (Rev A).nes")
	console.InsertCartridge(&gamePak)

	//
	console.Start()

	for !r.WindowShouldClose() {
		// Update emulator

		// Draw
		draw(console)
	}
	r.CloseWindow()
}

func draw(console nes.Nes) {
	r.BeginDrawing()
	r.ClearBackground(r.Black)

	drawEmulation()
	drawDebug(console)

	r.EndDrawing()
}

func drawEmulation() {
	// TODO
}

func drawDebug(console nes.Nes) {
	x := 380
	y := 10
	textColor := r.RayWhite
	font := r.LoadFont("./assets/Pixel_NES.otf")
	r.SetTextureFilter(font.Texture, r.FilterPoint)

	// Status Register
	graphics.DrawText("STATUS:", x, y, textColor)
	graphics.DrawText("N", x+70, y, textColor)

	// Program counter
	msg := fmt.Sprintf("PC: 0x%X", console.Debugger().ProgramCounter())
	graphics.DrawText(msg, x, y+15, textColor)

	//registers := fmt.Sprintf("A:0x%0X X:0x%X Y:0x%X P:0x%X SP:0x%X", 0, 0, 0, 0, 0)
	//position := r.Vector2{X: 380, Y: 20}
	//r.DrawTextEx(*font, registers, position, fontSize, 0, r.RayWhite)

	drawASM(console)
	drawCHR(font)
}

func drawASM(console nes.Nes) {
	textColor := r.RayWhite
	yOffset := 40
	ySeparation := 15
	disassembled := console.Debugger().Disassembled()

	for i := 0; i < 20; i++ {
		currentAddress := console.Debugger().ProgramCounter() - 10 + defs.Address(i)
		//0xc000 + defs.Address(i)
		if currentAddress == console.Debugger().ProgramCounter() {
			textColor = r.GopherBlue
		} else {
			textColor = r.White
		}

		code := disassembled[currentAddress]
		graphics.DrawText(code, 380, yOffset+(i*ySeparation), textColor)
	}
}

func drawCHR(font *r.Font) {

}
