package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const DefaultMaxGas = 2_000_000

// Parameter store keys
var (
	WhitelistCodeID = []byte("WhitelistCodeID")
	MaxGasExecute   = []byte("MaxGasExecute")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(whitelist []*CodeID, limitGas uint64) Params {
	return Params{
		WhitelistCodeID: whitelist,
		MaxGasExecute:   limitGas,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	empty := make([]*CodeID, 0)

	return NewParams(empty, DefaultMaxGas)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// code_id whitelist indicates which contract can be initialized as smart account
		// using gov proposal for updates
		paramtypes.NewParamSetPair(WhitelistCodeID, &p.WhitelistCodeID, validateWhitelistCodeID),
		// max_gas_query limits the amount of gas that the validation query can use
		paramtypes.NewParamSetPair(MaxGasExecute, &p.MaxGasExecute, validateMaxGasExecute),
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

func validateMaxGasExecute(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("zero max gas execute")
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

	// validate max gas execute
	err = validateMaxGasExecute(p.MaxGasExecute)
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
