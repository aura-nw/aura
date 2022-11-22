package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k Keeper
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	//TODO implement me
	panic("implement me")
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	//TODO implement me
	panic("implement me")
}
