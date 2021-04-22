package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/component"
	"io/ioutil"
)

func ReadRom(path string) component.GamePak {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read Header
	inesHeader := component.CreateINes1Header(data[0:16])

	if inesHeader.HasTrainer() {
		fmt.Println("Rom has trainer")
	} else {
		fmt.Println("Rom has no trainer")
	}

	fmt.Println("PRG:", inesHeader.ProgramSize(), "x 16KB Banks")
	fmt.Println("CHR:", inesHeader.CHRSize(), "x 8KB Banks")
	fmt.Println("Mapper:", inesHeader.MapperNumber())
	fmt.Println("Tv System:", inesHeader.TvSystem())

	return component.CreateGamePak(inesHeader, data[16:])
}
