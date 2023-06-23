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
	ErrSmartAccountAddress       = "Invalid address: %s"
	ErrAddressFromBech32         = "cannot convert string to address: %s"
	ErrSetSmartAccountStatus     = "cannot set smartaccount status: %s"
	ErrStoreSmartAccount         = "cannot store smartaccount value: %s"
	ErrInvalidTx                 = "invalid tx: %s"
	ErrSmartAccountCall          = "smart-account call fail: %s"
	ErrInvalidMsg                = "invalid smart-account message: %s"
	ErrNoSuchCodeID              = "code id not found: %d"
	ErrNilPubkey                 = "smart-account PublicKey must not be null"
	ErrAccountAlreadyExists      = "account already exists"
	ErrInstantiateDuplicate      = "instance with this code id, sender and label exists: try a different label"
)
