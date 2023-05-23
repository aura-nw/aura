package keeper

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
)

type msgServer struct {
	Keeper        Keeper
	WasmKeeper    *wasmkeeper.PermissionedKeeper
	AccountKeeper types.AccountKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, wasmKeeper *wasmkeeper.PermissionedKeeper, accountKeeper types.AccountKeeper) types.MsgServer {
	return &msgServer{
		Keeper:        keeper,
		WasmKeeper:    wasmKeeper,
		AccountKeeper: accountKeeper,
	}
}

var _ types.MsgServer = msgServer{}
