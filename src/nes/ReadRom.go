package nes

import (
	"fmt"
	"io/ioutil"
)

func readRom(path string) GamePak {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read Header
	inesHeader := CreateINes1Header(data[0:16])

	if inesHeader.HasTrainer() {
		fmt.Println("Rom has trainer")
	} else {
		fmt.Println("Rom has no trainer")
	}

	fmt.Println("PRG:", inesHeader.ProgramSize(), "x 16KB Banks")
	fmt.Println("CHR:", inesHeader.CHRSize(), "x 8KB Banks")
	fmt.Println("Mapper:", inesHeader.MapperNumber())
	fmt.Println("Tv System:", inesHeader.TvSystem())

	return GamePak{inesHeader, data[16:]}
}
