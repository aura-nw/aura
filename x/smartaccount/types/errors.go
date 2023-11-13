package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/aa module sentinel errors
var (
	ErrBadInstantiateMsg         = errorsmod.Register(ModuleName, 1, "bad instantiate message")
	ErrInvalidPubKey             = errorsmod.Register(ModuleName, 2, "invalid publickey")
	ErrAccountNotFoundForAddress = errorsmod.Register(ModuleName, 3, "account not found for address")
	ErrInvalidAddress            = errorsmod.Register(ModuleName, 4, "invalid address")
	ErrInvalidBench32            = errorsmod.Register(ModuleName, 5, "invalid bench32")
	ErrInvalidTx                 = errorsmod.Register(ModuleName, 6, "invalid tx")
	ErrInvalidMsg                = errorsmod.Register(ModuleName, 7, "invalid smart-account messages")
	ErrNoSuchCodeID              = errorsmod.Register(ModuleName, 8, "code id not found")
	ErrNilPubkey                 = errorsmod.Register(ModuleName, 9, "nil pubkey")
	ErrAccountAlreadyExists      = errorsmod.Register(ModuleName, 10, "account already exists")
	ErrInstantiateDuplicate      = errorsmod.Register(ModuleName, 11, "instance with this code id, sender, salt and label exists")
	ErrInvalidCredentials        = errorsmod.Register(ModuleName, 12, "invalid credentials")
	ErrInvalidCodeID             = errorsmod.Register(ModuleName, 13, "invalid code id")
	ErrNotSupported              = errorsmod.Register(ModuleName, 14, "not supported")
	ErrNotAllowedMsg             = errorsmod.Register(ModuleName, 15, "not allowed message")
)
