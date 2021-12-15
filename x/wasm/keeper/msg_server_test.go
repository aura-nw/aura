package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/aura-nw/aura/testutil/keeper"
	"github.com/aura-nw/aura/x/wasm/keeper"
	"github.com/aura-nw/aura/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.WasmKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
