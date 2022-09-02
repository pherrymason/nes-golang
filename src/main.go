package main

import (
	"flag"
	"github.com/raulferras/nes-golang/src/app"
	_ "net/http/pprof"
)

func main() {
	appOptions := cmdLineArguments()
	app.RunEmulator(appOptions)
}

func cmdLineArguments() app.Options {
	var cpuprofile = flag.Bool("cpuprofile", false, "write cpu profile to file")
	var romPath = flag.String("rom", "", "path to rom")
	var logCPU = flag.Bool("logCPU", false, "enables CPU log")
	var debugPPU = flag.Bool("debugPPU", false, "Displays PPU debug information")
	var scale = flag.Int("scale", 1, "scale resolution")
	var breakpoint = flag.String("breakpoint", "", "defines a breakpoint on start")
	flag.Parse()

	return app.NewOptions(*scale, *romPath, *logCPU, *debugPPU, *breakpoint, *cpuprofile)
}
