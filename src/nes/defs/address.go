package defs

// Address is an address representation
type Address = Word

// CreateAddress creates an Address
func CreateAddress(low byte, high byte) Address {
	return Address(uint16(low) + uint16(high)<<8)
}
