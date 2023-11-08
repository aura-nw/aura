package keeper_test

import (
	"testing"
	"time"

	"github.com/aura-nw/aura/app"
	tests "github.com/aura-nw/aura/tests"
	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	App *app.App
	Ctx sdk.Context
	ctx sdk.Context
}

func (s *KeeperTestSuite) SetupSuite() {
	s.App = tests.Setup(s.T(), false)
	s.Ctx = s.App.NewContext(false, tmproto.Header{
		ChainID: "aura-testnet",
		Time:    time.Now(),
	})

	smartaccount.InitGenesis(s.Ctx, s.App.SaKeeper, *helper.GenesisState)

	/* ======== store wasm ======== */
	creator := s.App.AccountKeeper.GetAllAccounts(s.Ctx)[0]

	// store code
	codeID, _, err := helper.StoreCodeID(s.App, s.Ctx, creator.GetAddress(), helper.WasmPath2+"base.wasm")
	require.NoError(s.T(), err)
	require.Equal(s.T(), codeID, uint64(1))

	codeID, _, err = helper.StoreCodeID(s.App, s.Ctx, creator.GetAddress(), helper.WasmPath2+"recovery.wasm")
	require.NoError(s.T(), err)
	require.Equal(s.T(), codeID, uint64(2))
}

func (s *KeeperTestSuite) SetupTest() {
	s.ctx, _ = s.Ctx.CacheContext()
}

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
