package defs

// AddressMode is an enum of the available Addressing Modes in this cpu
type AddressMode int

const (
	Implicit AddressMode = iota
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteXIndexed
	AbsoluteYIndexed
	Indirect
	IndirectX
	IndirectY
	Relative
)
