package v1_test

import (
	"testing"

	"github.com/aura-nw/aura/app"
	v1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestPubkeyToAny(t *testing.T) {
	for _, tc := range []struct {
		desc string
		raw  string
		err  bool
	}{
		{
			desc: "error, empty pubkey",
			raw:  "", // error pubkey string
			err:  true,
		},
		{
			desc: "convert pubkey string to type Any successfully",
			raw:  "{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A/2t0ru/iZ4HoiX0DkTuMy9rC2mMeXmiN6luM3pa+IvT\"}",
			err:  false,
		},
	} {

		_, err := v1.PubKeyToAny(app.MakeEncodingConfig().Marshaler, []byte(tc.raw))
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestPubKeyDecode(t *testing.T) {
	_, err := v1.PubKeyDecode(nil)
	require.Error(t, err)

	pubKey := &codectypes.Any{
		TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
		Value:   []byte(nil),
	}
	_, err = v1.PubKeyDecode(pubKey)
	require.Error(t, err)

	raw := "{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A/2t0ru/iZ4HoiX0DkTuMy9rC2mMeXmiN6luM3pa+IvT\"}"
	any, err := v1.PubKeyToAny(app.MakeEncodingConfig().Marshaler, []byte(raw))
	require.NoError(t, err)
	dPubKey, err := v1.PubKeyDecode(any)
	require.NoError(t, err)
	require.Equal(t, any.Value[2:], dPubKey.Bytes())
}
