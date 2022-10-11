package keeper

import (
	"encoding/json"
	"github.com/aura-nw/aura/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetClaimInfo(ctx sdk.Context, delAddr sdk.AccAddress) (types.ClaimInfo, error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorClaimAddrKey(delAddr))
	if b == nil {
		return types.ClaimInfo{}, nil
	}

	var claimInfo types.ClaimInfo
	if err := k.UnmarshalClaimInfo(b, &claimInfo); err != nil {
		return types.ClaimInfo{}, err
	}

	return claimInfo, nil
}

func (k Keeper) SetClaimInfo(ctx sdk.Context, delAddr sdk.AccAddress, info types.ClaimInfo) {
	store := ctx.KVStore(k.storeKey)

	value, err := k.MarshalClaimInfo(info)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetDelegatorClaimAddrKey(delAddr), value)
}

func (k Keeper) MarshalClaimInfo(claimInfo types.ClaimInfo) ([]byte, error) {
	bz, err := json.Marshal(&claimInfo)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (k Keeper) UnmarshalClaimInfo(bz []byte, claimInfo *types.ClaimInfo) error {
	if err := json.Unmarshal(bz, claimInfo); err != nil {
		return err
	}

	return nil
}
