package keeper

import (
	"context"

	types "github.com/aura-nw/aura/x/smartaccount/types"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper Keeper
}

var _ typesv1.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) typesv1.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) ActivateAccount(goCtx context.Context, msg *typesv1.MsgActivateAccount) (*typesv1.MsgActivateAccountResponse, error) {
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

	return &typesv1.MsgActivateAccountResponse{
		Address: sAccount.String(),
	}, nil
}

func (k msgServer) Recover(goCtx context.Context, msg *typesv1.MsgRecover) (*typesv1.MsgRecoverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	smartAccount, err := k.Keeper.ValidateRecoverSA(ctx, msg)
	if err != nil {
		return nil, err
	}

	// public key
	pubKey, err := typesv1.PubKeyDecode(msg.PubKey)
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
			types.EventTypeSmartAccountRecovery,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Address),
		),
	)

	return &typesv1.MsgRecoverResponse{}, nil
}
