package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/aa module sentinel errors
var (
	ErrSample                    = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrWasmKeeper                = "nil wasm keeper"
	ErrBadInstantiateMsg         = "cannot instantiate contract: %v"
	ErrBadPublicKey              = "cannot convert publickey: %v"
	ErrSetPublickey              = "cannot set public key: %v"
	ErrAccountNotFoundForAddress = "account not found for smartcontract address: %s"
	ErrAddressFromBech32         = "cannot convert string to address: %s"
	ErrSetSmartAccountStatus     = "cannot set smartaccount status: %s"
	ErrStoreSmartAccount         = "cannot store smartaccount value: %s"
)
