package v700

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	samodulekeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	// SDK v47 modules
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// custom
	smartaccounttypesauranw "github.com/aura-nw/aura/x/smartaccount/types/auranw"
	smartaccounttypesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
)

// UpgradeName is the name of upgrade. This upgrade added new module
const UpgradeName = "v0.7.0"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	saKeeper samodulekeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
	consensusParamKeeper consensusparamkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		baseAppLegacySS := paramKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &consensusParamKeeper)

		// Migrate smartaccounts from `auranw` to `v1` verson
		// Change typeUrl from "auranw.aura.smartaccount.SmartAccount" to "aura.smartaccount.v1.SmartAccount"
		var iterErr error
		authKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
			if oldSa, ok := account.(*smartaccounttypesauranw.SmartAccount); ok {
				newSa := smartaccounttypesv1.NewSmartAccount(oldSa.Address, oldSa.AccountNumber, oldSa.Sequence)
				err := newSa.SetPubKey(oldSa.GetPubKey())
				if err != nil {
					iterErr = err
					return true
				}

				authKeeper.SetAccount(ctx, newSa)
			}
			return false
		})

		if iterErr != nil {
			return nil, iterErr
		}

		// update smartaccount params
		smartaccountParams := smartaccounttypesv1.DefaultParams()
		err := saKeeper.SetParams(ctx, smartaccountParams)
		if err != nil {
			return nil, err
		}

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		return versionMap, err
	}
}