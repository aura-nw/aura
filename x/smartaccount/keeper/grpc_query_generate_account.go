package keeper

import (
	"context"

	"github.com/aura-nw/aura/x/smartaccount/types"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GenerateAccount(goCtx context.Context, req *typesv1.QueryGenerateAccountRequest) (*typesv1.QueryGenerateAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pubKey, err := typesv1.PubKeyDecode(req.PubKey)
	if err != nil {
		return nil, err
	}

	contractAddress, err := types.Instantiate2Address(ctx, k.WasmKeeper, req.CodeID, req.InitMsg, req.Salt, pubKey)
	if err != nil {
		return nil, err
	}

	return &typesv1.QueryGenerateAccountResponse{
		Address: contractAddress.String(),
	}, nil
}
