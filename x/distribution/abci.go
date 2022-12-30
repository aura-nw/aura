package distribution

import (
	customdistkeeper "github.com/aura-nw/aura/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k customdistkeeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

}
