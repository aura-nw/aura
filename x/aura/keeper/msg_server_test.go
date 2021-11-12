package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/aura-nw/aura/testutil/keeper"
	"github.com/aura-nw/aura/x/aura/keeper"
	"github.com/aura-nw/aura/x/aura/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AuraKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
