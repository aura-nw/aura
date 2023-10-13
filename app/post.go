package app

import (
	smartaccount "github.com/aura-nw/aura/x/smartaccount"
	smartaccountkeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
)

type PostHandlerOptions struct {
	posthandler.HandlerOptions

	SmartAccountKeeper smartaccountkeeper.Keeper
}

func NewPostHandler(options PostHandlerOptions) (sdk.PostHandler, error) {
	postDecorators := []sdk.PostDecorator{
		smartaccount.NewAfterTxDecorator(options.SmartAccountKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...), nil
}
