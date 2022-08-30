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
	"image"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

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
		controllerState := readController()
		console.UpdateController(1, controllerState)

		if !console.Paused() {
			//console.TickForTime(dt)
			console.TickTillFrameComplete()
		} else {
			console.PausedTick()
		}

		// Draw --------------------

		r.BeginDrawing()
		r.ClearBackground(r.Black)
		texture := drawEmulation(console.Frame(), scale)
		debuggerGUI.Tick()
		r.EndDrawing()
		r.UnloadTexture(texture)
		// End Draw --------------------
	}

	debuggerGUI.Close()
	console.Stop()
}

func readController() nes.ControllerState {
	state := nes.ControllerState{
		A:      r.IsKeyDown(r.KeyZ),
		B:      r.IsKeyDown(r.KeyX),
		Select: r.IsKeyDown(r.KeyA),
		Start:  r.IsKeyDown(r.KeyS),
		Up:     r.IsKeyDown(r.KeyUp),
		Down:   r.IsKeyDown(r.KeyDown),
		Left:   r.IsKeyDown(r.KeyLeft),
		Right:  r.IsKeyDown(r.KeyRight),
	}

	return state
}

func drawEmulation(frame image.Image, scale int) r.Texture2D {
	padding := 20
	paddingY := 20
	screenWidth := types.SCREEN_WIDTH * scale
	screenHeight := types.SCREEN_HEIGHT * scale
	r.DrawRectangle(padding-1, paddingY-1, screenWidth+2, screenHeight+2, r.RayWhite)
	texture := r.LoadTextureFromGo(frame)
	r.DrawTextureEx(texture, r.Vector2{X: float32(padding), Y: float32(paddingY)}, 0, float32(scale), r.White)
	return texture
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
