package keeper

import (
	"github.com/aura-nw/aura/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastEpochVotesInfo(ctx sdk.Context) types.EpochVotesInfo {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.EpochVotesInfoKey)

	if bz == nil {
		return types.EpochVotesInfo{}
	}
	var epochVotesInfo types.EpochVotesInfo
	k.cdc.MustUnmarshal(bz, &epochVotesInfo)
	return epochVotesInfo
}

func (k Keeper) UpdateLastEpochVoteInfo(ctx sdk.Context, listVotes []types.ValidatorEpochVoteInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.EpochVotesInfo{Validators: listVotes})
	store.Set(types.EpochVotesInfoKey, bz)
}

func (k Keeper) ResetEpochVotesInfo(ctx sdk.Context) {
	k.UpdateLastEpochVoteInfo(ctx, []types.ValidatorEpochVoteInfo{})
}
