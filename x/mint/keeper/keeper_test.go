package keeper_test

import (
	"testing"

	"github.com/aura-nw/aura/utils/testutil"
	"github.com/aura-nw/aura/x/mint"
	"github.com/aura-nw/aura/x/mint/keeper"
	minttestutil "github.com/aura-nw/aura/x/mint/testutil"
	"github.com/aura-nw/aura/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	mintKeeper    keeper.Keeper
	stakingKeeper *minttestutil.MockStakingKeeper
	bankKeeper    *minttestutil.MockBankKeeper
	accountKeeper *minttestutil.MockAccountKeeper
	auraKeeper    *minttestutil.MockAuraKeeper
	pk            paramskeeper.Keeper
}

func (s *KeeperTestSuite) SetupTest() {
	encCfg := testutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))

	s.ctx = testCtx.Ctx

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	auraKeeper := minttestutil.NewMockAuraKeeper(ctrl)
	pk := testutil.GetParamsKeeper()

	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{})

	feeCollector := authTypes.FeeCollectorName
	subspace := pk.Subspace(types.ModuleName)

	s.mintKeeper = keeper.NewKeeper(encCfg.Codec, key, subspace, stakingKeeper, accountKeeper, bankKeeper, auraKeeper, feeCollector)

	s.stakingKeeper = stakingKeeper
	s.bankKeeper = bankKeeper
	s.accountKeeper = accountKeeper
	s.auraKeeper = auraKeeper
	s.pk = pk
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestCustomStakingTokenSupply() {
	stakingTokenSupply := sdk.NewIntFromUint64(1_000_000)
	s.stakingKeeper.EXPECT().StakingTokenSupply(s.ctx).Return(stakingTokenSupply)

	excludeAmount := sdk.NewIntFromUint64(100_000)
	s.Require().Equal(s.mintKeeper.CustomStakingTokenSupply(s.ctx, excludeAmount), stakingTokenSupply.Sub(excludeAmount))
}

func (s *KeeperTestSuite) TestMintCoins() {
	coins := sdk.NewCoins(sdk.NewCoin("uaura", sdk.NewInt(1000000)))
	s.bankKeeper.EXPECT().MintCoins(s.ctx, types.ModuleName, coins).Return(nil)
	s.Require().Equal(s.mintKeeper.MintCoins(s.ctx, sdk.NewCoins()), nil)
	s.Require().Nil(s.mintKeeper.MintCoins(s.ctx, coins))
}

func (s *KeeperTestSuite) TestCustomBondedRatio() {
	stakingSupply := sdk.NewIntFromUint64(1_000_000)
	s.stakingKeeper.EXPECT().StakingTokenSupply(s.ctx).Return(stakingSupply).AnyTimes()

	excludeAmount := sdk.NewIntFromUint64(100_000)

	customStaking := s.mintKeeper.CustomStakingTokenSupply(s.ctx, excludeAmount)

	bonded := sdk.NewIntFromUint64(500_000)
	s.stakingKeeper.EXPECT().TotalBondedTokens(s.ctx).Return(bonded)

	s.Require().Equal(s.mintKeeper.CustomBondedRatio(s.ctx, excludeAmount), bonded.ToDec().QuoInt(customStaking))

}
