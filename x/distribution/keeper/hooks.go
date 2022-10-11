package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Hooks struct {
	k     Keeper
	hooks keeper.Hooks
}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k, k.Keeper.Hooks()} }

func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.hooks.AfterValidatorCreated(ctx, valAddr)
}

func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.hooks.BeforeValidatorModified(ctx, valAddr)
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
}

func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.hooks.AfterValidatorBonded(ctx, consAddr, valAddr)
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.hooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
}

func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
}

func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	// check condition claim reward
	claimInfo, err := h.k.IsClaimReward(ctx, delAddr)
	if err != nil {
		panic(err)
	}

	h.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	h.k.SetClaimInfo(ctx, delAddr, claimInfo)
}

func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}

func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
}

func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
}

var _ stakingtypes.StakingHooks = Hooks{}
