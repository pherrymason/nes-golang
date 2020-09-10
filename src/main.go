package main

import (
	"fmt"

	"github.com/raulferras/nes-golang/src/nes"
)

func main() {
	fmt.Printf("Nes Emulator\n")

	console := nes.CreateNes()
	console.Start()
}
