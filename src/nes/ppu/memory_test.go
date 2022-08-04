package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPU_writes_and_reads_into_palette(t *testing.T) {
	ppu := aPPU()

	for i := 0; i < 32; i++ {
		colorIndex := byte(i + 1)
		address := PaletteLowAddress + types.Address(i)
		ppu.Write(address, colorIndex)
		readValue := ppu.Read(address)
		assert.Equal(
			t,
			colorIndex,
			readValue,
			fmt.Sprintf("@%X has unexpected value", address),
		)
	}
}

type mirrorAssertion struct {
	address types.Address
	mirrors types.Address
}

func TestPPU_read_nameTables_with_mirroring(t *testing.T) {
	tests := []struct {
		name       string
		mirrorMode byte
		assertions []mirrorAssertion
	}{
		{
			name:       "vertical mirror",
			mirrorMode: gamePak.VerticalMirroring,
			assertions: []mirrorAssertion{
				// NameTable 2 mirrors NameTable 0
				{0x2000, 0x0000},
				{0x23FF, 0x03FF},
				{0x2800, 0x0000},
				{0x2BFF, 0x03FF},
				// NameTable 3 mirrors NameTable 1
				{0x2400, 0x400},
				{0x27FF, 0x7FF},
				{0x2C00, 0x400},
				{0x2FFF, 0x7FF},

				// 0x3000 mirrors 0x2000
				{0x3000, 0x000},
				{0x33FF, 0x03FF},
				{0x3800, 0x000},
				{0x3eff, 0x400 + 0x2ff},

				{0x3400, 0x400},
				{0x37FF, 0x7FF},
				{0x3C00, 0x400},
				{0x3EFF, 0x6FF},

				//{0x3B27},
			},
		},
		{
			name:       "horizontal mirroring",
			mirrorMode: gamePak.HorizontalMirroring,
			assertions: []mirrorAssertion{
				// NameTable 2 mirrors NameTable 0
				{0x2000, 0x000},
				{0x23FF, 0x3FF},
				{0x2400, 0x000},
				{0x27ff, 0x3ff},
				// NameTable 3 mirrors NameTable 1
				{0x2800, 0x400},
				{0x2BFF, 0x7FF},
				{0x2C00, 0x400},
				{0x2FFF, 0x7FF},
				// 0x3000 mirrors 0x2000
				{0x3000, 0x000},
				{0x33FF, 0x3FF},
				{0x3400, 0x000},
				{0x37FF, 0x3FF},

				{0x3800, 0x400},
				{0x3Bff, 0x7ff},
				{0x3C00, 0x400},
				{0x3EFF, 0x6FF},
			},
		},
		{
			name:       "One screen mirroring",
			mirrorMode: gamePak.OneScreenMirroring,
			assertions: []mirrorAssertion{},
		},
	}

	for _, tt := range tests {
		for _, ass := range tt.assertions {
			t.Run(
				fmt.Sprintf("%s: %x should mirror %x", tt.name, ass.address, ass.mirrors),

				func(t *testing.T) {
					assert.Equal(
						t,
						ass.mirrors,
						getNameTableAddress(tt.mirrorMode, ass.address),
					)
				})
		}
	}
}

func TestPPUMemory_write_nametables(t *testing.T) {
	type mirrored struct {
		address         types.Address
		addressMirrored types.Address
		realVRAM        types.Address
	}
	tests := []struct {
		name       string
		mirrorMode byte
		mirrorA    mirrored
		mirrorB    mirrored
	}{
		{
			"vertical mirroring: NameTable 2 mirrors NameTable 0",
			gamePak.VerticalMirroring,
			mirrored{0x2000, 0x2800, 0},
			mirrored{0x23FF, 0x2BFF, 0x3FF},
		},
		{
			"vertical mirroring NameTable 3 mirrors NameTable 1",
			gamePak.VerticalMirroring,
			mirrored{0x2400, 0x2C00, 0x400},
			mirrored{0x27FF, 0x2FFF, 0x7FF},
		},
		{
			"horizontal mirroring NameTable 2 mirrors NameTable 0",
			gamePak.HorizontalMirroring,
			mirrored{0x2000, 0x2400, 0},
			mirrored{0x23FF, 0x27FF, 0x3FF},
		},
		{
			"horizontal mirroring NameTable 3 mirrors NameTable 1",
			gamePak.HorizontalMirroring,
			mirrored{0x2800, 0x2C00, 0x400},
			mirrored{0x2BFF, 0x2FFF, 0x7FF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mirrorMode := byte(0)
			if tt.mirrorMode == gamePak.VerticalMirroring {
				mirrorMode = 1
			}
			header := gamePak.CreateINes1Header(10, 10, mirrorMode, 0, 0, 0, 0)
			pak := gamePak.CreateGamePak(header, []byte{0}, []byte{0})
			ppu := CreatePPU(&pak, false, "")

			ppu.Write(tt.mirrorA.address, 0xFF)
			assert.Equal(t, byte(0xFF), ppu.nameTables[tt.mirrorA.realVRAM])
			ppu.Write(tt.mirrorA.addressMirrored, 0xF1)
			assert.Equal(t, byte(0xF1), ppu.nameTables[tt.mirrorA.realVRAM])

			ppu.Write(tt.mirrorB.address, 0xFF)
			assert.Equal(t, byte(0xFF), ppu.nameTables[tt.mirrorB.realVRAM])
			ppu.Write(tt.mirrorB.addressMirrored, 0xF1)
			assert.Equal(t, byte(0xF1), ppu.nameTables[tt.mirrorB.realVRAM])
		})
	}
}
