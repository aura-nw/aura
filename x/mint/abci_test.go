package mint_test

import (
	"github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/mint"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestBeginBlocker_Basic(t *testing.T) {
	app := tests.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	dparams := app.MintKeeper.GetParams(ctx)
	params := types.Params{
		MintDenom:           "uaura",
		InflationMax:        sdk.NewDecWithPrec(12, 2),
		InflationRateChange: sdk.NewDecWithPrec(8, 2),
		InflationMin:        sdk.NewDecWithPrec(4, 2),
		BlocksPerYear:       5373084,
		GoalBonded:          dparams.GoalBonded,
	}
	app.MintKeeper.SetParams(ctx, params)

	mint.BeginBlocker(ctx, app.MintKeeper)
}
