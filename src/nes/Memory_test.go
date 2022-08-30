package nes

import (
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakePPU struct {
	register types.Address
	value    byte
}

func (f *fakePPU) WriteRegister(register types.Address, value byte) {
	f.register = register
	f.value = value
}

func (f *fakePPU) ReadRegister(register types.Address) byte {
	return f.value
}

type fakeMapper struct{}

func (f fakeMapper) prgBanks() byte {
	panic("implement me")
}

func (f fakeMapper) chrBanks() byte {
	panic("implement me")
}

func (f fakeMapper) Read(address types.Address) byte {
	panic("implement me")
}

func (f fakeMapper) Write(address types.Address, value byte) {
	panic("implement me")
}

type fields struct {
	ram     [0xFFFF + 1]byte
	gamePak *gamePak2.GamePak
	mapper  gamePak2.Mapper
	ppu     *fakePPU
}

type args struct {
	address types.Address
	value   byte
}

func TestCPUMemory_Read_into_cpu_ram(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Cpu reading into ram (low edge)", fields{}, args{0x0000, 0x01}},
		{"Cpu reading into ram (high edge)", fields{}, args{RAM_HIGHER_ADDRESS, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				ppu:     tt.fields.ppu,
			}
			cm.ram[tt.args.address&RAM_LAST_REAL_ADDRESS] = tt.args.value

			actualValue := cm.Read(tt.args.address)

			assert.Equal(t, tt.args.value, actualValue)
		})
	}
}

func TestCPUMemory_Read_into_ppu(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Cpu reading into ppu (low edge)", fields{ppu: &fakePPU{}}, args{0x2000, 0x01}},
		{"Cpu reading into ppu (high edge)", fields{ppu: &fakePPU{}}, args{0x3FFF, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := tt.fields.ppu
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				ppu:     ppu,
			}
			ppu.register = tt.args.address & 0x2007
			ppu.value = tt.args.value

			actualValue := cm.Read(tt.args.address)

			assert.Equal(t, actualValue, tt.args.value)
		})
	}
}

func TestCPUMemory_Write_into_cpu_ram(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Cpu writing into ram (low edge)", fields{}, args{0x0000, 0x01}},
		{"Cpu writing into ram (high edge)", fields{}, args{RAM_HIGHER_ADDRESS, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				ppu:     tt.fields.ppu,
			}

			cm.Write(tt.args.address, tt.args.value)

			assert.Equal(t, byte(0x01), cm.ram[tt.args.address&RAM_LAST_REAL_ADDRESS])
		})
	}
}

func TestCPUMemory_Write_into_ppu(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Cpu writing into ppu (low edge)", fields{ppu: &fakePPU{}}, args{0x2000, 0x01}},
		{"Cpu writing into ppu (high edge)", fields{ppu: &fakePPU{}}, args{0x3FFF, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := tt.fields.ppu
			cm := &CPUMemory{
				ram:     tt.fields.ram,
				gamePak: tt.fields.gamePak,
				ppu:     ppu,
			}

			cm.Write(tt.args.address, tt.args.value)

			assert.Equal(t, tt.args.address&0x2007, ppu.register)
			assert.Equal(t, tt.args.value, ppu.value)
		})
	}
}

func allPressedButtons() ControllerState {
	controllerState := ControllerState{
		A:      true,
		B:      true,
		Select: true,
		Start:  true,
		Up:     true,
		Down:   true,
		Left:   true,
		Right:  true,
	}
	return controllerState
}

func TestCPUMemory_Writing_into_controller_takes_a_snapshot_of_the_controller(t *testing.T) {
	bus := CPUMemory{}
	buttons := allPressedButtons()
	bus.controllers[0] = buttons.value()
	bus.controllers[1] = buttons.value()

	bus.Write(0x4016, 1)
	assert.Equal(t, buttons.value(), bus.controllersState[0])
	bus.Write(0x4017, 1)
	assert.Equal(t, buttons.value(), bus.controllersState[1])
}

func TestCPUMemory_Reading_controller_state_gets_most_significant_bit_everytime(t *testing.T) {
	bus := CPUMemory{}
	bus.controllersState[0] = 0b10101010
	expectedReads := [8]byte{1, 0, 1, 0, 1, 0, 1, 0}

	for _, expected := range expectedReads {
		value := bus.Read(0x4016)
		assert.Equal(t, expected, value)
	}
}

func TestCPUMemory_Reading_controller_state_shifts_controllerState_everytime(t *testing.T) {
	bus := CPUMemory{}
	bus.controllersState[0] = 0b10101010

	i := 0
	expectedShiftedState := bus.controllersState[0]
	for i < 8 {
		bus.Read(0x4016)
		expectedShiftedState <<= 1
		assert.Equal(t, expectedShiftedState, bus.controllersState[0])
		i++
	}
}
