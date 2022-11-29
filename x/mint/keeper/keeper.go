package keeper

import (
	"github.com/aura-nw/aura/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         sdk.StoreKey
	paramSpace       paramtypes.Subspace
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	auraKeeper       types.AuraKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper,
	auraKeeper types.AuraKeeper, feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		auraKeeper:       auraKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Return the wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}

func (k Keeper) GetMaxSupply(ctx sdk.Context) string {
	return k.auraKeeper.GetMaxSupply(ctx)
}

func (k Keeper) GetExcludeCirculatingAmount(ctx sdk.Context, denom string) sdk.Coin {
	return k.bankKeeper.GetExcludeCirculatingAmount(ctx, denom)
}

// CustomStakingTokenSupply implements an alias call to the underlying staking keeper's
// CustomStakingTokenSupply to be used in BeginBlocker.
func (k Keeper) CustomStakingTokenSupply(ctx sdk.Context, excludeAmount sdk.Int) sdk.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx).Sub(excludeAmount)
}

// CustomBondedRatio implements an alias call to the underlying staking keeper's
// CustomBondedRatio to be used in BeginBlocker.
func (k Keeper) CustomBondedRatio(ctx sdk.Context, excludeAmount sdk.Int) sdk.Dec {
	stakeSupply := k.CustomStakingTokenSupply(ctx, excludeAmount)
	if stakeSupply.IsPositive() {
		return k.stakingKeeper.TotalBondedTokens(ctx).ToDec().QuoInt(stakeSupply)
	}

	return sdk.ZeroDec()
}

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
