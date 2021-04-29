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
	cases := []struct {
		name     string
		hi       byte
		lo       byte
		expected Address
	}{
		{"writes address", 0x28, 0x10, 0x2810},
		{"writes address > 0x3FFF is mirrored down", 0x40, 0x20, 0x0020},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gamePak := CreateDummyGamePak()
			memory := CreatePPUMemory(gamePak)
			ppu := CreatePPU(memory)

			ppu.WriteRegister(PPUADDR, tt.hi)
			assert.Equal(t, Address(tt.hi)<<8, ppu.registers.ppuAddr)

			ppu.WriteRegister(PPUADDR, tt.lo)
			assert.Equal(t, tt.expected, ppu.registers.ppuAddr)
		})
	}
}

func TestPPU_PPUData_read(t *testing.T) {
	gamePak := CreateDummyGamePak()
	memory := CreatePPUMemory(gamePak)

	cases := []struct {
		name          string
		addressToRead Address
		incrementMode byte
	}{
		{"buffered read, increment mode going across", 0x2600, 0},
		{"buffered read, increment mode going down", 0x2600, 1},
	}

	ppu := CreatePPU(memory)
	ppu.Write(0x2600, 0x15)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ppu.registers.ppuAddr = tt.addressToRead
			ppu.ppuctrlWriteFlag(incrementMode, tt.incrementMode)
			expectedIncrement := Address(1)
			if tt.incrementMode == 1 {
				expectedIncrement = 32
			}

			// Dummy Read
			firstRead := ppu.ReadRegister(PPUDATA)
			assert.Equal(t, byte(0x00), firstRead)
			assert.Equal(t, tt.addressToRead+expectedIncrement, ppu.registers.ppuAddr, "ppuAddr(cpu@0x%X) must increment on each read to cpu@0x%X")

			secondRead := ppu.ReadRegister(PPUDATA)

			assert.Equal(t, byte(0x15), secondRead)

			assert.Equal(t, tt.addressToRead+expectedIncrement*2, ppu.registers.ppuAddr, "unexpected ppuAddr increment")
		})
	}
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
