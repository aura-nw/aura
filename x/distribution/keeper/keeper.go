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

	storeKey   storetypes.StoreKey
	auraKeeper customdistrtypes.AuraKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName string, blockedAddrs map[string]bool, auraKeeper customdistrtypes.AuraKeeper,
) Keeper {
	return Keeper{
		Keeper:     distributionkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName, blockedAddrs),
		storeKey:   key,
		auraKeeper: auraKeeper,
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
	blockTime := ctx.BlockTime()

	claimInfo, err := k.GetClaimInfo(ctx, delAddr)
	if err != nil {
		return claimInfo, err
	}

	if claimInfo.IsEmpty() {
		claimInfo.ClaimBlockNum = ctx.BlockHeight()
		claimInfo.ClaimTime = blockTime
		return claimInfo, nil
	}

	if blockTime.Sub(claimInfo.ClaimTime).Milliseconds() < k.auraKeeper.GetClaimDuration(ctx) {
		return customdistrtypes.ClaimInfo{}, errors.New("unable claim reward in the period")
	}

	claimInfo.ClaimBlockNum = ctx.BlockHeight()
	claimInfo.ClaimTime = ctx.BlockTime()
	return claimInfo, nil
}
