package keeper

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/aura-nw/aura/x/aura/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramSpace paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,

) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramSpace: paramSpace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the total set of aura parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of aura parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetMaxSupply return max supply of aura coin
func (k Keeper) GetMaxSupply(ctx sdk.Context) string {
	params := k.GetParams(ctx)
	return params.MaxSupply
}

// GetExcludeCirculatingAddr return list exclude address do not calculator for circulating
func (k Keeper) GetExcludeCirculatingAddr(ctx sdk.Context) []sdk.AccAddress {
	params := k.GetParams(ctx)
	excludeAddr := make([]sdk.AccAddress, 0, len(params.ExcludeCirculatingAddr))
	for _, addrBech32 := range params.ExcludeCirculatingAddr {
		addr, err := sdk.AccAddressFromBech32(addrBech32)
		if err != nil {
			panic(err)
		}
		excludeAddr = append(excludeAddr, addr)
	}

	return excludeAddr
}

// GetClaimDuration return claim reward duration
func (k Keeper) GetClaimDuration(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return int64(params.ClaimDuration)
}
