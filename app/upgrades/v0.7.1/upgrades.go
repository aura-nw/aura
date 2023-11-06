package v701

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	samodulekeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName is the name of upgrade. This upgrade added new module
const UpgradeName = "v0.7.1"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	saKeeper samodulekeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
	consensusParamKeeper consensusparamkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {

	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		// hard code upgrades for mainnet
		if ctx.ChainID() == "xstaxy-1" {
			return UpgradeMainnetHandler(ctx, plan, vm, mm, configurator, saKeeper, paramKeeper, consensusParamKeeper, ibcKeeper, authKeeper)
		} else {
			logger := ctx.Logger().With("upgrade", UpgradeName)

			// Run migrations
			logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
			versionMap, err := mm.RunMigrations(ctx, configurator, vm)
			if err != nil {
				return nil, err
			}
			logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

			return versionMap, err
		}
	}
}
