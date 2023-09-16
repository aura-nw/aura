package v601

import (
	samodulekeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	samoduletypes "github.com/aura-nw/aura/x/smartaccount/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName is the name of upgrade. This upgrade added new module
const UpgradeName = "v0.6.1"

func CreateUpgradeHandler(
	mm *module.Manager,
	saKeeper samodulekeeper.Keeper,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		// update smartaccount params
		smartaccountParams := samoduletypes.DefaultParams()
		err := saKeeper.SetParams(ctx, smartaccountParams)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
