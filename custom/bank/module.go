package bank

import (
	"github.com/aura-nw/aura/custom/bank/client/cli"
	custombank "github.com/aura-nw/aura/custom/bank/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

// AppModule implements an application module for the custom bank module.
type AppModule struct {
	bank.AppModule

	keeper custombank.BaseKeeper
}

type AppModuleBasic struct {
	bank.AppModuleBasic
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombank.BaseKeeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModule: bank.NewAppModule(cdc, keeper.BaseKeeper, accountKeeper),
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := keeper.NewMigrator(am.keeper.BaseKeeper)
	cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
}

// GetTxCmd returns the root tx command for the bank module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}
