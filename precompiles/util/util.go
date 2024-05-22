package util

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v16/precompiles/common"
)

func EvmToAuraBigInt(amount *big.Int) *big.Int {
	return new(big.Int).Div(amount, big.NewInt(1e12))
}

func AuraToEvmBigInt(amount *big.Int) *big.Int {
	return new(big.Int).Mul(amount, big.NewInt(1e12))
}

func EvmToAuraInt(amount types.Int) types.Int {
	return amount.Quo(types.NewInt(1e12))
}

func AuraToEvmInt(amount types.Int) types.Int {
	return amount.Mul(types.NewInt(1e12))
}

func NewDecCoinsResponseEVM(amount types.DecCoins) []cmn.DecCoin {
	// Create a new output for each coin and add it to the output array.
	outputs := make([]cmn.DecCoin, len(amount))
	for i, coin := range amount {
		outputs[i] = cmn.DecCoin{
			Denom:     coin.Denom,
			Amount:    AuraToEvmBigInt(coin.Amount.TruncateInt().BigInt()),
			Precision: math.LegacyPrecision,
		}
	}
	return outputs
}
