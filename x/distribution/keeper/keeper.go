package keeper

import (
	"github.com/aura-nw/aura/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
	distkeeper.Keeper
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName string, blockedAddrs map[string]bool) Keeper {
	baseKeeper := distkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName, blockedAddrs)
	return Keeper{
		Keeper:        baseKeeper,
		authKeeper:    ak,
		bankKeeper:    bk,
		stakingKeeper: sk,
		storeKey:      key,
		cdc:           cdc,
	}
}
