package keeper

import (
	"context"
	"fmt"

	db "github.com/aura-nw/aura/database"
	"github.com/aura-nw/aura/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
)

type msgServer struct {
	feegrant.MsgServer
	Indexer *db.Db
}

// NewMsgServerImpl returns an implementation of the feegrant MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k keeper.Keeper, indexer *db.Db) feegrant.MsgServer {
	return &msgServer{
		MsgServer: keeper.NewMsgServerImpl(k),
		Indexer:   indexer,
	}
}

var _ feegrant.MsgServer = msgServer{}

// GrantAllowance grants an allowance from the granter's funds to be used by the grantee.
func (k msgServer) GrantAllowance(goCtx context.Context, msg *feegrant.MsgGrantAllowance) (*feegrant.MsgGrantAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Println(k.Indexer)
	k.MsgServer.GrantAllowance(goCtx, msg)
	err := k.Indexer.SaveFeeGrantAllowance(types.FeeGrant{Grant: feegrant.Grant{Granter: msg.Granter, Grantee: msg.Grantee, Allowance: msg.Allowance}, Height: ctx.BlockHeight()})
	if err != nil {
		fmt.Println(err)
	}
	return &feegrant.MsgGrantAllowanceResponse{}, nil
}

// RevokeAllowance revokes a fee allowance between a granter and grantee.
func (k msgServer) RevokeAllowance(goCtx context.Context, msg *feegrant.MsgRevokeAllowance) (*feegrant.MsgRevokeAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.MsgServer.RevokeAllowance(goCtx, msg)
	k.Indexer.DeleteFeeGrantAllowance(types.GrantRemoval{Granter: msg.Granter, Grantee: msg.Grantee, Height: ctx.BlockHeight()})
	return &feegrant.MsgRevokeAllowanceResponse{}, nil
}
