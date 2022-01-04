package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/types"
)

func (cpu6502 *Cpu6502) Init() {

}

func (cpu6502 *Cpu6502) Reset() {
	cpu6502.registers.Reset()
	cpu6502.cycle = 0

	// Read Reset Vector
	address := cpu6502.read16(cpu6502.Registers().Pc)
	cpu6502.registers.Pc = address
	cpu6502.cycle = 7
}

func (cpu6502 *Cpu6502) ResetToAddress(programCounter types.Address) {
	cpu6502.registers.Reset()
	cpu6502.registers.Pc = programCounter
	cpu6502.cycle = 7
}

func (cpu6502 *Cpu6502) Tick() byte {
	if cpu6502.opCyclesLeft > 0 {
		cpu6502.opCyclesLeft--
		return cpu6502.opCyclesLeft
	}

	registersCopy := *cpu6502.Registers()

	opcode := cpu6502.memory.Read(cpu6502.registers.Pc)
	instruction := cpu6502.instructions[opcode]
	cpu6502.opCyclesLeft = instruction.Cycles()

	if instruction.Method() == nil {
		msg := fmt.Errorf("opcode 0x%X not implemented", opcode)
		panic(msg)
	}

	operandAddress, operand, pageCrossed := cpu6502.evaluateOperandAddress(
		instruction.AddressMode(),
		cpu6502.registers.Pc+1,
	)
	step := OperationMethodArgument{
		instruction.AddressMode(),
		operandAddress,
	}
	if cpu6502.debug {
		cpu6502.logStep(registersCopy, opcode, operand, instruction, step, cpu6502.cycle)
	}

	cpu6502.registers.Pc += types.Address(instruction.Size())

	opMightNeedExtraCycle := instruction.Method()(step)

	if pageCrossed && opMightNeedExtraCycle {
		cpu6502.opCyclesLeft++
	}

	cpu6502.cycle += uint32(cpu6502.opCyclesLeft)

	return cpu6502.opCyclesLeft
}

func (cpu6502 *Cpu6502) logStep(registers cpu.Registers, opcode byte, operand [3]byte, instruction Instruction, step OperationMethodArgument, cpuCycle uint32) {
	//state := CreateStateFromCPU(*cpu6502)
	state := CreateState(
		registers,
		[3]byte{opcode, operand[0], operand[1]},
		instruction,
		step,
		cpuCycle,
	)

	cpu6502.Logger.Log(state)
}

func (cpu6502 *Cpu6502) Stop() {
	if cpu6502.debug {
		cpu6502.Logger.Close()
	}
}
