package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// At end of epoch, allocate token rewards to all validators contributed

	if epochIdentifier == "day" {
		epochVotesInfo := k.GetLastEpochVotesInfo(ctx)
		k.AllocateTokens(ctx, epochVotesInfo)
	}

	return nil
}

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// Reset Epoch vote info when starting new epoch
	k.ResetEpochVotesInfo(ctx)
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
