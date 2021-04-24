package nes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_read_gamePak_without_trainer(t *testing.T) {
	gamePak := aSampleGamePak(false)

	value := gamePak.read(GAMEPAK_ROM_LOWER_BANK_START)
	assert.Equal(t, byte(0x01), value)
}

func Test_read_gamePak_on_mirror_space(t *testing.T) {
	gamePak := aSampleGamePak(false)

	value := gamePak.read(0xC100)
	assert.Equal(t, byte(0xFF), value)
}

/*
func Test_read_gamePak_with_trainer(t *testing.T) {
	gamePak := aSampleGamePak(true)

	value := gamePak.read(GAMEPAK_ROM_LOWER_BANK_START)
	assert.Equal(t, byte(0x02), value)
}*/

func Test_write(t *testing.T) {
	gamePak := aSampleGamePak(false)
	gamePak.write(GAMEPAK_ROM_LOWER_BANK_START+1, 0x02)

	assert.Equal(t, byte(0x02), gamePak.prgROM[0x01])
}

func Test_write_on_mirror_space(t *testing.T) {
	gamePak := aSampleGamePak(false)
	gamePak.write(0xC100, 0x02)

	assert.Equal(t, byte(0x02), gamePak.prgROM[0xc100-0x4000-GAMEPAK_ROM_LOWER_BANK_START])
}

func aSampleGamePak(withTrainer bool) GamePak {
	data := make([]byte, GAMEPAK_MEMORY_SIZE)

	// Set data in limits
	data[0x00] = 0x01
	data[512] = 0x02

	// NROM mirroring from 0xC000
	data[0xC100-0x4000-GAMEPAK_ROM_LOWER_BANK_START] = 0xFF

	flag6 := byte(0x00)
	if withTrainer {
		flag6 = 0x06
	}
	return CreateGamePak(
		Header{
			255,
			255,
			flag6,
			0,
			0,
			0,
			0,
		},
		data,
	)
}
