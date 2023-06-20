package keeper

import (
	"context"

	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ActivateAccount(goCtx context.Context, msg *types.MsgActivateAccount) (*types.MsgActivateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgActivateAccountResponse{}, nil
}
