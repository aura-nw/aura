package v0_3_4

import (
	"github.com/aura-nw/aura/x/aura"
	auramodulekeeper "github.com/aura-nw/aura/x/aura/keeper"
	"github.com/aura-nw/aura/x/aura/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const UpgradeName = "v0.3.4"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	auraKeeper auramodulekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		aura.InitGenesis(ctx, auraKeeper, types.GenesisState{Params: types.DefaultParams()})

		ctx.Logger().Error("CreateUpgradeHandler", "max", auraKeeper.GetMaxSupply(ctx))

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
