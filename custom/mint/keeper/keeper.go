package keeper

import (
	custommint "github.com/aura-nw/aura/custom/mint/types"
	aurakeeper "github.com/aura-nw/aura/x/aura/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the mint store
type Keeper struct {
	mintkeeper.Keeper

	bankKeeper custommint.BankKeeper
	auraKeeper aurakeeper.Keeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk custommint.StakingKeeper, ak custommint.AccountKeeper, bk custommint.BankKeeper,
	auraKeeper aurakeeper.Keeper, feeCollectorName string,
) Keeper {
	return Keeper{
		Keeper:     mintkeeper.NewKeeper(cdc, key, paramSpace, sk, ak, bk, feeCollectorName),
		bankKeeper: bk,
		auraKeeper: auraKeeper,
	}
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}

func (k Keeper) GetMaxSupply(ctx sdk.Context) string {
	return k.auraKeeper.GetMaxSupply(ctx)
}
