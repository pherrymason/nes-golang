package nes

import (
	"fmt"
)

func (cpu *Cpu6502) Init() {

}

func (cpu *Cpu6502) Reset() {
	cpu.registers.Reset()
	cpu.cycle = 0

	// Read Reset Vector
	address := cpu.read16(0xFFFC)
	cpu.registers.Pc = Address(address)
	cpu.cycle = 7
}

func (cpu *Cpu6502) ResetToAddress(programCounter Address) {
	cpu.registers.Reset()
	cpu.registers.Pc = programCounter
	cpu.cycle = 7
}

func (cpu *Cpu6502) Tick() byte {
	if cpu.opCyclesLeft == 0 {
		registersCopy := *cpu.Registers()

		opcode := cpu.memory.Read(cpu.registers.Pc)
		instruction := cpu.instructions[opcode]
		cpu.opCyclesLeft = instruction.Cycles()

		if instruction.Method() == nil {
			msg := fmt.Errorf("opcode 0x%X not implemented", opcode)
			panic(msg)
		}

	operandAddress, operand, pageCrossed := cpu.evaluateOperandAddress(instruction.AddressMode(), cpu.registers.Pc+1)
	cpu.registers.Pc += Address(instruction.Size())

		step := OperationMethodArgument{
			instruction.AddressMode(),
			operandAddress,
		}
		opMightNeedExtraCycle := instruction.Method()(step)

		if pageCrossed && opMightNeedExtraCycle {
			cpu.opCyclesLeft++
		}

		if cpu.debug {
			cpu.logStep(registersCopy, opcode, instruction, step, cpu.cycle)
		}

		cpu.cycle += uint32(cpu.opCyclesLeft)
	}

	cpu.opCyclesLeft--

	return cpu.opCyclesLeft
}

func (cpu *Cpu6502) logStep(registers Cpu6502Registers, opcode byte, instruction Instruction, step OperationMethodArgument, cpuCycle uint32) {
	//state := CreateStateFromCPU(*cpu)
	state := CreateState(registers, opcode, instruction, step, cpuCycle)

	cpu.Logger.Log(state)
}

func (cpu *Cpu6502) Stop() {
	if cpu.debug {
		cpu.Logger.Close()
	}
}
