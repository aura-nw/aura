package types

import (
	bytes "bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

///////////////
// EVM -> Cosmos SDK
///////////////

// NewConversionPair returns a new ConversionPair.
func NewConversionPair(address InternalEVMAddress, denom string) ConversionPair {
	return ConversionPair{
		KavaERC20Address: address.Address.Bytes(),
		Denom:            denom,
	}
}

// GetAddress returns the InternalEVMAddress of the Kava ERC20 address.
func (pair ConversionPair) GetAddress() InternalEVMAddress {
	return NewInternalEVMAddress(common.BytesToAddress(pair.KavaERC20Address))
}

// Validate returns an error if the ConversionPair is invalid.
func (pair ConversionPair) Validate() error {
	if err := sdk.ValidateDenom(pair.Denom); err != nil {
		return fmt.Errorf("conversion pair denom invalid: %v", err)
	}

	if len(pair.KavaERC20Address) != common.AddressLength {
		return fmt.Errorf("address length is %v but expected %v", len(pair.KavaERC20Address), common.AddressLength)
	}

	if bytes.Equal(pair.KavaERC20Address, common.Address{}.Bytes()) {
		return fmt.Errorf("address cannot be zero value %v", hex.EncodeToString(pair.KavaERC20Address))
	}

	return nil
}

// ConversionPairs defines a slice of ConversionPair.
type ConversionPairs []ConversionPair

// NewConversionPairs returns ConversionPairs from the provided values.
func NewConversionPairs(pairs ...ConversionPair) ConversionPairs {
	return ConversionPairs(pairs)
}

func (pairs ConversionPairs) Validate() error {
	// Check for duplicates for both addrs and denoms
	addrs := map[string]bool{}
	denoms := map[string]bool{}

	for _, pair := range pairs {
		if addrs[hex.EncodeToString(pair.KavaERC20Address)] {
			return fmt.Errorf(
				"found duplicate enabled conversion pair internal ERC20 address %s",
				hex.EncodeToString(pair.KavaERC20Address),
			)
		}

		if denoms[pair.Denom] {
			return fmt.Errorf(
				"found duplicate enabled conversion pair denom %s",
				pair.Denom,
			)
		}

		if err := pair.Validate(); err != nil {
			return err
		}

		addrs[hex.EncodeToString(pair.KavaERC20Address)] = true
		denoms[pair.Denom] = true
	}

	return nil
}

// validateConversionPairs validates an interface as ConversionPairs
func validateConversionPairs(i interface{}) error {
	pairs, ok := i.(ConversionPairs)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return pairs.Validate()
}

///////////////
// Cosmos SDK -> EVM
///////////////

// NewDeployedCosmosCoinContract returns a new DeployedCosmosCoinContract
func NewDeployedCosmosCoinContract(denom string, address InternalEVMAddress) DeployedCosmosCoinContract {
	return DeployedCosmosCoinContract{
		CosmosDenom: denom,
		Address:     &address,
	}
}

// NewAllowedCosmosCoinERC20Token returns an AllowedCosmosCoinERC20Token
func NewAllowedCosmosCoinERC20Token(
	cosmosDenom, name, symbol string,
	decimal uint32,
) AllowedCosmosCoinERC20Token {
	return AllowedCosmosCoinERC20Token{
		CosmosDenom: cosmosDenom,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimal,
	}
}

// Validate validates the fields of a single AllowedCosmosCoinERC20Token
func (token AllowedCosmosCoinERC20Token) Validate() error {
	// disallow empty string fields
	if err := sdk.ValidateDenom(token.CosmosDenom); err != nil {
		return fmt.Errorf("allowed cosmos coin erc20 token's sdk denom is invalid: %v", err)
	}

	if token.Name == "" {
		return errors.New("allowed cosmos coin erc20 token's name cannot be empty")
	}

	if token.Symbol == "" {
		return errors.New("allowed cosmos coin erc20 token's symbol cannot be empty")
	}

	// ensure decimals will properly cast to uint8 of erc20 spec
	if token.Decimals > math.MaxUint8 {
		return fmt.Errorf("allowed cosmos coin erc20 token's decimals must be less than 256, found %d", token.Decimals)
	}

	return nil
}

// AllowedCosmosCoinERC20Tokens defines a slice of AllowedCosmosCoinERC20Token
type AllowedCosmosCoinERC20Tokens []AllowedCosmosCoinERC20Token

// NewAllowedCosmosCoinERC20Tokens returns AllowedCosmosCoinERC20Tokens from the provided values.
func NewAllowedCosmosCoinERC20Tokens(pairs ...AllowedCosmosCoinERC20Token) AllowedCosmosCoinERC20Tokens {
	return AllowedCosmosCoinERC20Tokens(pairs)
}

// Validate checks that all containing tokens are valid and that there are
// no duplicate denoms or symbols.
func (tokens AllowedCosmosCoinERC20Tokens) Validate() error {
	// Disallow multiple instances of a single sdk_denom or evm symbol
	denoms := make(map[string]struct{}, len(tokens))
	symbols := make(map[string]struct{}, len(tokens))

	for i, t := range tokens {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("invalid token at index %d: %s", i, err)
		}

		if _, found := denoms[t.CosmosDenom]; found {
			return fmt.Errorf("found duplicate token with sdk denom %s", t.CosmosDenom)
		}
		if _, found := symbols[t.Symbol]; found {
			return fmt.Errorf("found duplicate token with symbol %s", t.Symbol)
		}

		denoms[t.CosmosDenom] = struct{}{}
		symbols[t.Symbol] = struct{}{}
	}

	return nil
}

// validateAllowedCosmosCoinERC20Tokens validates an interface as AllowedCosmosCoinERC20Tokens
func validateAllowedCosmosCoinERC20Tokens(i interface{}) error {
	pairs, ok := i.(AllowedCosmosCoinERC20Tokens)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return pairs.Validate()
}
