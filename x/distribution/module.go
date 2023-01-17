package distribution

import (
	customkeeper "github.com/aura-nw/aura/x/distribution/keeper"
	"github.com/aura-nw/aura/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	distribution.AppModule

	keeper customkeeper.Keeper
}

func NewAppModule(cdc codec.Codec, customKeeper customkeeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, stackingKeeper types.StakingKeeper) AppModule {
	return AppModule{
		AppModule: distribution.NewAppModule(cdc, customKeeper.Keeper, accountKeeper, bankKeeper, stackingKeeper),
		keeper:    customKeeper,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	disttypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper.Keeper))
	disttypes.RegisterQueryServer(cfg.QueryServer(), am.keeper.Keeper)

	m := keeper.NewMigrator(am.keeper.Keeper)

	_ = cfg.RegisterMigration(disttypes.ModuleName, 1, m.Migrate1to2)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}
