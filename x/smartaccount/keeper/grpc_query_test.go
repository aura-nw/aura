package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestQueryParams(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	queryServer := app.SaKeeper

	for _, tc := range []struct {
		desc string
		msg  *typesv1.QueryParamsRequest
		res  typesv1.Params
		err  bool
	}{
		{
			desc: "query params successfully",
			msg:  &typesv1.QueryParamsRequest{},
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
	ctx, app := helper.SetupGenesisTest(t)

	creator := app.AccountKeeper.GetAllAccounts(ctx)[0]

	codeID, _, err := helper.StoreCodeID(app, ctx, creator.GetAddress(), helper.WasmPath2+"base.wasm")
	require.NoError(t, err)
	require.Equal(t, codeID, helper.DefaultCodeID)

	queryServer := app.SaKeeper

	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *typesv1.QueryGenerateAccountRequest
		err  bool
	}{
		{
			desc: "error, nil message when query",
			msg:  nil, // nil message
			err:  true,
		},
		{
			desc: "error, empty message when query",
			msg:  &typesv1.QueryGenerateAccountRequest{}, // empty message
			err:  true,
		},
		{
			desc: "error, invalid public key value",
			msg: &typesv1.QueryGenerateAccountRequest{
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
			msg: &typesv1.QueryGenerateAccountRequest{
				CodeID:  uint64(2), // code_id not exist
				PubKey:  pubKey,
				Salt:    helper.DefaultSalt,
				InitMsg: helper.DefaultMsg,
			},
			err: true,
		},
		{
			desc: "query generate account successfully",
			msg: &typesv1.QueryGenerateAccountRequest{
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
