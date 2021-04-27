package nes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateMapper000ForTest(prgROMSize byte) Mapper {
	rom := prgROM()

	gamePak := CreateGamePak(
		Header{prgROMSize, 1, 0, 0, 0, 0, 0},
		rom,
		make([]byte, 100),
	)

	return CreateMapper(&gamePak)
}

func prgROM() []byte {
	rom := make([]byte, 0xFFFF)
	// Write on boundaries
	rom[0x0000] = 0x01 // First byte
	rom[0x3FFF] = 0xFF // 16KB Rom boundary
	rom[0x4000] = 0x40 // 32KB Rom low boundary
	rom[0x7FFF] = 0x7F // 32KB Rom boundary
	return rom
}

func TestReads_from_last_byte_in_16KB_ROM(t *testing.T) {
	oneBank := byte(1)
	mapper := CreateMapper000ForTest(oneBank)
	startOfCPUMap := Address(0x8000)

	result := mapper.Read(startOfCPUMap)
	assert.Equal(t, byte(0x01), result)

	result = mapper.Read(startOfCPUMap + 0x3FFF)
	assert.Equal(t, byte(0xFF), result)

	result = mapper.Read(startOfCPUMap + 0x4000)
	assert.Equal(t, byte(0x01), result, "mirroring failed")
}

func TestReads_from_last_byte_in_32KB_ROM(t *testing.T) {
	twoBanks := byte(2)
	mapper := CreateMapper000ForTest(twoBanks)
	startOfCPUMap := Address(0x8000)

	result := mapper.Read(startOfCPUMap)
	assert.Equal(t, byte(0x01), result)

	result = mapper.Read(startOfCPUMap + 0x3FFF)
	assert.Equal(t, byte(0xFF), result)

	result = mapper.Read(startOfCPUMap + 0x4000)
	assert.Equal(t, byte(0x40), result)

	result = mapper.Read(startOfCPUMap + 0x7FFF)
	assert.Equal(t, byte(0x7F), result)
}
