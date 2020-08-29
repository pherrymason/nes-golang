package nes

// Word is an unsigned int 16 bits
type Word uint16

func (word *Word) low() byte {
	return byte(*word & 0x00FF)
}

func (word *Word) high() byte {
	return byte((*word & 0xFF00) >> 8)
}

// CreateWord creates a Word
func CreateWord(low byte, high byte) Word {
	return Word(uint16(low) + uint16(high)<<8)
}

// Address is an address representation
type Address uint16

// CreateAddress creates an Address
func CreateAddress(low byte, high byte) Address {
	return Address(uint16(low) + uint16(high)<<8)
}
