package main

import (
	"fmt"
	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"math/rand"
	"time"
)

var font *r.Font

func main() {
	rand.Seed(time.Now().UnixNano())

	// Init Window System
	r.InitWindow(800, 800, "NES golang")
	r.SetTraceLogLevel(r.LogWarning)
	r.SetTargetFPS(60)
	font = r.LoadFont("./assets/Pixel_NES.otf")

	graphics.InitDrawer()

	fmt.Printf("Nes Emulator\n")
	//gamePak := nes.CreateGamePakFromROMFile("./roms/nestest/nestest.nes")
	gamePak := nes.CreateGamePakFromROMFile("./roms/Donkey Kong (World) (Rev A).nes")

	printRomInfo(&gamePak)

	console := nes.CreateNes(&gamePak, nes.NesDebugger{})
	//console.InsertGamePak(&gamePak)

	//
	console.Start()

	for !r.WindowShouldClose() {
		// Update emulator

		// Draw
		draw(console)
	}
	r.CloseWindow()
}

func printRomInfo(gamePak *nes.GamePak) {
	inesHeader := gamePak.Header()

	if inesHeader.HasTrainer() {
		fmt.Println("Rom has trainer")
	} else {
		fmt.Println("Rom has no trainer")
	}

	fmt.Println("PRG:", inesHeader.ProgramSize(), "x 16KB Banks")
	fmt.Println("CHR:", inesHeader.CHRSize(), "x 8KB Banks")
	fmt.Println("Mapper:", inesHeader.MapperNumber())
	fmt.Println("Tv System:", inesHeader.TvSystem())
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
	yIteration := 0
	ySeparation := 15
	disassembled := console.Debugger().Disassembled()

	for i := 0; i < 20; i++ {
		currentAddress := console.Debugger().ProgramCounter() - 10 + nes.Address(i)
		if currentAddress == console.Debugger().ProgramCounter() {
			textColor = r.GopherBlue
		} else {
			textColor = r.White
		}

		code := disassembled[currentAddress]
		if len(code) > 0 {
			graphics.DrawText(code, 380, yOffset+(yIteration*ySeparation), textColor)
			yIteration++
		}
	}
}

func drawCHR(font *r.Font) {

}
