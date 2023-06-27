package keeper

import (
	"context"

	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) ActivateAccount(goCtx context.Context, msg *types.MsgActivateAccount) (*types.MsgActivateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate smart account
	sAccount, err := k.Keeper.ValidateActiveSA(ctx, msg)
	if err != nil {
		return nil, err
	}

	// get current sequence of account
	currentSeq := sAccount.GetSequence()

	// set sequence to 0 so we can instantiate it later
	err = k.Keeper.PrepareBeforeActive(ctx, sAccount)
	if err != nil {
		return nil, err
	}

	pubKey, err := k.Keeper.ActiveSmartAccount(ctx, msg, sAccount)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.HandleAfterActive(ctx, sAccount, currentSeq, pubKey)
	if err != nil {
		return nil, err
	}

	return &types.MsgActivateAccountResponse{
		Address: sAccount.String(),
	}, nil
}

func (k msgServer) Recover(goCtx context.Context, msg *types.MsgRecover) (*types.MsgRecoverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	smartAccount, err := k.Keeper.ValidateRecoverSA(ctx, msg)
	if err != nil {
		return nil, err
	}

	// public key
	pubKey, err := types.PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, err
	}

	// verify logic recover from smart contract
	err = k.Keeper.CallSMValidate(ctx, msg, smartAccount.GetAddress(), pubKey)
	if err != nil {
		return nil, err
	}

	// recover public key for smart account
	err = k.Keeper.UpdateAccountPubKey(ctx, smartAccount, pubKey)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Info(
		"smart account recovery",
		types.AttributeKeyCreator, msg.Creator,
		types.AttributeKeyContractAddr, msg.Address,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAccountRecovery,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Address),
		),
	)

	return &types.MsgRecoverResponse{}, nil
}
