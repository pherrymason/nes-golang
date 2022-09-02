package app

import (
	r "github.com/gen2brain/raylib-go/raylib"
	"github.com/pkg/profile"
	"github.com/raulferras/nes-golang/src/audio"
	"github.com/raulferras/nes-golang/src/debugger"
	"github.com/raulferras/nes-golang/src/nes"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
)

type Options struct {
	videoScale int
	romPath    string
	logCPU     bool
	debugPPU   bool
	breakpoint string
	cpuProfile bool
}

func NewOptions(videoScale int,
	romPath string,
	logCPU bool,
	debugPPU bool,
	breakpoint string,
	cpuProfile bool) Options {
	return Options{
		videoScale: videoScale,
		romPath:    romPath,
		logCPU:     logCPU,
		debugPPU:   debugPPU,
		breakpoint: breakpoint,
		cpuProfile: cpuProfile,
	}
}

func RunEmulator(options Options) {
	// Init Window System
	var windowWidth int32
	windowWidth = 800
	if options.debugPPU {
		windowWidth += 400
	}
	r.InitWindow(windowWidth, 700, "NES golang")
	r.SetTraceLog(r.LogWarning)
	r.SetTargetFPS(30)
	font := r.LoadFont("./assets/Pixel_NES.otf")
	r.SetTextureFilter(font.Texture, r.FilterPoint)

	audioDevice := audio.NewAudio()
	audioDevice.Init()
	defer audioDevice.Stop()
	audio.Generate()

	nesDebugger := nes.CreateNesDebugger(
		"./var",
		options.logCPU,
		options.debugPPU,
	)

	cartridge := gamePak.CreateGamePakFromROMFile(options.romPath)
	console := nes.CreateNes(
		&cartridge,
		nesDebugger,
	)

	debugger.PrintRomInfo(&cartridge)
	if options.cpuProfile {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	loop(console, options.videoScale, audioDevice)

	r.UnloadFont(font)
	r.CloseAudioDevice()
	r.CloseWindow()
}

func loop(console *nes.Nes, videoScale int, audioDevice *audio.Audio) {
	console.Start()
	_timestamp := r.GetTime()
	debuggerGUI := debugger.NewDebugger(console, audioDevice)

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
		audioDevice.Update()

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
		texture := drawEmulation(console.Frame(), videoScale)
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
	padding := int32(20)
	paddingY := int32(20)
	screenWidth := int32(types.SCREEN_WIDTH) * int32(scale)
	screenHeight := int32(types.SCREEN_HEIGHT) * int32(scale)
	r.DrawRectangle(padding-1, paddingY-1, screenWidth+2, screenHeight+2, r.RayWhite)

	image := r.NewImageFromImage(frame)
	texture := r.LoadTextureFromImage(image)
	r.DrawTextureEx(texture, r.Vector2{X: float32(padding), Y: float32(paddingY)}, 0, float32(scale), r.White)
	return texture
}
