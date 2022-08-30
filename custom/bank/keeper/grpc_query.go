package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TotalSupply implements custom the Query/TotalSupply gRPC method
func (k BaseKeeper) TotalSupply(ctx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalSupply, pageRes, err := k.GetPaginatedTotalSupply(sdkCtx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var excludeTotalSupply sdk.Coins
	for _, coin := range totalSupply {
		excludeCoin := k.GetExcludeCirculatingAmount(sdkCtx, coin.GetDenom())
		excludeTotalSupply = excludeTotalSupply.Add(excludeCoin)
	}

	totalSupply = totalSupply.Sub(excludeTotalSupply)

	return &types.QueryTotalSupplyResponse{Supply: totalSupply, Pagination: pageRes}, nil
}

// SupplyOf implements custom the Query/SupplyOf gRPC method
func (k BaseKeeper) SupplyOf(c context.Context, req *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	supply := k.GetSupply(ctx, req.Denom)

	excludeCoin := k.GetExcludeCirculatingAmount(ctx, req.GetDenom())
	if !excludeCoin.IsZero() {
		supply = supply.Sub(excludeCoin)
	}

	return &types.QuerySupplyOfResponse{Amount: sdk.NewCoin(req.Denom, supply.Amount)}, nil
}
