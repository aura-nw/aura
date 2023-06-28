package keeper

import (
	"context"

	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GenerateAccount(goCtx context.Context, req *types.QueryGenerateAccountRequest) (*types.QueryGenerateAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pubKey, err := types.PubKeyDecode(req.PubKey)
	if err != nil {
		return nil, err
	}

	contractAddress, err := types.Instantiate2Address(ctx, k.wasmKeeper, req.Owner, req.CodeID, req.InitMsg, pubKey.Bytes())
	if err != nil {
		return nil, err
	}

	return &types.QueryGenerateAccountResponse{
		Address: contractAddress.String(),
	}, nil
}
