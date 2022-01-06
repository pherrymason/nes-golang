package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/mocks"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPUMemory_read(t *testing.T) {
	type fields struct {
		gamePak      *gamePak.GamePak
		vram         [2048]byte
		oamData      [256]byte
		paletteTable [32]byte
		mirroring    byte
	}
	type args struct {
		address types.Address
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   byte
	}{
		{"vertical mirroring, reading vram, low edge", fields{mirroring: gamePak.VerticalMirroring}, args{0x2000}, 0x01},
		{"vertical mirroring, reading vram, high edge", fields{mirroring: gamePak.VerticalMirroring}, args{0x2800}, 0x01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := new(mocks.MockableHeader)
			pak := gamePak.CreateGamePak(header, []byte{0}, []byte{0})
			ppu := &PPUMemory{
				gamePak:      &pak,
				vram:         tt.fields.vram,
				paletteTable: tt.fields.paletteTable,
			}

			header.On("Mirroring").Return(tt.fields.mirroring)
			ppu.vram[(tt.args.address-0x2000)&0x27FF] = tt.want
			if got := ppu.read(tt.args.address, false); got != tt.want {
				t.Errorf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPPUMemory_read_nametables(t *testing.T) {
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
		{
			"one screen mirroring",
			gamePak.OneScreenMirroring,
			mirrored{0x2000, 0x2400, 0x000},
			mirrored{0x23FF, 0x27FF, 0x3FF},
		},
		{
			"one screen mirroring > 0x2800",
			gamePak.OneScreenMirroring,
			mirrored{0x2000, 0x2800, 0x000},
			mirrored{0x23FF, 0x2BFF, 0x3FF},
		},
		{
			"one screen mirroring > 0x2C00",
			gamePak.OneScreenMirroring,
			mirrored{0x2000, 0x2C00, 0x000},
			mirrored{0x23FF, 0x2FFF, 0x3FF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := gamePak.CreateINes1Header(10, 10, tt.mirrorMode, 0, 0, 0, 0)
			pak := gamePak.CreateGamePak(header, []byte{0}, []byte{0})
			ppu := &PPUMemory{
				gamePak:      &pak,
				vram:         [2048]byte{},
				paletteTable: [32]byte{},
			}

			ppu.vram[tt.mirrorA.realVRAM] = 0xFF
			assert.Equal(t, byte(0xFF), ppu.Read(tt.mirrorA.address))
			assert.Equal(t, byte(0xFF), ppu.Read(tt.mirrorA.addressMirrored), fmt.Sprintf("failed %x mirrors %x", tt.mirrorA.addressMirrored, tt.mirrorA.address))

			ppu.vram[tt.mirrorB.realVRAM] = 0xF1
			assert.Equal(t, byte(0xF1), ppu.Read(tt.mirrorB.address))
			assert.Equal(t, byte(0xF1), ppu.Read(tt.mirrorB.addressMirrored), fmt.Sprintf("failed %x mirrors %x", tt.mirrorB.addressMirrored, tt.mirrorB.address))
		})
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
			ppu := &PPUMemory{
				gamePak:      &pak,
				vram:         [2048]byte{},
				paletteTable: [32]byte{},
			}

			ppu.Write(tt.mirrorA.address, 0xFF)
			assert.Equal(t, byte(0xFF), ppu.vram[tt.mirrorA.realVRAM])
			ppu.Write(tt.mirrorA.addressMirrored, 0xF1)
			assert.Equal(t, byte(0xF1), ppu.vram[tt.mirrorA.realVRAM])

			ppu.Write(tt.mirrorB.address, 0xFF)
			assert.Equal(t, byte(0xFF), ppu.vram[tt.mirrorB.realVRAM])
			ppu.Write(tt.mirrorB.addressMirrored, 0xF1)
			assert.Equal(t, byte(0xF1), ppu.vram[tt.mirrorB.realVRAM])
		})
	}
}