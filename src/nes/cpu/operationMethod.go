package cpu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

type OperationMethodArgument struct {
	AddressMode    AddressMode
	OperandAddress types.Address
}

type OperationMethod func(info OperationMethodArgument) bool
