package helpers

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Sudo execute contract, recover from panic
//
// referenced from Juno' helpers:
// https://github.com/CosmosContracts/juno/blob/3a0bb93303772936362c93e3979d900a39cda68a/app/helpers/contracts.go
func SudoContract(k wasmtypes.ContractOpsKeeper, childCtx sdk.Context, contractAddr sdk.AccAddress, msgBz []byte, err *error) {
	// Recover from panic, return error
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			// Determine error associated with panic
			if isOutofGas, msg := IsOutOfGasError(recoveryError); isOutofGas {
				*err = ErrOutOfGas.Wrapf("%s", msg)
			} else {
				*err = ErrContractExecutionPanic.Wrapf("%s", recoveryError)
			}
		}
	}()

	// Execute contract with sudo
	_, *err = k.Sudo(childCtx, contractAddr, msgBz)
}

// Check if error is out of gas error
func IsOutOfGasError(err any) (bool, string) {
	switch e := err.(type) {
	case storetypes.ErrorOutOfGas:
		return true, e.Descriptor
	case storetypes.ErrorGasOverflow:
		return true, e.Descriptor
	default:
		return false, ""
	}
}
