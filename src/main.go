package main

import (
	"flag"
	"fmt"
	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/raulferras/nes-golang/src/utils"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

var font *r.Font
var cpuAdvance bool

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var romPath = flag.String("rom", "", "path to rom")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Now().UnixNano())

	// Init Window System
	r.InitWindow(800, 600, "NES golang")
	r.SetTraceLogLevel(r.LogWarning)
	//r.SetTargetFPS(60)
	font = r.LoadFont("./assets/Pixel_NES.otf")

	graphics.InitDrawer()

	fmt.Printf("Nes Emulator\n")
	//path := "./assets/roms/nestest/nestest.nes"
	//path := "./assets/roms/full_palette/full_palette.nes"
	//path := "./assets/roms/snake.nes"
	//path := "./assets/roms/Pac-Man (USA) (Namco).nes"
	//path := "./assets/roms/Donkey Kong (World) (Rev A).nes"
	//path := "./assets/roms/Super Mario Bros. (World).nes"
	//path := "./assets/roms/Mega Man 2 (Europe).nes"
	cartridge := gamePak.CreateGamePakFromROMFile(*romPath)

	printRomInfo(&cartridge)

	console := nes.CreateNes(
		&cartridge,
		nes.CreateNesDebugger("./var/run.log", true),
	)

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
		draw(&console)

		cpuAdvance = true
	}
	r.CloseWindow()

	console.Stop()
}

func draw(console *nes.Nes) {
	r.BeginDrawing()
	r.ClearBackground(r.Black)

	drawEmulation(console)
	//drawBackgroundTileIDs(console)
	DrawDebug(console)
	r.EndDrawing()
}

func drawEmulation(console *nes.Nes) {
	frame := console.Frame()

	padding := 20
	paddingY := 20
	r.DrawRectangle(padding-1, paddingY-1, types.SCREEN_WIDTH+2, types.SCREEN_HEIGHT+2, r.RayWhite) /*
		for i := 0; i < types.SCREEN_WIDTH*types.SCREEN_HEIGHT; i++ {
			pixel := frame.Pixels[i]
			color := utils.RGBA2raylibColor(pixel)
			x := i % types.SCREEN_WIDTH
			y := i / types.SCREEN_WIDTH

			r.DrawPixel(padding+x, paddingY+y, color)
		}*/
	for x := 0; x < types.SCREEN_WIDTH; x++ {
		for y := 0; y < types.SCREEN_HEIGHT; y++ {
			pixel := frame.At(x, y)
			color := utils.RGBA2raylibColor(pixel)
			//x := i % types.SCREEN_WIDTH
			//y := i / types.SCREEN_WIDTH

			r.DrawPixel(padding+x, paddingY+y, color)
		}
	}
}
