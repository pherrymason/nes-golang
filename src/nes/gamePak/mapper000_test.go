package gamePak

import (
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateMapper000ForTest(prgROMSize byte) Mapper {
	rom := prgROM()

	header := CreateINes1Header(prgROMSize, 1, 0, 0, 0, 0, 0)

	return CreateMapper(header, rom, make([]byte, 100))
}

func prgROM() []byte {
	rom := make([]byte, 0xFFFF)
	// WritePrgROM on boundaries
	rom[0x0000] = 0x01 // First byte
	rom[0x3FFF] = 0xFF // 16KB Rom boundary
	rom[0x4000] = 0x40 // 32KB Rom low boundary
	rom[0x7FFF] = 0x7F // 32KB Rom boundary
	return rom
}

func TestReads_from_last_byte_in_16KB_ROM(t *testing.T) {
	oneBank := byte(1)
	mapper := CreateMapper000ForTest(oneBank)
	startOfCPUMap := types.Address(0x8000)

	result := mapper.ReadPrgROM(startOfCPUMap)
	assert.Equal(t, byte(0x01), result)

	result = mapper.ReadPrgROM(startOfCPUMap + 0x3FFF)
	assert.Equal(t, byte(0xFF), result)

	result = mapper.ReadPrgROM(startOfCPUMap + 0x4000)
	assert.Equal(t, byte(0x01), result, "mirroring failed")
}

func TestReads_from_last_byte_in_32KB_ROM(t *testing.T) {
	twoBanks := byte(2)
	mapper := CreateMapper000ForTest(twoBanks)
	startOfCPUMap := types.Address(0x8000)

	result := mapper.ReadPrgROM(startOfCPUMap)
	assert.Equal(t, byte(0x01), result)

	result = mapper.ReadPrgROM(startOfCPUMap + 0x3FFF)
	assert.Equal(t, byte(0xFF), result)

	result = mapper.ReadPrgROM(startOfCPUMap + 0x4000)
	assert.Equal(t, byte(0x40), result)

	result = mapper.ReadPrgROM(startOfCPUMap + 0x7FFF)
	assert.Equal(t, byte(0x7F), result)
}
