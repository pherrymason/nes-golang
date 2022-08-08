package nes

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func aDebugger() *Debugger {
	return CreateNesDebugger("", false, false)
}

func TestDebugger_isBreakpointTriggered_should_return_false(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x200

	debugger.AddBreakPoint(0x100)

	assert.False(t, debugger.isBreakpointTriggered())
}

func TestDebugger_isBreakpointTriggered_should_return_true(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x100

	debugger.AddBreakPoint(0x100)

	assert.True(t, debugger.isBreakpointTriggered())
}

func TestDebugger_shouldPauseEmulation_return_true_when_breakpoint_is_reached(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x100

	debugger.AddBreakPoint(0x100)

	assert.True(t, debugger.shouldPauseEmulation())
}

func TestDebugger_shouldPauseEmulation_return_true_until_resume_is_enabled(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x102

	debugger.AddBreakPoint(0x100)
	debugger.cpuStepByStepMode = true // We simulate breakpoint was enabled some instructions before

	assert.True(t, debugger.shouldPauseEmulation())
}

func TestDebugger_shouldPauseEmulation_when_cpu_is_paused_and_user_clicks_on_run_one_step_it_should_run_allow_emulation_until_next_cpu_operation_finishes(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x100

	debugger.AddBreakPoint(0x100)
	debugger.shouldPauseEmulation()       // This is to trigger breakpoint for first time
	debugger.RunOneCPUOperationAndPause() // Use then allows Cpu to run next instruction

	assert.False(t, debugger.shouldPauseEmulation(), "Should allow running one Cpu cycle.")

	debugger.oneCpuOperationRan()
	assert.True(t, debugger.shouldPauseEmulation(), "Should again stop emulation after running one Cpu cycle")
}

func TestDebugger_shouldPauseEmulation_after_pausing_due_breakpoint_and_user_clicking_on_run_one_step_it_should_pause_after_third_cpu_cycle(t *testing.T) {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	debugger := aDebugger()
	nes := CreateNes(cartridge, debugger)
	nes.Cpu.registers.Pc = 0x102

	debugger.AddBreakPoint(0x100)
	debugger.cpuStepByStepMode = true
	debugger.RunOneCPUOperationAndPause()

	assert.False(t, debugger.shouldPauseEmulation(), "Should allow running one Cpu cycle.")

	debugger.oneCpuOperationRan()
	assert.True(t, debugger.shouldPauseEmulation(), "Should again stop emulation after running one Cpu cycle")
}
