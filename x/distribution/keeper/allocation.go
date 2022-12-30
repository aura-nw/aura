package keeper

import (
	"github.com/aura-nw/aura/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens handles distribution of the collected fees on one epoch
func (k Keeper) AllocateTokens(
	ctx sdk.Context,
	numberBlocks int64,
	totalPreviousPower int64,
	bondedEpochVotes []types.EpochVoteInfo,
) {
	logger := k.Logger(ctx)

	// get total fees in epoch
	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	logger.Info("feesCollected = ", feesCollected)

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, disttypes.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		logger.Info("No previous power")
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return
	}

	// get reward config
	baseProposerReward := k.GetBaseProposerReward(ctx)
	bonusProposerReward := k.GetBonusProposerReward(ctx)
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward)

	communityTax := k.GetCommunityTax(ctx)
	feeCommunityTax := feesCollected.MulDecTruncate(communityTax)
	feePool.CommunityPool = feePool.CommunityPool.Add(feeCommunityTax...)
	feesCollectedAfterTax := feesCollected.Sub(feeCommunityTax)

	var parts int64
	for _, voteInfo := range bondedEpochVotes {
		parts += voteInfo.ActiveBlocks - voteInfo.ProposerBlocks
	}

	dec := sdk.NewDec(parts).Add(sdk.NewDec(numberBlocks).Mul(proposerMultiplier.Add(sdk.OneDec())))
	rewardUnit := feesCollectedAfterTax.MulDecTruncate(dec)

	remaining := feesCollectedAfterTax

	for _, voteInfo := range bondedEpochVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, voteInfo.Validator.Address)

		// Rewards for proposer
		proposerReward := rewardUnit.MulDecTruncate(proposerMultiplier.Add(sdk.OneDec())).MulDecTruncate(sdk.NewDec(voteInfo.ProposerBlocks))

		//  Rewards for non-proposer
		nonProposerReward := rewardUnit.MulDec(sdk.NewDec(voteInfo.ActiveBlocks - voteInfo.ProposerBlocks))

		reward := proposerReward.Add(nonProposerReward...)

		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
}

func (k Keeper) AllocateTokensToValidatorAndDelegator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	// update current commission
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			disttypes.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(disttypes.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			disttypes.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(disttypes.AttributeKeyValidator, val.GetOperator().String()),
		),
	)

	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}
