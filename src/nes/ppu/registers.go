package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type Control struct {
	nameTableX                    byte // Most significant bit of scrolling coordinates (X)
	nameTableY                    byte // Most significant bit of scrolling coordinates (Y)
	incrementMode                 byte // Address increment per CPU IO of PPUDATA. (0: add 1, going across; 1: add 32, going down)
	spritePatterTableAddress      byte // Sprite pattern table address for 8x8 sprites. (0: $0000; 1: $1000; ignored in 8x16 mode)
	backgroundPatternTableAddress byte // Background pattern table address (0: $0000; 1: $1000)
	spriteSize                    byte
	masterSlaveSelect             byte
	generateNMIAtVBlank           bool // NMI enabled/disabled
}

func (control *Control) value() byte {
	ctrl := byte(0)
	ctrl |= control.nameTableX
	ctrl |= control.nameTableY << 1
	ctrl |= control.incrementMode << 2
	ctrl |= control.spritePatterTableAddress << 3
	ctrl |= control.backgroundPatternTableAddress << 4
	ctrl |= control.spriteSize << 5
	ctrl |= control.masterSlaveSelect << 6
	if control.generateNMIAtVBlank {
		ctrl |= 1 << 7
	}

	return ctrl
}

type Scroll struct {
	scrollX byte
	scrollY byte
	latch   byte
}

func (scroll *Scroll) write(value byte) {
	if scroll.latch == 0 {
		scroll.scrollX = value
	} else {
		scroll.scrollY = value
	}

	// flip latch
	scroll.latch = (scroll.latch + 1) & 0x01
}

type MASKFlag int

const (
	greyScale MASKFlag = iota
	showBackgroundLeftEdge
	showSpritesLeftEdge
	showBackground
	showSprites
	emphasizeRed
	emphasizeGreen
	emphasizeBlue
)

type Status struct {
	spriteOverflow       byte
	sprite0Hit           byte
	verticalBlankStarted bool
}

func (status *Status) value() byte {
	value := byte(0)
	value |= status.spriteOverflow << 5
	value |= status.sprite0Hit << 6
	if status.verticalBlankStarted {
		value |= 1 << 7
	}

	return value
}

/*
loopyRegister Internal PPU Register
The 15 bit registers t and v are composed this way during rendering:
	yyy NN YYYYY XXXXX
	||| || ||||| +++++-- coarse X scroll
	||| || +++++-------- coarse Y scroll
	||| ++-------------- nametable select
	+++----------------- fine Y scroll
*/
type loopyRegister struct {
	/*
		coarseX    byte // 5 bits
		coarseY    byte // 5 bits
		nameTableX byte // 1 bit
		nameTableY byte // 1 bit
		fineY      byte // 3 bits
	*/
	latch   byte
	address types.Address
}

func (register *loopyRegister) setNameTableX(nameTableX byte) {
	register.address &= 0b111101111111111
	register.address |= types.Address(nameTableX) << 10
}

func (register *loopyRegister) setNameTableY(nameTableY byte) {
	register.address &= 0b111101111111111
	register.address |= types.Address(nameTableY) << 11
}

func (register *loopyRegister) setCoarseX(coarseX byte) {
	register.address &= 0b111111111100000
	register.address |= types.Address(coarseX)
}

func (register *loopyRegister) setCoarseY(coarseY byte) {
	register.address &= 0b111110000011111
	register.address |= types.Address(coarseY) << 5
}

func (register *loopyRegister) setAddress(address types.Address) {
	fmt.Printf("%X\n", address)

	register.address = address & 0x3FFF
	/*
		address &= 0x3FFF
		register.coarseX = byte(address & 0b11111)
		register.coarseY = byte(address&0b11111) >> 5
		register.nameTableX = byte(address>>10) & 1
		register.nameTableY = byte(address>>11) & 1
		register.fineY = byte(address>>12) & 0b111
	*/
}

func (register *loopyRegister) push(value byte) {
	if register.latch == 0 {
		register.address =
			types.Address(value&0x3F)<<8 |
				(register.address & 0x00FF)
	} else {
		register.address = register.address&0xFF00 | types.Address(value)
	}

	register.latch++
	register.latch &= 0b1

	/*
		if register.latch == 0 {
			register.address &= 0x3FFF
		}*/
}

/*
func (register *loopyRegister) addr() types.Address {
	address := types.Address(register.coarseX)
	address |= types.Address(register.coarseY << 5)
	address |= types.Address(register.nameTableX) << 10
	address |= types.Address(register.nameTableY) << 11
	address |= types.Address(register.fineY) << 12
	return address
}*/

func (register *loopyRegister) increment(incrementMode byte) {
	if incrementMode == 0 {
		register.address++
	} else {
		register.address += 32
	}
	register.address &= 0x3FFF
}

func (register *loopyRegister) resetLatch() {
	register.latch = 0
}

func (register *loopyRegister) nameTableY() byte {
	return byte((register.address >> 11) & 1)
}

func (register *loopyRegister) nameTableX() byte {
	return byte((register.address >> 10) & 1)
}

func (register *loopyRegister) coarseX() byte {
	return byte(register.address) & 0b11111
}

func (register *loopyRegister) coarseY() byte {
	return byte(register.address>>5) & 0b11111
}

func (register *loopyRegister) fineY() byte {
	return byte(register.address>>12) & 0b111
}

func (ppu *Ppu2c02) ppuCtrlWrite(value byte) {
	ppu.ppuControl.nameTableX = value & 0x01
	ppu.ppuControl.nameTableY = (value >> 1) & 1
	ppu.ppuControl.incrementMode = (value >> 2) & 1
	ppu.ppuControl.spritePatterTableAddress = (value >> 3) & 1
	ppu.ppuControl.backgroundPatternTableAddress = (value >> 4) & 1
	ppu.ppuControl.spriteSize = (value >> 5) & 1
	ppu.ppuControl.masterSlaveSelect = (value >> 6) & 1
	if (value>>7)&1 == 1 {
		ppu.ppuControl.generateNMIAtVBlank = true
	} else {
		ppu.ppuControl.generateNMIAtVBlank = false
	}
}

func (ppu *Ppu2c02) VBlank() bool {
	return ppu.ppuStatus.verticalBlankStarted
}
