package distribution

import (
	"bytes"
	customdistkeeper "github.com/aura-nw/aura/x/distribution/keeper"
	"github.com/aura-nw/aura/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k customdistkeeper.Keeper) {
	defer telemetry.ModuleMeasureSince(disttypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.Logger(ctx).Error("customdistr/beginblocker=======================")
	lastEpochVotesInfo := k.GetLastEpochVotesInfo(ctx)
	previousProposer := k.GetPreviousProposerConsAddr(ctx)

	updatedVotes := []types.ValidatorEpochVoteInfo{}

	previousVotes := req.LastCommitInfo.GetVotes()

	for _, previousVote := range previousVotes {
		for _, lastVote := range lastEpochVotesInfo.Validators {
			valAddr, err := sdk.ValAddressFromBech32(lastVote.ValidatorAddress)
			if err != nil {
				panic(err)
			}
			if bytes.Equal(previousVote.Validator.Address, valAddr.Bytes()) {
				proposerBlocks := lastVote.ProposerBlocks
				if bytes.Equal(previousVote.Validator.Address, previousProposer.Bytes()) {
					proposerBlocks += 1
				}
				updateVote := types.ValidatorEpochVoteInfo{
					ValidatorAddress: lastVote.ValidatorAddress,
					ActiveBlocks:     lastVote.ActiveBlocks + 1,
					AccPower:         lastVote.AccPower + previousVote.Validator.Power,
					ProposerBlocks:   proposerBlocks,
				}
				updatedVotes = append(updatedVotes, updateVote)
			}
		}
	}

	k.UpdateLastEpochVoteInfo(ctx, updatedVotes)

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}
