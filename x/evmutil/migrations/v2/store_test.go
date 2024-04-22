package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v2evmutil "github.com/aura-nw/aura/x/evmutil/migrations/v2"
	"github.com/aura-nw/aura/x/evmutil/types"
)

func TestStoreMigrationAddsKeyTableIncludingNewParam(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	evmutilKey := sdk.NewKVStoreKey(types.ModuleName)
	tEvmutilKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(evmutilKey, tEvmutilKey)
	paramstore := paramtypes.NewSubspace(encCfg.Codec, encCfg.Amino, evmutilKey, tEvmutilKey, types.ModuleName)

	// Check param doesn't exist before
	require.False(t, paramstore.Has(ctx, types.KeyAllowedCosmosDenoms))

	// Run migrations.
	err := v2evmutil.MigrateStore(ctx, paramstore)
	require.NoError(t, err)

	// Make sure the new params are set.
	require.True(t, paramstore.Has(ctx, types.KeyAllowedCosmosDenoms))
}

func TestStoreMigrationSetsNewParamOnExistingKeyTable(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	evmutilKey := sdk.NewKVStoreKey(types.ModuleName)
	tEvmutilKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(evmutilKey, tEvmutilKey)
	paramstore := paramtypes.NewSubspace(encCfg.Codec, encCfg.Amino, evmutilKey, tEvmutilKey, types.ModuleName)
	paramstore.WithKeyTable(types.ParamKeyTable())

	// expect it to have key table
	require.True(t, paramstore.HasKeyTable())
	// expect it to not have new param
	require.False(t, paramstore.Has(ctx, types.KeyAllowedCosmosDenoms))

	// Run migrations.
	err := v2evmutil.MigrateStore(ctx, paramstore)
	require.NoError(t, err)

	// Make sure the new params are set.
	require.True(t, paramstore.Has(ctx, types.KeyAllowedCosmosDenoms))
}
