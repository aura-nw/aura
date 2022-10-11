package keeper

import (
	customdistrtypes "github.com/aura-nw/aura/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/pkg/errors"
)

type Keeper struct {
	distributionkeeper.Keeper

	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName string, blockedAddrs map[string]bool,
) Keeper {
	return Keeper{
		Keeper:   distributionkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName, blockedAddrs),
		storeKey: key,
		cdc:      cdc,
	}
}

func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	claimInfo, err := k.IsClaimReward(ctx, delAddr)
	if err != nil {
		return nil, err
	}

	coins, err := k.Keeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	k.SetClaimInfo(ctx, delAddr, claimInfo)
	return coins, nil
}

func (k Keeper) IsClaimReward(ctx sdk.Context, delAddr sdk.AccAddress) (customdistrtypes.ClaimInfo, error) {
	currentHeight := ctx.BlockHeight()

	claimInfo, err := k.GetClaimInfo(ctx, delAddr)
	if err != nil {
		return claimInfo, err
	}
	k.Logger(ctx).Error("IsClaimReward", "number", claimInfo.ClaimBlockNum)

	if claimInfo.IsEmpty() {
		claimInfo.Address = delAddr.String()
		claimInfo.ClaimBlockNum = currentHeight
		claimInfo.ClaimTime = ctx.BlockTime()
		return claimInfo, nil
	}

	// TODO: hardcode number 5000
	if currentHeight-claimInfo.ClaimBlockNum < 500 {
		k.Logger(ctx).Error("IsClaimReward", "currentHeight", currentHeight)
		return customdistrtypes.ClaimInfo{}, errors.New("unable claim reward in the period")
	}

	claimInfo.ClaimBlockNum = currentHeight
	claimInfo.ClaimTime = ctx.BlockTime()
	return claimInfo, nil
}
