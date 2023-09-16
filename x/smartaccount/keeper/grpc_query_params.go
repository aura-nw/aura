package keeper

import (
	"context"

	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(c context.Context, req *typesv1.QueryParamsRequest) (*typesv1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &typesv1.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
