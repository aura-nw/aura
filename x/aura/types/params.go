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
	MaxSupply              = []byte("MaxSupply")
	ExcludeCirculatingAddr = []byte("ExcludeCirculatingAddr")
	ClaimDuration          = []byte("ClaimDuration")
)

// Regex using check string is number
var digitCheck = regexp.MustCompile(`^[0-9]+$`)

// ParamTable for aura module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(maxSupply string, excludeCirculatingAddr []string, claimDuration uint32) Params {
	return Params{
		MaxSupply:              maxSupply,
		ExcludeCirculatingAddr: excludeCirculatingAddr,
		ClaimDuration:          claimDuration,
	}
}

// default aura module parameters
func DefaultParams() Params {
	return Params{
		MaxSupply:              "1000000000000000000000000000",
		ExcludeCirculatingAddr: []string{},
		ClaimDuration:          86400000, // default 1 day
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}

	if err := validateClaimDuration(p.ClaimDuration); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(MaxSupply, &p.MaxSupply, validateMaxSupply),
		paramtypes.NewParamSetPair(ExcludeCirculatingAddr, &p.ExcludeCirculatingAddr, validateExcludeCirculatingAddr),
		paramtypes.NewParamSetPair(ClaimDuration, &p.ClaimDuration, validateClaimDuration),
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

func validateExcludeCirculatingAddr(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, addBech32 := range v {
		if strings.TrimSpace(addBech32) == "" {
			return errors.New("exclude circulating address can not contain blank")
		}
	}
	return nil
}

func validateClaimDuration(i interface{}) error {
	return nil
}
