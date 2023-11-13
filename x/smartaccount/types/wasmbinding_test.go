package types_test

import (
	"testing"

	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
)

func TestParseMessagesString(t *testing.T) {
	addr1 := "cosmos1p3ucd3ptpw902fluyjzhq3ffgq4ntddac9sa3s"
	acc1, err := sdk.AccAddressFromBech32(addr1)
	require.NoError(t, err)

	addr2 := "cosmos15hmqrc245kryaehxlch7scl9d9znxa58qkpjet"
	acc2, err := sdk.AccAddressFromBech32(addr2)
	require.NoError(t, err)

	bankMsg1 := banktypes.NewMsgSend(acc1, acc2, sdk.NewCoins())
	value1, err := proto.Marshal(bankMsg1)
	require.NoError(t, err)

	bankMsg2 := banktypes.NewMsgSend(acc2, acc1, sdk.NewCoins())
	value2, err := proto.Marshal(bankMsg2)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msgs []sdk.Msg
		data []types.Any
	}{
		{
			desc: "parse zero message successfully",
			msgs: []sdk.Msg{},
			data: []types.Any{},
		},
		{
			desc: "parse many messages successfully",
			msgs: []sdk.Msg{
				bankMsg1,
				bankMsg2,
			},
			data: []types.Any{
				{
					TypeURL: "/cosmos.bank.v1beta1.MsgSend",
					Value:   value1,
				},
				{
					TypeURL: "/cosmos.bank.v1beta1.MsgSend",
					Value:   value2,
				},
			},
		},
	} {
		data, err := types.ParseMessagesString(tc.msgs)
		require.NoError(t, err)

		require.Equal(t, tc.data, data)
	}
}
