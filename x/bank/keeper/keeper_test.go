package keeper_test

import (
	"github.com/aura-nw/aura/tests"
	keeper "github.com/aura-nw/aura/x/bank/keeper"
	"github.com/aura-nw/aura/x/bank/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	bankKeeper keeper.BaseKeeper

	accountKeeper *testutil.MockAccountKeeper
	auraKeeper    *testutil.MockAuraKeeper
}

func (s *KeeperTestSuite) SetupTest() {
	encConfig := tests.MakeTestEncodingConfig(bank.AppModuleBasic{})
	key := sdk.NewKVStoreKey(banktypes.StoreKey)
	testCtx := tests.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_key"))

	s.ctx = testCtx.Ctx

	ctrl := gomock.NewController(s.T())

	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	auraKeeper := testutil.NewMockAuraKeeper(ctrl)

	pk := tests.GetParamsKeeper()
	subspace := pk.Subspace(banktypes.ModuleName)

	s.bankKeeper = keeper.NewBaseKeeper(encConfig.Codec, key, accountKeeper, subspace, map[string]bool{}, auraKeeper)

	s.accountKeeper = accountKeeper
	s.auraKeeper = auraKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestGetExcludeCirculatingAmount() {
	s.auraKeeper.EXPECT().GetExcludeCirculatingAddr(s.ctx).Return([]sdk.AccAddress{})
	s.Require().NotNil(s.bankKeeper.GetExcludeCirculatingAmount(s.ctx, "uaura"))
}
