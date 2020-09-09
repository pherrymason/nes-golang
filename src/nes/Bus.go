package nes

type Bus struct {
	ram *RAM
}

func (bus *Bus) read(address Address) byte {
	return bus.ram.read(address)
}

func (bus *Bus) read16(address Address) Word {
	low := bus.read(address)
	high := bus.read(address + 1)

	return CreateWord(low, high)
}

func (bus *Bus) read16Bugged(address Address) Word {
	lsb := address
	msb := (lsb & 0xFF00) | (lsb & 0xFF) + 1

	low := bus.ram.read(lsb)
	high := bus.ram.read(msb)

	return CreateWord(low, high)
}

func (bus *Bus) write(address Address, value byte) {
	bus.ram.write(address, value)
}

func CreateBus(ram *RAM) Bus {
	return Bus{ram}
}
