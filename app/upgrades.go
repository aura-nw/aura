package app

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// next upgrade name
const upgradeName = "v0.2"

// RegisterUpgradeHandlers returns upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// consensus params
		// increase max gas as part of the upgrade to handle cosmwam
		consensusParams := app.BaseApp.GetConsensusParams(ctx)
		consensusParams.Block.MaxGas = 75_000_000 // 75M
		app.BaseApp.StoreConsensusParams(ctx, consensusParams)
		// wasm params
		wasmParams := app.WasmKeeper.GetParams(ctx)
		wasmParams.CodeUploadAccess = wasmtypes.AllowNobody
		wasmParams.MaxWasmCodeSize = DefaultMaxWasmCodeSize
		app.WasmKeeper.SetParams(ctx, wasmParams)

		govVotingParams := app.GovKeeper.GetVotingParams(ctx)
		govVotingParams.VotingPeriod = DefaultVotingPeriod
		app.GovKeeper.SetVotingParams(ctx, govVotingParams)

		govDepositParams := app.GovKeeper.GetDepositParams(ctx)
		govDepositParams.MaxDepositPeriod = DefaultDepositPeriod
		app.GovKeeper.SetDepositParams(ctx, govDepositParams)
		return vm, nil
	})
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := store.StoreUpgrades{}
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
