package keeper

import (
	"context"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type msgServer struct {
	Keeper        Keeper
	WasmKeeper    *wasmkeeper.PermissionedKeeper
	AccountKeeper types.AccountKeeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, wasmKeeper *wasmkeeper.PermissionedKeeper, accountKeeper types.AccountKeeper) types.MsgServer {
	return &msgServer{
		Keeper:        keeper,
		WasmKeeper:    wasmKeeper,
		AccountKeeper: accountKeeper,
	}
}

func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	saAddress, data, err := InstantiateSmartAccount(ctx, k.Keeper, k.WasmKeeper, msg)
	if err != nil {
		return nil, err
	}

	saAddressStr := saAddress.String()

	// get smart contract account by address
	scAccount := k.AccountKeeper.GetAccount(ctx, saAddress)
	if _, ok := scAccount.(*authtypes.BaseAccount); !ok {
		return &types.MsgCreateAccountResponse{},
			fmt.Errorf(types.ErrAccountNotFoundForAddress, saAddressStr)
	}

	smartAccount := types.NewSmartAccountFromAccount(scAccount)
	err = smartAccount.SetPubKey(&msg.PubKey)
	if err != nil {
		return nil, err
	}

	k.AccountKeeper.SetAccount(ctx, smartAccount)

	return &types.MsgCreateAccountResponse{
		Address: saAddressStr,
		Data:    data,
	}, nil
}
