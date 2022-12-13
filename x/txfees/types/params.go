package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"strings"
)

// Parameter store keys
var (
	KeyFeeDenom        = []byte("FeeDenom")
	KeyEpochIdentifier = []byte("EpochIdentifier")
)

// ParamKeyTable for tx fees module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(feeDenom string, epochIdentifier string) Params {
	return Params{
		FeeDenom:        feeDenom,
		EpochIdentifier: epochIdentifier,
	}
}

// default tx fees module parameters
func DefaultParams() Params {
	return Params{
		FeeDenom:        sdk.DefaultBondDenom,
		EpochIdentifier: "day",
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateFeeDenom(p.FeeDenom); err != nil {
		return err
	}
	if err := validateEpochIdentifier(p.EpochIdentifier); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeDenom, &p.FeeDenom, validateFeeDenom),
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateEpochIdentifier),
	}
}

func validateFeeDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("fee denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateEpochIdentifier(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("epoch identifier cannot be blank")
	}

	return nil
}
