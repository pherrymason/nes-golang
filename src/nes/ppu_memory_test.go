package nes

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPUMemory_read(t *testing.T) {
	type fields struct {
		gamePak      *GamePak
		vram         [2048]byte
		oamData      [256]byte
		paletteTable [32]byte
	}
	type args struct {
		address Address
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   byte
	}{
		{"reading vram, low edge", fields{}, args{0x2000}, 0x01},
		{"reading vram, high edge", fields{}, args{0x2800}, 0x01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := &PPUMemory{
				gamePak:      tt.fields.gamePak,
				vram:         tt.fields.vram,
				paletteTable: tt.fields.paletteTable,
			}

			ppu.vram[(tt.args.address-0x2000)&0x27FF] = tt.want
			if got := ppu.read(tt.args.address, false); got != tt.want {
				t.Errorf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPPUMemory_read_nametables(t *testing.T) {
	type mirrored struct {
		address         Address
		addressMirrored Address
		realVRAM        Address
	}
	tests := []struct {
		name       string
		mirrorMode byte
		mirrorA    mirrored
		mirrorB    mirrored
	}{
		{
			"vertical mirroring",
			gamePak.VerticalMirroring,
			mirrored{0x2000, 0x2800, 0},
			mirrored{0x2400, 0x2C00, 0x400},
		},
		{
			"horizontal mirroring",
			gamePak.HorizontalMirroring,
			mirrored{0x2000, 0x2400, 0},
			mirrored{0x2800, 0x2C00, 0x400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mirrorMode := byte(0)
			if tt.mirrorMode == gamePak.VerticalMirroring {
				mirrorMode = 1
			}
			header := gamePak.CreateINes1Header(10, 10, mirrorMode, 0, 0, 0, 0)
			pak := CreateGamePak(header, []byte{0}, []byte{0})
			ppu := &PPUMemory{
				gamePak:      &pak,
				vram:         [2048]byte{},
				paletteTable: [32]byte{},
			}

			ppu.vram[tt.mirrorA.realVRAM] = 0xFF
			assert.Equal(t, byte(0xFF), ppu.Read(tt.mirrorA.address))
			assert.Equal(t, byte(0xFF), ppu.Read(tt.mirrorA.addressMirrored))

			ppu.vram[tt.mirrorB.realVRAM] = 0xF1
			assert.Equal(t, byte(0xF1), ppu.Read(tt.mirrorB.address))
			assert.Equal(t, byte(0xF1), ppu.Read(tt.mirrorB.addressMirrored))
		})
	}
}

func TestPPUMemory_write(t *testing.T) {
	type fields struct {
		gamePak      *GamePak
		vram         [2048]byte
		oamData      [256]byte
		paletteTable [32]byte
	}
	type args struct {
		address Address
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   byte
	}{
		{"writing vram, low edge", fields{}, args{0x2000}, 0x01},
		{"writing vram, high edge", fields{}, args{0x2800}, 0x01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := &PPUMemory{
				gamePak:      tt.fields.gamePak,
				vram:         tt.fields.vram,
				paletteTable: tt.fields.paletteTable,
			}

			ppu.Write(tt.args.address, tt.want)

			got := ppu.vram[(tt.args.address-0x2000)&0x27FF]
			assert.Equal(t, tt.want, got)
		})
	}
}
