package keeper

import (
	"encoding/json"
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

// StorePacketCallback stores which contract will be listening for the ack or timeout of a packet
func (k Keeper) StoreSmartAccount(ctx sdk.Context, accountAddress string) error {
	store := ctx.KVStore(k.storeKey)
	value := types.SmartAccountValue{
		Type:   types.SmartAccountI,
		Active: true,
	}

	// json marshal should always work
	bytes, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf(types.ErrStoreSmartAccount, err.Error())
	}

	store.Set([]byte(accountAddress), bytes)

	return nil
}

// GetPacketCallback returns the bech32 addr of the contract that is expecting a callback from a packet
func (k Keeper) GetSmartAccount(ctx sdk.Context, accountAddress string) types.SmartAccountValue {
	store := ctx.KVStore(k.storeKey)
	accountValue := store.Get([]byte(accountAddress))

	var value types.SmartAccountValue

	if accountValue == nil {
		return value
	}

	_ = json.Unmarshal(accountValue, &value)

	return value
}

func (k Keeper) SetSmartAccountStatus(ctx sdk.Context, accountAddress string, status bool) error {
	store := ctx.KVStore(k.storeKey)

	accountValue := store.Get([]byte(accountAddress))

	var value types.SmartAccountValue

	if accountValue == nil {
		return fmt.Errorf(types.ErrSetSmartAccountStatus, "account address not found")
	}

	err := json.Unmarshal(accountValue, &value)
	if err != nil {
		return fmt.Errorf(types.ErrSetSmartAccountStatus, err.Error())
	}

	return nil
}
