package nes

const CONTROLLER_A byte = 0x80
const CONTROLLER_B byte = 0x40
const CONTROLLER_SELECT byte = 0x20
const CONTROLLER_START byte = 0x10
const CONTROLLER_ARROW_UP byte = 0x08
const CONTROLLER_ARROW_DOWN byte = 0x04
const CONTROLLER_ARROW_LEFT byte = 0x02
const CONTROLLER_ARROW_RIGHT byte = 0x01

type ControllerState struct {
	A      bool
	B      bool
	Select bool
	Start  bool
	Up     bool
	Down   bool
	Left   bool
	Right  bool
}

func (state *ControllerState) value() byte {
	value := byte(0)

	if state.A {
		value |= CONTROLLER_A
	}
	if state.B {
		value |= CONTROLLER_B
	}
	if state.Select {
		value |= CONTROLLER_SELECT
	}
	if state.Start {
		value |= CONTROLLER_START
	}
	if state.Up {
		value |= CONTROLLER_ARROW_UP
	}
	if state.Down {
		value |= CONTROLLER_ARROW_DOWN
	}
	if state.Left {
		value |= CONTROLLER_ARROW_LEFT
	}
	if state.Right {
		value |= CONTROLLER_ARROW_RIGHT
	}

	return value
}
