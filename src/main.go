package main

import (
	"github.com/pkg/profile"
	"github.com/raulferras/nes-golang/src/gui"
)

func main() {
	defer profile.Start(profile.ProfilePath("."), profile.CPUProfile).Stop()
	gui.Run()
}
