package main

import (
	"flag"
	"fmt"
	"github.com/FMNSSun/hexit"
	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/pkg/profile"
	"github.com/raulferras/nes-golang/src/debugger"
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/raulferras/nes-golang/src/utils"
	"image"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

var cpuAdvance bool

func main() {
	scale, cpuprofile, romPath, logCPU, debugPPU, breakpoint := cmdLineArguments()
	if cpuprofile != "" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	rand.Seed(time.Now().UnixNano())

	// Init Window System
	var windowWidth int
	windowWidth = 800
	if debugPPU {
		windowWidth += 400
	}
	r.InitWindow(windowWidth, 700, "NES golang")
	r.SetTraceLogLevel(r.LogWarning)
	r.SetTargetFPS(60)

	graphics.InitDrawer()

	fmt.Printf("Nes Emulator\n")
	cartridge := gamePak.CreateGamePakFromROMFile(romPath)
	debugger.PrintRomInfo(&cartridge)

	nesDebugger := nes.CreateNesDebugger(
		"./var",
		logCPU,
		debugPPU,
	)
	if len(breakpoint) > 0 {
		nesDebugger.AddBreakPoint(types.Address(hexit.UnhexUint16Str(breakpoint)))
	}
	console := nes.CreateNes(
		&cartridge,
		nesDebugger,
	)

	loop(console, scale)
	r.CloseWindow()
}

func loop(console *nes.Nes, scale int) {
	cpuAdvance = true
	console.Start()
	_timestamp := r.GetTime()
	debuggerGUI := debugger.NewDebugger(console)

	for !r.WindowShouldClose() {
		if console.Finished() {
			break
		}

		timestamp := r.GetTime()
		dt := timestamp - _timestamp
		_timestamp = timestamp
		if dt > 1 {
			// difference too big
			dt = 0
		}

		// Update emulator

		if !console.Paused() {
			//console.TickForTime(dt)
			console.TickTillFrameComplete()
		} else {
			console.PausedTick()
		}

		// Draw --------------------

		r.BeginDrawing()
		r.ClearBackground(r.Black)
		drawEmulation(console.Frame())
		debuggerGUI.Tick()
		r.EndDrawing()

		// End Draw --------------------

		cpuAdvance = true
	}

	debuggerGUI.Close()
	console.Stop()
}

func drawEmulation(frame *image.RGBA) {
	padding := 20
	paddingY := 20
	r.DrawRectangle(padding-1, paddingY-1, types.SCREEN_WIDTH+2, types.SCREEN_HEIGHT+2, r.RayWhite)
	for x := 0; x < types.SCREEN_WIDTH; x++ {
		for y := 0; y < types.SCREEN_HEIGHT; y++ {
			pixel := frame.At(x, y)
			color := utils.RGBA2raylibColor(pixel)
			r.DrawPixel(padding+x, paddingY+y, color)
		}
	}
}

func cmdLineArguments() (int, string, string, bool, bool, string) {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var romPath = flag.String("rom", "", "path to rom")
	var logCPU = flag.Bool("logCPU", false, "enables CPU log")
	var debugPPU = flag.Bool("debugPPU", false, "Displays PPU debug information")
	var scale = flag.Int("scale", 1, "scale resolution")
	var breakpoint = flag.String("breakpoint", "", "defines a breakpoint on start")
	flag.Parse()

	return *scale, *cpuprofile, *romPath, *logCPU, *debugPPU, *breakpoint
}
