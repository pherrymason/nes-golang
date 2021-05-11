package nes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakePPU struct {
	register Address
	value    byte
}

func (f *fakePPU) WriteRegister(register Address, value byte) {
	f.register = register
	f.value = value
}

func (f *fakePPU) ReadRegister(register Address) byte {
	panic("implement me")
}

type fakeMapper struct{}

func (f fakeMapper) prgBanks() byte {
	panic("implement me")
}

func (f fakeMapper) chrBanks() byte {
	panic("implement me")
}

func (f fakeMapper) Read(address Address) byte {
	panic("implement me")
}

func (f fakeMapper) Write(address Address, value byte) {
	panic("implement me")
}

func TestCPUMemory_Write_into_cpu_ram(t *testing.T) {
	type fields struct {
		ram     [0xFFFF + 1]byte
		gamePak *GamePak
		mapper  Mapper
		ppu     PPU
	}
	type args struct {
		address Address
		value   byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"cpu writing into ram (low edge)", fields{}, args{0x0000, 0x01}},
		{"cpu writing into ram (high edge)", fields{}, args{RAM_HIGHER_ADDRESS, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				mapper:  tt.fields.mapper,
				ppu:     tt.fields.ppu,
			}

			cm.Write(tt.args.address, tt.args.value)

			assert.Equal(t, byte(0x01), cm.ram[tt.args.address&RAM_LAST_REAL_ADDRESS])
		})
	}
}

func TestCPUMemory_Write_into_ppu(t *testing.T) {
	type fields struct {
		ram     [0xFFFF + 1]byte
		gamePak *GamePak
		mapper  Mapper
		ppu     fakePPU
	}
	type args struct {
		address Address
		value   byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"cpu writing into ppu (low edge)", fields{}, args{0x2000, 0x01}},
		{"cpu writing into ppu (high edge)", fields{}, args{0x3FFF, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := tt.fields.ppu
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				mapper:  tt.fields.mapper,
				ppu:     &ppu,
			}

			cm.Write(tt.args.address, tt.args.value)

			assert.Equal(t, tt.args.address&0x2007, ppu.register)
			assert.Equal(t, tt.args.value, ppu.value)
		})
	}
}
