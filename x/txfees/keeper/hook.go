package keeper

import (
	"github.com/aura-nw/aura/utils"
	epochstypes "github.com/aura-nw/aura/x/epochs/types"
	"github.com/aura-nw/aura/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Hooks struct {
	k Keeper
}

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := k.GetParams(ctx)
	if epochIdentifier == params.EpochIdentifier {
		feeAddress := k.accountKeeper.GetModuleAddress(types.TxFeeCollectorName)
		baseDenomCoins := sdk.NewCoins(k.bankKeeper.GetBalance(ctx, feeAddress, params.FeeDenom))
		if !baseDenomCoins.IsZero() {
			utils.ApplyFuncIfNoError(ctx, func(cacheCtx sdk.Context) error {
				err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.TxFeeCollectorName, types.FeeCollectorName, baseDenomCoins)
				return err
			})
		}
	}
	return nil
}

var _ epochstypes.EpochHooks = Hooks{}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}
