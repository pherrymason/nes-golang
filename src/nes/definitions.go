package nes

// Address is an address representation
type Address = Word

// CreateAddress creates an Address
func CreateAddress(low byte, high byte) Address {
	return Address(uint16(low) + uint16(high)<<8)
}

// Word is an unsigned int 16 bits
type Word uint16

func (word *Word) LowNibble() byte {
	return byte(*word & 0x00FF)
}

func (word *Word) HighNibble() byte {
	return byte((*word & 0xFF00) >> 8)
}

func (word *Word) ToBytes() []byte {
	return []byte{word.HighNibble(), word.LowNibble()}
}

func LowNibble(word Word) byte {
	return byte(word & 0x00FF)
}

func HighNibble(word Word) byte {
	return byte((word & 0xFF00) >> 8)
}

// CreateWord creates a Word
func CreateWord(low byte, high byte) Word {
	return Word(uint16(low) + uint16(high)<<8)
}
