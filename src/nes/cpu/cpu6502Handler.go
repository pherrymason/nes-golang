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
	cpu.instructionCycle = 0

	// Read Reset Vector
	address := cpu.bus.Read16(0xFFFC)
	cpu.registers.Pc = defs.Address(address)
}

func (cpu *Cpu6502) ResetToAddress(programCounter defs.Address) {
	cpu.registers.reset()
	cpu.registers.Pc = programCounter
	cpu.cycle = 7
}

func (cpu *Cpu6502) Tick() byte {
	if cpu.instructionCycle == 0 {
		if cpu.debug {
			cpu.logStep()
		}

		opcode := cpu.Read(cpu.registers.Pc)
		instruction := cpu.instructions[opcode]
		cpu.instructionCycle = instruction.Cycles()

		if instruction.Method() == nil {
			msg := fmt.Errorf("opcode 0x%X not implemented", opcode)
			panic(msg)
		}

		operandAddress, pageCrossed := cpu.evaluateOperandAddress(instruction.AddressMode(), cpu.registers.Pc+1)

		cpu.registers.Pc += defs.Address(instruction.Size())

		step := defs.OperationMethodArgument{
			instruction.AddressMode(),
			operandAddress,
		}
		opMightNeedExtraCycle := instruction.Method()(step)

		if pageCrossed && opMightNeedExtraCycle {
			cpu.instructionCycle++
		}
		cpu.cycle += uint16(cpu.instructionCycle)
	}

	cpu.instructionCycle--

	return cpu.instructionCycle
}

func (cpu *Cpu6502) logStep() {
	state := CreateState(*cpu)

	cpu.Logger.Log(state)
}
