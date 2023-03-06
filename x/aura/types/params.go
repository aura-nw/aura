package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"strings"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	MaxSupply              = []byte("MaxSupply")
	ExcludeCirculatingAddr = []byte("ExcludeCirculatingAddr")
)

const (
	MinMaxSupply                      = 1_000_000
	LimitLengthExcludeCirculatingAddr = 10
)

// Regex using check string is number
var digitCheck = regexp.MustCompile(`^[0-9]+$`)

// ParamTable for aura module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(maxSupply string, excludeCirculatingAddr []string) Params {
	return Params{
		MaxSupply:              maxSupply,
		ExcludeCirculatingAddr: excludeCirculatingAddr,
	}
}

// default aura module parameters
func DefaultParams() Params {
	return Params{
		MaxSupply:              "1000000000000000000000000000",
		ExcludeCirculatingAddr: []string{},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}

	if err := validateExcludeCirculatingAddr(p.ExcludeCirculatingAddr); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(MaxSupply, &p.MaxSupply, validateMaxSupply),
		paramtypes.NewParamSetPair(ExcludeCirculatingAddr, &p.ExcludeCirculatingAddr, validateExcludeCirculatingAddr),
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

	vi, err := sdk.NewDecFromStr(v)
	if err != nil {
		return fmt.Errorf("can not parse max supply to Dec, value=%s", v)
	}

	if vi.LT(sdk.NewDec(MinMaxSupply)) {
		return fmt.Errorf("required max supply greater than %d", MinMaxSupply)
	}

	return nil
}

func validateExcludeCirculatingAddr(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) > LimitLengthExcludeCirculatingAddr {
		return errors.New("len of exclude exclude circulating address reach limit")
	}

	for _, addBech32 := range v {
		if strings.TrimSpace(addBech32) == "" {
			return errors.New("exclude circulating address can not contain blank")
		}
	}

	if checkDuplicate(v) {
		return errors.New("duplicated address in exclude circulating address")
	}

	return nil
}

func checkDuplicate(ss []string) bool {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	return !(len(ss) == len(m))
}
