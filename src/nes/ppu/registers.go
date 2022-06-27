package ppu

import "github.com/raulferras/nes-golang/src/nes/types"

type Registers struct {
	ctrl   byte // Controls PPU operation
	mask   byte // Controls the rendering of sprites and backgrounds
	status byte // Reflects state of various functions inside PPU

	scrollX     byte // Changes scroll position
	scrollY     byte // Changes scroll position
	scrollLatch byte // Controls which scroll needs to be written

	oamAddr      byte
	ppuAddr      types.Address
	addressLatch byte
	readBuffer   byte
}

type CtrlFlag int

const (
	baseNameTableAddress0 CtrlFlag = iota // Most significant bit of scrolling coordinates (X)
	baseNameTableAddress1                 // Most significant bit of scrolling coordinates (Y)
	incrementMode
	spritePatternTableAddress
	backgroundPatternTableAddress
	spriteSize
	masterSlaveSelect
	generateNMIAtVBlank
)

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

type STATUSFlag int

const (
	spriteOverflow       = 5
	sprite0Hit           = 6
	verticalBlankStarted = 7
)
