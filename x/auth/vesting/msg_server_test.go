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

	msgServer := vesting.NewMsgServerImpl(app.AccountKeeper, app.BankKeeper)

	message := &types.MsgCreatePeriodicVestingAccount{}

	_, err := msgServer.CreatePeriodicVestingAccount(sdk.WrapSDKContext(ctx), message)

	require.NotNil(t, err)

}
