package nes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateDummyGamePak() *GamePak {
	return &GamePak{
		Header{1, 1, 0, 0, 0, 0, 0},
		make([]byte, 100),
		make([]byte, 0x01FFF),
	}
}

func TestPPU_PPUADDR_write_twice_to_set_address(t *testing.T) {
	gamePak := CreateDummyGamePak()
	memory := CreatePPUMemory(gamePak)
	ppu := CreatePPU(memory)

	ppu.WriteRegister(PPUADDR, 0x28)
	assert.Equal(t, Address(0x2800), ppu.registers.ppuAddr)

	ppu.WriteRegister(PPUADDR, 0x10)
	assert.Equal(t, Address(0x2810), ppu.registers.ppuAddr)
}

func TestPPU_is_instructed_to_read_address(t *testing.T) {
	gamePak := CreateDummyGamePak()
	memory := CreatePPUMemory(gamePak)
	ppu := CreatePPU(memory)
	ppu.Write(0x2600, 0x15)

	// CPU wants to access memory cell at 0x0600 PPU memory space
	// LDA #$06
	// STA $2006
	// LDA #$00
	// STA $2006
	ppu.WriteRegister(PPUADDR, 0x26)
	ppu.WriteRegister(PPUADDR, 0x00)

	// Dummy Read
	firstRead := ppu.ReadRegister(PPUDATA)
	assert.Equal(t, byte(0x00), firstRead)
	assert.Equal(t, Address(0x2601), ppu.registers.ppuAddr, "ppuAddr(cpu@0x2006) must increment on each read to cpu@0x2007")

	secondRead := ppu.ReadRegister(PPUDATA)

	assert.Equal(t, byte(0x15), secondRead)
	assert.Equal(t, Address(0x2602), ppu.registers.ppuAddr, "ppuAddr(cpu@0x2006) must increment on each read to cpu@0x2007")
}

func TestPPU_is_instructed_to_read_address_and_mirrors(t *testing.T) {
	t.Skipf("Mirror still not implemented")
	gamePak := CreateDummyGamePak()
	memory := CreatePPUMemory(gamePak)
	ppu := CreatePPU(memory)

	ppu.WriteRegister(PPUADDR, 0x3F)
	ppu.WriteRegister(PPUADDR, 0xFF)

	// Dummy Read
	ppu.ReadRegister(PPUDATA)
	assert.Equal(t, Address(0x0000), ppu.registers.ppuAddr, "ppuAddr(cpu@0x2006) must increment on each read to cpu@0x2007")
}
