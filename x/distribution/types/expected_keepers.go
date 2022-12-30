package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type StakingKeeper interface {
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation
	IterateValidators(ctx sdk.Context, f func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateBondedValidatorsByPower(ctx sdk.Context, f func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateLastValidators(ctx sdk.Context, f func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingtypes.ValidatorI
	ValidatorByConsAddr(ctx sdk.Context, address sdk.ConsAddress) stakingtypes.ValidatorI
	Slash(ctx sdk.Context, address sdk.ConsAddress, i int64, i2 int64, dec sdk.Dec)
	Jail(ctx sdk.Context, address sdk.ConsAddress)
	Unjail(ctx sdk.Context, address sdk.ConsAddress)
	MaxValidators(ctx sdk.Context) uint32
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	GetLastTotalPower(ctx sdk.Context) sdk.Int
	GetLastValidatorPower(ctx sdk.Context, valAddr sdk.ValAddress) int64
}
