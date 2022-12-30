package keeper

import (
	"github.com/aura-nw/aura/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetListPreviousEpochVoteInfo(ctx sdk.Context) types.ListEpochVoteInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochVoteInfoKey)

	if bz == nil {
		panic("list previous epoch info not set")
	}

	var votes types.ListEpochVoteInfo
	k.cdc.MustUnmarshal(bz, &votes)
	return votes
}

func (k Keeper) UpdateListPreviousEpochVoteInfo(ctx sdk.Context, listVotes types.ListEpochVoteInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&listVotes)
	store.Set(types.EpochVoteInfoKey, bz)
}
