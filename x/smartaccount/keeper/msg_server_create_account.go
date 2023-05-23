package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	SmartAccountLabel = "smart account"
)

func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.WasmKeeper == nil {
		return nil, fmt.Errorf(types.ErrWasmKeeper)
	}

	creator, aMsg := sdk.AccAddressFromBech32(msg.Creator)
	if aMsg != nil {
		return nil, fmt.Errorf(types.ErrAddressFromBech32, aMsg)
	}

	// instantiate smartcontract by code id
	address, _, iErr := k.WasmKeeper.Instantiate(
		ctx,
		uint64(msg.CodeId),
		creator,                 // owner
		creator,                 // admin
		[]byte(msg.InitMessage), // message
		SmartAccountLabel,       // label
		sdk.NewCoins(sdk.NewCoin("uaura", sdk.NewInt(0))), // test funds
	)
	if iErr != nil {
		return nil, fmt.Errorf(types.ErrBadInstantiateMsg, iErr.Error())
	}

	// get smart account by address
	smartAccount := k.AccountKeeper.GetAccount(ctx, address)
	if smartAccount == nil {
		return &types.MsgCreateAccountResponse{
			Address: address.String(),
		}, fmt.Errorf(types.ErrAccountNotFoundForAddress, address.String())
	}

	bz, err := hex.DecodeString(msg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf(types.ErrBadPublicKey, err.Error())
	}

	newPubkey := &secp256k1.PubKey{Key: bz} // new secp25k61 public key

	// set new public key
	err = smartAccount.SetPubKey(newPubkey)
	if err != nil {
		return nil, fmt.Errorf(types.ErrSetPublickey, err.Error())
	}

	// save account to db after updated public key
	k.AccountKeeper.SetAccount(ctx, smartAccount)

	// save smart account address to db, so we can check later
	k.Keeper.StoreSmartAccount(ctx, address.String())

	return &types.MsgCreateAccountResponse{
		Address: address.String(),
	}, nil
}
