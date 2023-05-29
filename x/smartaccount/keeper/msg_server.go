package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	if k.WasmKeeper == nil {
		return nil, fmt.Errorf(types.ErrWasmKeeper)
	}

	saAddress, err := instantiateSmartAccount(ctx, k.WasmKeeper, msg)
	if err != nil {
		return nil, err
	}

	saAddressStr := saAddress.String()

	// get smart account by address
	smartAccount := k.AccountKeeper.GetAccount(ctx, saAddress)
	if smartAccount == nil {
		return &types.MsgCreateAccountResponse{},
			fmt.Errorf(types.ErrAccountNotFoundForAddress, saAddressStr)
	}

	bz, err := hex.DecodeString(msg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf(types.ErrBadPublicKey, err.Error())
	}

	// new secp25k61 public key
	newPubkey := &secp256k1.PubKey{Key: nil}
	keyErr := newPubkey.UnmarshalAmino(bz)
	if keyErr != nil {
		return nil, fmt.Errorf(types.ErrBadPublicKey, keyErr.Error())
	}

	// set new public key
	err = smartAccount.SetPubKey(newPubkey)
	if err != nil {
		return nil, fmt.Errorf(types.ErrSetPublickey, err.Error())
	}

	// save account to db after updated public key
	k.AccountKeeper.SetAccount(ctx, smartAccount)

	// save smart account address to db, so we can check later
	k.Keeper.StoreSmartAccount(ctx, saAddressStr)

	return &types.MsgCreateAccountResponse{}, nil
}
