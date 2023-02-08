package mint_test

import (
	"github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/mint"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	basemint "github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := tests.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}

func TestModule_BeginBlocker(t *testing.T) {
	app := tests.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	encCfg := tests.MakeTestEncodingConfig(basemint.AppModuleBasic{})

	m := mint.NewAppModule(encCfg.Codec, app.MintKeeper, app.AccountKeeper)
	req := abci.RequestBeginBlock{}
	m.BeginBlock(ctx, req)
}
