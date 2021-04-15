package cpu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

func (cpu *Cpu6502) Init() {
	cpu.initInstructionsTable()
	cpu.initAddressModeEvaluators()
}

func (cpu *Cpu6502) Reset() {
	cpu.registers.reset()

	// Read Reset Vector
	address := cpu.bus.Read16(0xFFFC)
	cpu.registers.Pc = defs.Address(address)
}

func (cpu *Cpu6502) ResetToAddress(programCounter defs.Address) {
	cpu.registers.reset()
	cpu.registers.Pc = programCounter
}

func (cpu *Cpu6502) Tick() {

	// Read opcode
	if cpu.debug {
		cpu.logStep()
	}

	opcode := cpu.Read(cpu.registers.Pc)
	//cpu.registers.Pc++

	instruction := cpu.instructions[opcode]
	if instruction.Method() == nil {
		msg := fmt.Errorf("opcode 0x%X not implemented", opcode)
		panic(msg)
	}

	operandAddress := cpu.evaluateOperandAddress(instruction.AddressMode(), cpu.registers.Pc+1)

	cpu.registers.Pc += defs.Address(instruction.Size())

	step := defs.InfoStep{
		instruction.AddressMode(),
		operandAddress,
	}
	instruction.Method()(step)

	// -analyze opcode:
	//	-address mode
	//  -get operand
	//  - update PC accordingly
	//  - run InfoStep

}

func (cpu *Cpu6502) logStep() {
	state := CreateState(*cpu)

	cpu.Logger.Log(state)

	/*

		pc := cpu.registers.Pc
		opcode := cpu.Read(pc)
		pc++
		instruction := cpu.instructions[opcode]
		if instruction.Method() == nil {
			msg := fmt.Errorf("Error: Opcode 0x%X not implemented", opcode)
			panic(msg)
		}

		evaluatedAddress := cpu.evaluateOperandAddress(instruction.AddressMode(), pc)

		var msg string
		msg += fmt.Sprintf("%X", (pc-1)) + "  "
		msg += fmt.Sprintf("%X ", opcode) + " "

		for i := byte(0); i < (instruction.Size() - 1); i++ {
			msg += fmt.Sprintf("%X ", cpu.Read(pc+defs.Address(i)))
		}

		for i := len(msg); i <= 16; i++ {
			msg += " "
		}

		msg += instruction.Name() + " "

		if instruction.AddressMode() == defs.Immediate {
			msg += "#"
		} else {
			msg += fmt.Sprintf("$%X", evaluatedAddress)
		}

		for i := len(msg); i <= 48; i++ {
			msg += " "
		}

		msg += fmt.Sprintf(
			"A:%X X:%X Y:%X P:%X SP:%X PPU:___,___ CYC:%d",
			cpu.registers.A,
			cpu.registers.X,
			cpu.registers.Y,
			cpu.registers.Status,
			cpu.registers.Sp,
			0,
		)

		cpu.Logger.Log(msg)
	*/
}
