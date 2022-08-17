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

func (cpu6502 *Cpu6502) Tick() (byte, cpu.CpuState) {
	var state cpu.CpuState

	if cpu6502.opCyclesLeft == 0 {
		registersCopy := *cpu6502.Registers()

		opcode := cpu6502.memory.Read(cpu6502.registers.Pc)
		instruction := cpu6502.instructions[opcode]
		cpu6502.opCyclesLeft = instruction.Cycles()

		if instruction.Method() == nil {
			msg := fmt.Errorf("opcode 0x%X not implemented", opcode)
			cpu6502.Stop()
			panic(msg)
		}

		operandAddress, operand, pageCrossed := cpu6502.evaluateOperandAddress(
			instruction.AddressMode(),
			cpu6502.registers.Pc+1,
		)
		step := cpu.OperationMethodArgument{
			instruction.AddressMode(),
			operandAddress,
		}

		state = cpu.CreateState(
			registersCopy,
			[3]byte{opcode, operand[0], operand[1]},
			instruction,
			step,
			cpu6502.cycle,
		)

		cpu6502.registers.Pc += types.Address(instruction.Size())

		opMightNeedExtraCycle := instruction.Method()(step)

		if pageCrossed && opMightNeedExtraCycle {
			cpu6502.opCyclesLeft++
		}

		cpu6502.cycle += uint32(cpu6502.opCyclesLeft)
	} else {
		state = cpu.CreateWaitingState()
	}

	cpu6502.opCyclesLeft--

	return cpu6502.opCyclesLeft, state
}

func (cpu6502 *Cpu6502) Stop() {
	cpu6502.debugger.Stop()
}
