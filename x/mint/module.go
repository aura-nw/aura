package mint

import (
	custommint "github.com/aura-nw/aura/x/mint/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/mint/exported"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// AppModule implements an application module for the mint module.
type AppModule struct {
	mint.AppModule

	keeper     custommint.Keeper
	authKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper custommint.Keeper,
	ak types.AccountKeeper,
	ic types.InflationCalculationFn,
	ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModule:  mint.NewAppModule(cdc, keeper.Keeper, ak, ic, ss),
		keeper:     keeper,
		authKeeper: ak,
	}
}

// BeginBlock returns the begin blocker for the mint module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}
