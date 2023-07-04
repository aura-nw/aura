package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/aa module sentinel errors
var (
	ErrBadInstantiateMsg         = sdkerrors.Register(ModuleName, 1, "bad instantiate message")
	ErrInvalidPubKey             = sdkerrors.Register(ModuleName, 2, "invalid publickey")
	ErrAccountNotFoundForAddress = sdkerrors.Register(ModuleName, 3, "account not found for address")
	ErrInvalidAddress            = sdkerrors.Register(ModuleName, 4, "invalid address")
	ErrInvalidBench32            = sdkerrors.Register(ModuleName, 5, "invalid bench32")
	ErrInvalidTx                 = sdkerrors.Register(ModuleName, 6, "invalid tx")
	ErrInvalidMsg                = sdkerrors.Register(ModuleName, 7, "invalid smart-account messages")
	ErrNoSuchCodeID              = sdkerrors.Register(ModuleName, 8, "code id not found")
	ErrNilPubkey                 = sdkerrors.Register(ModuleName, 9, "nil pubkey")
	ErrAccountAlreadyExists      = sdkerrors.Register(ModuleName, 10, "account already exists")
	ErrInstantiateDuplicate      = sdkerrors.Register(ModuleName, 11, "instance with this code id, sender, salt and label exists")
	ErrInvalidCredentials        = sdkerrors.Register(ModuleName, 12, "invalid credentials")
	ErrInvalidCodeID             = sdkerrors.Register(ModuleName, 13, "invalid code id")
	ErrNotSupported              = sdkerrors.Register(ModuleName, 14, "not supported")
)
