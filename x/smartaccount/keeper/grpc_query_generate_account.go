package keeper

import (
	"context"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
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

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, fmt.Errorf(types.ErrAddressFromBech32, err)
	}

	pub_key, err := PubKeyDecode(req.PubKey)
	if err != nil {
		return nil, err
	}

	codeInfo := k.wasmKeeper.GetCodeInfo(ctx, req.CodeID)
	if codeInfo == nil {
		return nil, fmt.Errorf(types.ErrNoSuchCodeID, req.CodeID)
	}

	salt, err := types.GenerateSalt(req.Owner, req.CodeID, req.InitMsg, pub_key.Key)
	if err != nil {
		return nil, err
	}

	addrGenerator := wasmkeeper.PredicableAddressGenerator(owner, salt, req.InitMsg, true)
	contractAddress := addrGenerator(ctx, req.CodeID, codeInfo.CodeHash)
	if k.wasmKeeper.HasContractInfo(ctx, contractAddress) {
		return nil, fmt.Errorf(types.ErrInstantiateDuplicate)
	}

	return &types.QueryGenerateAccountResponse{
		Address: contractAddress.String(),
	}, nil
}
