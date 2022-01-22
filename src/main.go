package main

import (
	"fmt"
	"github.com/pkg/profile"
	"github.com/raulferras/nes-golang/src/gui"
	"net/http"
)
import _ "net/http/pprof"

func main() {
	// Server for pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	defer profile.Start(profile.ProfilePath("."), profile.CPUProfile).Stop()
	gui.Run()
}
