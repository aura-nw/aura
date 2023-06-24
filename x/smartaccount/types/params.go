package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	WhitelistCodeID = []byte("WhitelistCodeID")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(whitelist []*CodeID) Params {
	return Params{
		WhitelistCodeID: whitelist,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	empty := make([]*CodeID, 0)

	return NewParams(empty)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// code_id whitelist indicates which contract can be initialized as smart account
		// using gov proposal for updates
		paramtypes.NewParamSetPair(WhitelistCodeID, &p.WhitelistCodeID, validateWhitelistCodeID),
	}
}
func validateWhitelistCodeID(i interface{}) error {
	v, ok := i.([]*CodeID)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// not allowed duplicate code_id in whitelist
	visited := make(map[uint64]bool, 0)
	for _, codeID := range v {
		if visited[codeID.CodeID] {
			return fmt.Errorf("duplicate code_id %d in whitelist_code_id", codeID.CodeID)
		} else {
			visited[codeID.CodeID] = true
		}
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	// validate whitelist_code_id param
	err := validateWhitelistCodeID(p.WhitelistCodeID)
	if err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
