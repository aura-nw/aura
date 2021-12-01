package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	MaxSupply = []byte("MaxSupply")
)

// Regex using check string is number
var digitCheck = regexp.MustCompile(`^[0-9]+$`)

// ParamTable for aura module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(maxSupply string) Params {
	return Params{
		MaxSupply: maxSupply,
	}
}

// default aura module parameters
func DefaultParams() Params {
	return Params{
		MaxSupply: "1000000000000000000000000000",
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(MaxSupply, &p.MaxSupply, validateMaxSupply),
	}
}

func validateMaxSupply(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("max supply cannot be blank")
	}

	if !digitCheck.MatchString(strings.TrimSpace(v)) {
		return errors.New("invalid max supply parameter, expected string as number")
	}

	return nil
}
