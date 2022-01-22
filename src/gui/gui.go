package gui

import (
	"fmt"
	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"math/rand"
	"time"
)

var font *r.Font
var cpuAdvance bool

func Run() {
	rand.Seed(time.Now().UnixNano())

	// Init Window System
	r.InitWindow(800, 800, "NES golang")
	r.SetTraceLogLevel(r.LogWarning)
	//r.SetTargetFPS(60)
	font = r.LoadFont("./assets/Pixel_NES.otf")

	graphics.InitDrawer()

	fmt.Printf("Nes Emulator\n")
	path := "./assets/roms/nestest/nestest.nes"
	//path := "./assets/roms/full_palette/full_palette.nes"
	//path := "./assets/roms/snake.nes"
	//path := "./assets/roms/Pac-Man (USA) (Namco).nes"
	//path := "./assets/roms/Donkey Kong (World) (Rev A).nes"
	//path := "./assets/roms/Super Mario Bros. (World).nes"
	//path = "./assets/roms/Mega Man 2 (Europe).nes"
	gamePak := gamePak.CreateGamePakFromROMFile(path)

	printRomInfo(&gamePak)

	console := nes.CreateNes(
		&gamePak,
		nes.CreateNesDebugger("./var/run.log", true),
	)
	//console.InsertGamePak(&gamePak)

	//
	console.Start()
	cpuAdvance = true
	_timestamp := r.GetTime()
	for !r.WindowShouldClose() {
		timestamp := r.GetTime()
		dt := timestamp - _timestamp
		_timestamp = timestamp
		if dt > 1 {
			dt = 0
		}
		//fmt.Printf("%f sec\n", dt)
		//if r.IsKeyPressed(r.KeySpace) {
		//	cpuAdvance = true
		//}

		// Update emulator
		if cpuAdvance {
			//console.Tick()
			console.TickForTime(dt)
		}

		// Draw
		draw(console)
		//if i > 4000 {
		//	break
		//}

		cpuAdvance = true
	}
	r.CloseWindow()

	console.Stop()
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

func draw(console nes.Nes) {
	r.BeginDrawing()
	r.ClearBackground(r.Black)

	drawEmulation(console)
	drawBackgroundTileIDs(console)
	DrawDebug(console)
	r.EndDrawing()
}

func drawEmulation(console nes.Nes) {
	frame := console.Frame()

	padding := 20
	paddingY := 20
	r.DrawRectangle(padding, paddingY, types.WIDTH, types.HEIGHT, r.DarkBrown)
	for i := 0; i < types.WIDTH*types.HEIGHT; i++ {
		pixel := frame.Pixels[i]
		color := types.PixelColor2RaylibColor(pixel)
		x := i % types.WIDTH
		y := i / types.WIDTH

		r.DrawPixel(padding+x, paddingY+y, color)
	}
}

func drawBackgroundTileIDs(console nes.Nes) {
	padding := 20
	paddingY := 100
	// Debug background tiles IDS
	offsetY := paddingY + types.HEIGHT + 10
	framePattern := console.FramePattern()
	tilesWidth := 32
	//tilesHeight := 30
	for i := 0; i < 0x3C0; i++ {
		x := i % tilesWidth * 17
		y := (i / tilesWidth) * 17
		r.DrawText(
			fmt.Sprintf("%X", framePattern[i]),
			padding+x,
			offsetY+y,
			8,
			r.Violet,
		)
	}
}
