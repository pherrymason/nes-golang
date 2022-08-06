package debugger

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
)

func PrintRomInfo(cartridge *gamePak.GamePak) {
	inesHeader := cartridge.Header()

	if inesHeader.HasTrainer() {
		fmt.Println("Rom has trainer")
	} else {
		fmt.Println("Rom has no trainer")
	}

	if inesHeader.Mirroring() == gamePak.VerticalMirroring {
		fmt.Println("Vertical Mirroring")
	} else {
		fmt.Println("Horizontal Mirroring")
	}

	fmt.Println("PRG:", inesHeader.ProgramSize(), "x 16KB Banks")
	fmt.Println("CHR:", inesHeader.CHRSize(), "x 8KB Banks")
	fmt.Println("Mapper:", inesHeader.MapperNumber())
	fmt.Println("Tv System:", inesHeader.TvSystem())
}
