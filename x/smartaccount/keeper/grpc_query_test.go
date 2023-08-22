package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestQueryParams(t *testing.T) {
	ctx, app := helper.SetupGenesisTest()

	queryServer := app.SaKeeper

	for _, tc := range []struct {
		desc string
		msg  *types.QueryParamsRequest
		res  types.Params
		err  bool
	}{
		{
			desc: "query params successfully",
			msg:  &types.QueryParamsRequest{},
			res:  helper.GenesisState.Params,
			err:  false,
		},
		{
			desc: "error, query fail with nil message",
			msg:  nil,
			err:  true,
		},
	} {
		res, err := queryServer.Params(sdk.WrapSDKContext(ctx), tc.msg)

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.res, res.Params)
		}
	}
}

func TestQueryGenerateAccount(t *testing.T) {
	ctx, app := helper.SetupGenesisTest()

	creator := app.AccountKeeper.GetAllAccounts(ctx)[0]

	codeID, _, err := helper.StoreCodeID(app, ctx, creator.GetAddress(), helper.WasmPath2+"base.wasm")
	require.NoError(t, err)
	require.Equal(t, codeID, helper.DefaultCodeID)

	queryServer := app.SaKeeper

	pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *types.QueryGenerateAccountRequest
		err  bool
	}{
		{
			desc: "error, nil message when query",
			msg:  nil, // nil message
			err:  true,
		},
		{
			desc: "error, empty message when query",
			msg:  &types.QueryGenerateAccountRequest{}, // empty message
			err:  true,
		},
		{
			desc: "error, invalid public key value",
			msg: &types.QueryGenerateAccountRequest{
				// invlid pubkey
				PubKey: &codectypes.Any{
					TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
					Value:   []byte("error key value"),
				},
			},
			err: true,
		},
		{
			desc: "error, codeID not exist on chain",
			msg: &types.QueryGenerateAccountRequest{
				CodeID:  uint64(2), // code_id not exist
				PubKey:  pubKey,
				Salt:    helper.DefaultSalt,
				InitMsg: helper.DefaultMsg,
			},
			err: true,
		},
		{
			desc: "query generate account successfully",
			msg: &types.QueryGenerateAccountRequest{
				CodeID:  helper.DefaultCodeID,
				PubKey:  pubKey,
				Salt:    helper.DefaultSalt,
				InitMsg: helper.DefaultMsg,
			},
			err: false,
		},
	} {
		_, err := queryServer.GenerateAccount(sdk.WrapSDKContext(ctx), tc.msg)

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
