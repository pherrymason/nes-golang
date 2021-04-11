package defs

// Address is an address representation
type Address uint16

// CreateAddress creates an Address
func CreateAddress(low byte, high byte) Address {
	return Address(uint16(low) + uint16(high)<<8)
}
