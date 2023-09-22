package smartaccount

import (
	"github.com/aura-nw/aura/x/smartaccount/keeper"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState typesv1.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	k.SetNextAccountID(ctx, genState.GetSmartAccountId())
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *typesv1.GenesisState {
	genesis := typesv1.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.SmartAccountId = k.GetNextAccountID(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
