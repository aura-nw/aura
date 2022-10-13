package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuraKeeper interface {
	GetClaimDuration(ctx sdk.Context) int64
}
