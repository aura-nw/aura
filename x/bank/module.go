package bank

import (
	custombank "github.com/aura-nw/aura/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AppModule implements an application module for the custom bank module.
type AppModule struct {
	bank.AppModule

	keeper         custombank.BaseKeeper
	legacySubspace bankexported.Subspace // need review
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombank.BaseKeeper, accountKeeper types.AccountKeeper, ss bankexported.Subspace) AppModule {
	return AppModule{
		AppModule:      bank.NewAppModule(cdc, keeper.BaseKeeper, accountKeeper, ss),
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := keeper.NewMigrator(am.keeper.BaseKeeper, am.legacySubspace) // need review
	_ = cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
}
