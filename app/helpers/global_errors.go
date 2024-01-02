package helpers

import (
	errorsmod "cosmossdk.io/errors"
)

const codespace = "aura-global"

var (
	ErrOutOfGas               = errorsmod.Register(codespace, 1, "contract execution ran out of gas")
	ErrContractExecutionPanic = errorsmod.Register(codespace, 2, "contract execution panicked")
)
