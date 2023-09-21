package types_test

import (
	"testing"

	"github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/auth/vesting/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	orgtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMsg_Basic(t *testing.T) {
	app := tests.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	_ = ctx

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	sdk.GetConfig().SetBech32PrefixForAccount("aura", "aurapubkey")

	fromAddr, err := sdk.AccAddressFromBech32("aura1txe6y425gk7ef8xp6r7ze4da09nvwfr2fhafjl")
	require.NoError(t, err)

	toAddr, err := sdk.AccAddressFromBech32("aura1fqqrll4l62hlx36kw3mhav57n00lsy4kskvat8")
	require.NoError(t, err)

	startTime := int64(1676618561)

	periods := []orgtypes.Period{}

	msgAcc := types.NewMsgCreatePeriodicVestingAccount(fromAddr, toAddr, startTime, periods)

	require.NotNil(t, msgAcc)

	route := msgAcc.Route()
	require.Greater(t, len(route), 0)

	typeStr := msgAcc.Type()
	require.Equal(t, typeStr, "msg_create_periodic_vesting_account")

	signers := msgAcc.GetSigners()
	require.Equal(t, len(signers) > 0, true)

	signByte := msgAcc.GetSignBytes()
	require.NotNil(t, signByte)

	require.NoError(t, msgAcc.ValidateBasic())

}

func TestMsg_NotValid(t *testing.T) {
	app := tests.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	_ = ctx

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	sdk.GetConfig().SetBech32PrefixForAccount("aura", "aurapubkey")

	fromAddr, err := sdk.AccAddressFromBech32("aura1txe6y425gk7ef8xp6r7ze4da09nvwfr2fhafjl")
	require.NoError(t, err)

	emptyFromAddr, err := sdk.AccAddressFromBech32("")
	require.Error(t, err)

	toAddr, err := sdk.AccAddressFromBech32("aura1fqqrll4l62hlx36kw3mhav57n00lsy4kskvat8")
	require.NoError(t, err)

	emptyToAddr, err := sdk.AccAddressFromBech32("")
	require.Error(t, err)

	startTime := int64(1676618561)

	periods := []orgtypes.Period{}

	msgAccNotValidFrom := types.NewMsgCreatePeriodicVestingAccount(emptyFromAddr, toAddr, startTime, periods)

	require.Error(t, msgAccNotValidFrom.ValidateBasic())

	msgAccNotValidTo := types.NewMsgCreatePeriodicVestingAccount(fromAddr, emptyToAddr, startTime, periods)

	require.Error(t, msgAccNotValidTo.ValidateBasic())

	zeroTime := 0

	msgAccZeroTime := types.NewMsgCreatePeriodicVestingAccount(fromAddr, toAddr, int64(zeroTime), periods)

	require.Error(t, msgAccZeroTime.ValidateBasic())

	invalidPeriod := orgtypes.Period{
		Length: 0,
	}

	periods = []orgtypes.Period{invalidPeriod}
	msgAccInvalidPeriod := types.NewMsgCreatePeriodicVestingAccount(fromAddr, toAddr, startTime, periods)

	require.Error(t, msgAccInvalidPeriod.ValidateBasic())
}
