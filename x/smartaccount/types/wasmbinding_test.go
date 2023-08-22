package types_test

import (
	"testing"

	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

const (
	AccountAddressPrefix = "aura"
)

func TestParseMessagesString(t *testing.T) {

	sdk.GetConfig().SetBech32PrefixForAccount(AccountAddressPrefix, "")

	addr1 := "aura17n858c4urvenkf7edjs4uumudej3ekyv432e34"
	acc1, err := sdk.AccAddressFromBech32(addr1)
	require.NoError(t, err)

	addr2 := "aura1pxf370s3smrfcpkx5x06a6jl0k0cfy9evq7ddr"
	acc2, err := sdk.AccAddressFromBech32(addr2)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msgs []sdk.Msg
		data []types.MsgData
	}{
		{
			desc: "parse zero message successfully",
			msgs: []sdk.Msg{},
			data: []types.MsgData{},
		},
		{
			desc: "parse many messages successfully",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1, acc2, sdk.NewCoins()),
				banktypes.NewMsgSend(acc2, acc1, sdk.NewCoins()),
			},
			data: []types.MsgData{
				{
					TypeURL: "/cosmos.bank.v1beta1.MsgSend",
					Value:   "{\"from_address\":\"" + addr1 + "\",\"to_address\":\"" + addr2 + "\",\"amount\":[]}",
				},
				{
					TypeURL: "/cosmos.bank.v1beta1.MsgSend",
					Value:   "{\"from_address\":\"" + addr2 + "\",\"to_address\":\"" + addr1 + "\",\"amount\":[]}",
				},
			},
		},
	} {
		data, err := types.ParseMessagesString(tc.msgs)
		require.NoError(t, err)

		require.Equal(t, tc.data, data)
	}
}
