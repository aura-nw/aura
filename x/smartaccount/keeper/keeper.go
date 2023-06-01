package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ------------------------------- NextAccountId -------------------------------

func (k Keeper) GetAndIncrementNextAccountID(ctx sdk.Context) uint64 {
	id := k.GetNextAccountID(ctx)

	k.SetNextAccountID(ctx, id+1)

	return id
}

func (k Keeper) GetNextAccountID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	return sdk.BigEndianToUint64(store.Get(types.KeyPrefix(types.AccountIDKey)))
}

func (k Keeper) SetNextAccountID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.AccountIDKey), sdk.Uint64ToBigEndian(id))
}
