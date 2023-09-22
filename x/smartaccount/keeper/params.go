package keeper

import (
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params typesv1.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params typesv1.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.paramstore.SetParamSet(ctx, &params)
	return nil
}
