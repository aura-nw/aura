package vesting_test

import (
	"testing"

	"github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/auth/vesting"
	"github.com/aura-nw/aura/x/auth/vesting/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestCreatePeriodicVestingAccount(t *testing.T) {
	app := tests.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	sdk.GetConfig().SetBech32PrefixForAccount("aura", "aurapubkey")

	msgServer := vesting.NewMsgServerImpl(app.AccountKeeper, app.BankKeeper)

	messageEmpty := &types.MsgCreatePeriodicVestingAccount{}

	_, err := msgServer.CreatePeriodicVestingAccount(sdk.WrapSDKContext(ctx), messageEmpty)

	require.NotNil(t, err)

	messageBasic := &types.MsgCreatePeriodicVestingAccount{
		FromAddress: "aura1txe6y425gk7ef8xp6r7ze4da09nvwfr2fhafjl",
		ToAddress:   "aura1fqqrll4l62hlx36kw3mhav57n00lsy4kskvat8",
		StartTime:   int64(18388373),
	}

	_, err = msgServer.CreatePeriodicVestingAccount(sdk.WrapSDKContext(ctx), messageBasic)

	t.Logf("%v", err)

	require.NoError(t, err)
}
