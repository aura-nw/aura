package v1_test

import (
	"testing"

	"github.com/aura-nw/aura/app"
	helper "github.com/aura-nw/aura/tests/smartaccount"
	v1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestActivateAccountValidateBasic(t *testing.T) {
	pubKey, err := v1.PubKeyToAny(app.MakeEncodingConfig().Marshaler, helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *v1.MsgActivateAccount
		err  bool
	}{
		{
			desc: "error, account address invalid bench32 string",
			msg: &v1.MsgActivateAccount{
				AccountAddress: "abcde",
			},
			err: true,
		},
		{
			desc: "error, length of salt too long",
			msg: &v1.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			},
			err: true,
		},
		{
			desc: "error, invalid pubkey",
			msg: &v1.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				PubKey: &codectypes.Any{
					TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
					Value:   []byte(nil),
				},
			},
			err: true,
		},
		{
			desc: "error, codeID is zero",
			msg: &v1.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				PubKey:         pubKey,
				CodeID:         uint64(0),
			},
			err: true,
		},
		{
			desc: "error, invalid json msg",
			msg: &v1.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				PubKey:         pubKey,
				CodeID:         uint64(1),
				InitMsg:        []byte("{]"),
			},
			err: true,
		},
		{
			desc: "validate basic successfully",
			msg: &v1.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				PubKey:         pubKey,
				CodeID:         uint64(1),
				InitMsg:        []byte("{}"),
			},
			err: false,
		},
	} {
		err := tc.msg.ValidateBasic()

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestRecoverValidateBasic(t *testing.T) {
	pubKey, err := v1.PubKeyToAny(app.MakeEncodingConfig().Marshaler, helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *v1.MsgRecover
		err  bool
	}{
		{
			desc: "error, creator address invalid bench32 string",
			msg: &v1.MsgRecover{
				Creator: "abcde",
			},
			err: true,
		},
		{
			desc: "error, account address invalid bench32 string",
			msg: &v1.MsgRecover{
				Creator: helper.UserAddr,
				Address: "abcde",
			},
			err: true,
		},
		{
			desc: "error, invalid pubkey",
			msg: &v1.MsgRecover{
				Creator: helper.UserAddr,
				Address: helper.UserAddr,
				PubKey: &codectypes.Any{
					TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
					Value:   []byte(nil),
				},
			},
			err: true,
		},
		{
			desc: "error, credentials invalid base64 string",
			msg: &v1.MsgRecover{
				Creator:     helper.UserAddr,
				Address:     helper.UserAddr,
				PubKey:      pubKey,
				Credentials: "abcde",
			},
			err: true,
		},
		{
			desc: "validate basic successfully",
			msg: &v1.MsgRecover{
				Creator:     helper.UserAddr,
				Address:     helper.UserAddr,
				PubKey:      pubKey,
				Credentials: "eyIifQ==",
			},
			err: false,
		},
	} {
		err := tc.msg.ValidateBasic()

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
