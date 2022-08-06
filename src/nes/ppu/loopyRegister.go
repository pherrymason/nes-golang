package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

/*
LoopyRegister Internal PPU Register
The 15 bit registers t and v are composed this way during rendering:
	yyy NN YYYYY XXXXX
	||| || ||||| +++++-- coarse X scroll
	||| || +++++-------- coarse Y scroll
	||| ++-------------- nametable select
	+++----------------- fine Y scroll
*/
type LoopyRegister struct {
	_coarseX    byte // 5 bits
	_coarseY    byte // 5 bits
	_nameTableX byte // 1 bit
	_nameTableY byte // 1 bit
	_fineY      byte // 3 bits

	latch byte
}

func (register *LoopyRegister) address() types.Address {
	var address uint16 = 0
	address |= uint16(register._coarseX & 0b11111)
	address |= uint16(register._coarseY&0b11111) << 5
	address |= uint16(register._nameTableX&1) << 10
	address |= uint16(register._nameTableY&1) << 11
	address |= uint16(register._fineY&0b111) << 12
	return types.Address(address) & 0x3FFF
}

func (register *LoopyRegister) nameTableAddress() types.Address {
	var address uint16 = 0
	address |= uint16(register._coarseX & 0b11111)
	address |= uint16(register._coarseY&0b11111) << 5
	address |= uint16(register._nameTableX&1) << 10
	address |= uint16(register._nameTableY&1) << 11
	return types.Address(address) & 0x3FFF
}

func (register *LoopyRegister) Value() types.Address {
	var address uint16 = 0
	address |= uint16(register._coarseX & 0b11111)
	address |= uint16(register._coarseY&0b11111) << 5
	address |= uint16(register._nameTableX&1) << 10
	address |= uint16(register._nameTableY&1) << 11
	address |= uint16(register._fineY&0b111) << 12
	return types.Address(address)
}

func (register *LoopyRegister) setValue(address types.Address) {
	register._coarseX = byte(address & 0b11111)
	register._coarseY = byte(address>>5) & 0b11111
	register._nameTableX = byte(address>>10) & 1
	register._nameTableY = byte(address>>11) & 1
	register._fineY = byte(address>>12) & 0b111
}

func (register *LoopyRegister) setNameTableX(nameTableX byte) {
	register._nameTableX = nameTableX & 1
}

func (register *LoopyRegister) setNameTableY(nameTableY byte) {
	register._nameTableY = nameTableY & 1
}

func (register *LoopyRegister) setCoarseX(coarseX byte) {
	register._coarseX = coarseX & 0b11111
}

func (register *LoopyRegister) setCoarseY(coarseY byte) {
	register._coarseY = coarseY & 0b11111
}

func (register *LoopyRegister) push(value byte) {
	/*
		if register.latch == 0 {
			register.address = types.Address(Value&0x3F)<<8 |
					(register.address & 0x00FF)
		} else {
			register.address = register.address&0xFF00 | types.Address(Value)
		}
	*/
	// fY fY fy NY NX YYYYY XXXXX
	if register.latch == 0 {
		//     fy fy fy NY NX Y Y] [Y Y Y X X X X X]
		// [15 14 13 12 11 10 9 8   7 6 5 4 3 2 1 0]
		//   v  v  v  v  v  v v v
		// 0x3F => 111111 <<8 => 11111100000000
		// Bits 0 and 1 go into bits 3 and 4 of CoarseY
		register._coarseY = (value & 0x03 << 3) | register._coarseY&0b111

		// bit 2 goes into NameTableX
		register._nameTableX = (value & 0x04) >> 2
		// bit 3 goes into NameTableY
		register._nameTableY = (value & 0x08) >> 3
		// bit 5,6 and 7 goes into FineY. bit 7 is cleared
		register._fineY = (value >> 4) & 0b011
		//register.address =
		//	types.Address(Value&0x3F)<<8 | (register.address & 0x00FF)
		register.latch = 1
	} else {
		// bits 0,1,2,3, and 4 goes into CoarseX
		register._coarseX = value & 0b11111

		// bits 5,6 and 7 go into CoarseY[0,1,2]
		register._coarseY = register._coarseY&0b11111000 | (value&0b11100000)>>5
		register.latch = 0
	}
}

func (register *LoopyRegister) increment(incrementMode byte) {
	address := register.Value()
	if incrementMode == 0 {
		address++
	} else {
		address += 32
	}
	register.setValue(address)
}

func (register *LoopyRegister) resetLatch() {
	register.latch = 0
}

func (register *LoopyRegister) NameTableY() byte {
	return register._nameTableY
}

func (register *LoopyRegister) NameTableX() byte {
	return register._nameTableX
}

func (register *LoopyRegister) CoarseX() byte {
	return register._coarseX
}

func (register *LoopyRegister) CoarseY() byte {
	return register._coarseY
}

func (register *LoopyRegister) FineY() byte {
	return register._fineY
}

func (register *LoopyRegister) resetFineY() {
	register._fineY = 0
}
