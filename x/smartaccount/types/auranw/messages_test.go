package auranw_test

import (
	"testing"

	"github.com/aura-nw/aura/app"
	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/types/auranw"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestActivateAccountValidateBasic(t *testing.T) {
	pubKey, err := auranw.PubKeyToAny(app.MakeEncodingConfig().Marshaler, helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *auranw.MsgActivateAccount
		err  bool
	}{
		{
			desc: "error, account address invalid bench32 string",
			msg: &auranw.MsgActivateAccount{
				AccountAddress: "abcde",
			},
			err: true,
		},
		{
			desc: "error, length of salt too long",
			msg: &auranw.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			},
			err: true,
		},
		{
			desc: "error, invalid pubkey",
			msg: &auranw.MsgActivateAccount{
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
			msg: &auranw.MsgActivateAccount{
				AccountAddress: helper.UserAddr,
				Salt:           []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				PubKey:         pubKey,
				CodeID:         uint64(0),
			},
			err: true,
		},
		{
			desc: "error, invalid json msg",
			msg: &auranw.MsgActivateAccount{
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
			msg: &auranw.MsgActivateAccount{
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
	pubKey, err := auranw.PubKeyToAny(app.MakeEncodingConfig().Marshaler, helper.DefaultPubKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msg  *auranw.MsgRecover
		err  bool
	}{
		{
			desc: "error, creator address invalid bench32 string",
			msg: &auranw.MsgRecover{
				Creator: "abcde",
			},
			err: true,
		},
		{
			desc: "error, account address invalid bench32 string",
			msg: &auranw.MsgRecover{
				Creator: helper.UserAddr,
				Address: "abcde",
			},
			err: true,
		},
		{
			desc: "error, invalid pubkey",
			msg: &auranw.MsgRecover{
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
			msg: &auranw.MsgRecover{
				Creator:     helper.UserAddr,
				Address:     helper.UserAddr,
				PubKey:      pubKey,
				Credentials: "abcde",
			},
			err: true,
		},
		{
			desc: "validate basic successfully",
			msg: &auranw.MsgRecover{
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
