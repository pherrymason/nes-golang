package nes

import (
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
