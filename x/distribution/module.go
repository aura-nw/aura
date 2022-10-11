package distribution

import (
	customdistrkeeper "github.com/aura-nw/aura/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type AppModule struct {
	distr.AppModule

	keeper customdistrkeeper.Keeper
}

func NewAppModule(
	cdc codec.Codec, keeper customdistrkeeper.Keeper, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper,
) AppModule {
	return AppModule{
		AppModule: distr.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper, stakingKeeper),
		keeper:    keeper,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), customdistrkeeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := keeper.NewMigrator(am.keeper.Keeper)
	cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
}
