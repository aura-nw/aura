package tests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func GetParamsKeeper() paramskeeper.Keeper {
	encodingCfg := MakeTestEncodingConfig(params.AppModuleBasic{})
	key := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkey := sdk.NewTransientStoreKey("params_transient_test")

	pk := paramskeeper.NewKeeper(encodingCfg.Codec, encodingCfg.Amino, key, tkey)

	return pk
}
