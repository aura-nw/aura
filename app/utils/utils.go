package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	AuraExponent = 6
	BaseCoinUnit = "ueaura"
)

var (
	// DevnetChainID defines the Aura chain ID for devnet
	DevnetChainID = "aura-testnet"

	// SerenityChainID defines the Aura chain ID for serenity testnet
	SerenityChainID = "serenity-testnet"

	// EuphoriaChainID defines the Aura chain ID for euphoria testnet
	EuphoriaChainID = "euphoria"
)

// IsDevnet returns true if the chain-id has the Aura devnet chain prefix.
func IsDevnet(chainID string) bool {
	return strings.HasPrefix(chainID, DevnetChainID)
}

// IsSerenity returns true if the chain-id has the Aura serenity network chain prefix.
func IsSerenity(chainID string) bool {
	return strings.HasPrefix(chainID, SerenityChainID)
}

// IsEuphoria returns true if the chain-id has the Aura euphoria network chain prefix.
func IsEuphoria(chainID string) bool {
	return strings.HasPrefix(chainID, EuphoriaChainID)
}

// RegisterDenoms registers token denoms.
func RegisterDenoms() {
	err := sdk.RegisterDenom(BaseCoinUnit, sdk.NewDecWithPrec(1, AuraExponent))
	if err != nil {
		panic(err)
	}
}
