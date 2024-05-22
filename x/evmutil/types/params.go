package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamKeyTable for evmutil module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value
// pairs pairs of the evmutil module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// NewParams returns new evmutil module Params.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns the default parameters for evmutil.
func DefaultParams() Params {
	return NewParams()
}

// Validate returns an error if the Params is invalid.
func (p *Params) Validate() error {
	return nil
}
