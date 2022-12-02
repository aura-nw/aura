package keeper

import (
	"errors"
	"github.com/aura-nw/aura/x/mint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := k.GetParams(ctx)

	if epochIdentifier == params.EpochIdentifier {
		minter := k.GetMinter(ctx)

		// check over max supply
		maxSupplyString := k.GetMaxSupply(ctx)
		maxSupply, ok := sdk.NewIntFromString(maxSupplyString)
		if !ok {
			panic(errors.New("panic convert max supply string to bigInt"))
		}
		currentSupply := k.GetSupply(ctx, params.GetMintDenom())
		excludeAmount := k.GetExcludeCirculatingAmount(ctx, params.GetMintDenom())

		if currentSupply.LT(maxSupply) {
			// recalculate inflation rate
			totalStakingSupply := k.CustomStakingTokenSupply(ctx, excludeAmount.Amount)
			bondedRatio := k.CustomBondedRatio(ctx, excludeAmount.Amount)
			minter.Inflation = minter.NextInflationRate(params, bondedRatio)
			minter.EpochProvisions = minter.EpochProvision(params, totalStakingSupply)
			k.SetMinter(ctx, minter)

			// mint coins, update supply
			mintedCoin := minter.EpochReward(params)
			mintedCoins := sdk.NewCoins(mintedCoin)
			k.Logger(ctx).Info("AfterEpochEnd", "mintedCoin", mintedCoin)

			supplyNext := currentSupply.Add(mintedCoin.Amount)
			if supplyNext.GT(maxSupply) {
				mintedCoin.Amount = maxSupply.Sub(currentSupply)
				mintedCoins = sdk.NewCoins(mintedCoin)
			}
			err := k.MintCoins(ctx, mintedCoins)
			if err != nil {
				panic(err)
			}

			// send the minted coins to the fee collector account
			err = k.AddCollectedFees(ctx, mintedCoins)
			if err != nil {
				panic(err)
			}

			if mintedCoin.Amount.IsInt64() {
				defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeMint,
					sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
					sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
					sdk.NewAttribute(types.AttributeKeyEpochProvisions, minter.EpochProvisions.String()),
					sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
				),
			)

		} else {
			k.Logger(ctx).Info("Over the max supply", "currentSupply", currentSupply)
		}
	}
	return nil
}

// BeforeEpochStart is a hook which is executed before the start of an epoch. It is a no-op for mint module.
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// no-op
	return nil
}

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k Keeper
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}
