package nes

import "fmt"

/**
 * This Memory component handles the interfaces processors will use to
 * read and write through the communication bus.
 *
 * Two specific implementations to be used by CPU and PPU respectively
 * are defined.
 *
 * As Golang is not a OOP language, I've found that this is a good approach.
 * Another possibility would be to implement this in Nes space
 * and inject Read/Write methods into CPU and PPU
 */

// Nes Memory Map
// 	0x0000 - 0x00FF  ZeroPage
//	0x0100 - 0x01FF  Stack
//	0x0200 - 0x07FF  General Purpose RAM
//  0x0800 - 0x1FFF  Mirrors previous chunk of memory (0x0 - 0x7FF)
//  0x2000 - 0x2007  PPU registers
//  0x2008 - 0x4000  Mirrors PPU registers
//  0x4000 - 0x4020  I/O Registers
//  0x4020 - 0x5FFF  Expansion ROM
//  0x6000 - 0x7FFF  SRAM
//  0x8000 - 0xBFFF  GamePak prgROM lower bank
//  0xC000 - 0x10000 GamePak prgROM higher bank

// How I will structure all these stuff
// ram:     from 0x0000 -> 0x1FFF.
//          ZeroPage, Stack, General Purpose RAM and Mirroring from 0x07FF-0x1FFF
// gamePak: from 0x8000 -> 0x10000

const RAM_LOWER_ADDRESS = Address(0x0000)
const RAM_HIGHER_ADDRESS = Address(0x1FFF)
const RAM_LAST_REAL_ADDRESS = Address(0x07FF)

type Memory interface {
	// Peek Reads without side effects. Useful for debugging
	Peek(Address) byte
	Read(Address) byte
	Write(Address, byte)
}

type CPUMemory struct {
	ram     [0xFFFF + 1]byte
	gamePak *GamePak
	mapper  Mapper
	ppu     PPU
}

func CreateCPUMemory(ppu PPU, gamePak *GamePak) *CPUMemory {
	mapper := CreateMapper(gamePak)

	return &CPUMemory{gamePak: gamePak, mapper: mapper, ppu: ppu}
}

func (cm *CPUMemory) Peek(address Address) byte {
	return cm.read(address, false)
}

func (cm *CPUMemory) Read(address Address) byte {
	return cm.read(address, true)
}

func (cm *CPUMemory) read(address Address, readOnly bool) byte {
	if address <= RAM_HIGHER_ADDRESS {
		// Read with mirror after RAM_LAST_REAL_ADDRESS
		return cm.ram[address&RAM_LAST_REAL_ADDRESS]
	} else if address <= 0x3FFF {
		return cm.ppu.ReadRegister(address & 0x2007)
	} else if address >= GAMEPAK_LOW_RANGE {
		return cm.mapper.Read(address)
	}

	panic(fmt.Sprintf("reading from invalid address %X", address))
}

func (cm *CPUMemory) Write(address Address, value byte) {
	if address <= RAM_HIGHER_ADDRESS {
		cm.ram[address&RAM_LAST_REAL_ADDRESS] = value
	} else if address <= 0x3FFF {
		cm.ppu.WriteRegister(address&0x2007, value)
	} else if address >= GAMEPAK_ROM_LOWER_BANK_START {
		cm.mapper.Write(address, value)
	}
}
