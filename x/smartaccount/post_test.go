package smartaccount_test

import (
	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/keeper"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestAfterTxDecorator() {

	keybase := keyring.NewInMemory(s.App.AppCodec())

	newAcc, pubKey, err := helper.GenerateInActivateAccount(
		s.App,
		s.ctx,
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(s.T(), err)
	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(s.T(), err)

	msgServer := keeper.NewMsgServerImpl(s.App.SaKeeper)

	/* ======== activate smart account ======== */
	msg := &typesv1.MsgActivateAccount{
		AccountAddress: newAcc.Address,
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey,
	}

	// activate account
	res, err := msgServer.ActivateAccount(sdk.WrapSDKContext(s.ctx), msg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), newAcc.String(), res.Address)

	accSigner, err := makeMockAccount(keybase, "test")
	require.NoError(s.T(), err)
	err = accSigner.SetPubKey(dPubKey)
	require.NoError(s.T(), err)
	s.App.AccountKeeper.SetAccount(s.ctx, accSigner)

	signer := Signer{
		keyName:        "test",
		acc:            accSigner,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	acc1, err := makeMockAccount(keybase, "test1")
	require.NoError(s.T(), err)

	for _, tc := range []struct {
		desc    string
		msgs    []sdk.Msg
		signers []Signer
		err     bool
		isSa    bool
	}{
		{
			desc: "success post handler for smartaccount tx",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(newAcc.GetAddress(), acc1.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer},
			isSa:    true,
			err:     false,
		},
		{
			desc: "not smartaccount tx",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(newAcc.GetAddress(), acc1.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer},
			isSa:    false,
			err:     false,
		},
	} {
		if tc.isSa {
			s.App.SaKeeper.SetSignerAddress(s.ctx, newAcc.GetAddress())
		}

		sigTx, err := prepareTx(s.ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(s.T(), err)

		saat := smartaccount.NewAfterTxDecorator(s.App.SaKeeper)
		_, err = saat.PostHandle(s.ctx, sigTx, false, true, DefaultPostHandler())

		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}
}

func (s *TestSuite) TestPostValidateAuthzTxDecorator() {

	keybase := keyring.NewInMemory(s.App.AppCodec())

	newAcc, pubKey, err := helper.GenerateInActivateAccount(
		s.App,
		s.ctx,
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(s.T(), err)
	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(s.T(), err)

	msgServer := keeper.NewMsgServerImpl(s.App.SaKeeper)

	/* ======== activate smart account ======== */
	msg := &typesv1.MsgActivateAccount{
		AccountAddress: newAcc.Address,
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey,
	}

	// activate account
	res, err := msgServer.ActivateAccount(sdk.WrapSDKContext(s.ctx), msg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), newAcc.String(), res.Address)

	accSigner, err := makeMockAccount(keybase, "test")
	require.NoError(s.T(), err)
	err = accSigner.SetPubKey(dPubKey)
	require.NoError(s.T(), err)
	s.App.AccountKeeper.SetAccount(s.ctx, accSigner)

	signer := Signer{
		keyName:        "test",
		acc:            accSigner,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	acc1, err := makeMockAccount(keybase, "test1")
	require.NoError(s.T(), err)

	anyBankSend, err := codectypes.NewAnyWithValue(banktypes.NewMsgSend(newAcc.GetAddress(), acc1.GetAddress(), sdk.Coins{}))
	require.NoError(s.T(), err)

	anyMsgExec, err := codectypes.NewAnyWithValue(&authz.MsgExec{
		Grantee: signer.acc.GetAddress().String(),
		Msgs: []*codectypes.Any{
			anyBankSend,
		},
	})
	require.NoError(s.T(), err)

	for _, tc := range []struct {
		desc         string
		msgs         []sdk.Msg
		signers      []Signer
		setGas       bool
		gasRemaining bool
		err          bool
	}{
		{
			desc: "validate one level nested smartaccount success",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: signer.acc.GetAddress().String(),
					Msgs:    []*codectypes.Any{anyBankSend},
				},
			},
			signers: []Signer{signer},
			err:     false,
		},
		{
			desc: "validate multi level nested smartaccount success",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: signer.acc.GetAddress().String(),
					Msgs: []*codectypes.Any{
						anyMsgExec,
					},
				},
			},
			signers: []Signer{signer},
			err:     false,
		},
		{
			desc: "panic, validate nested smartaccount fail with out of gas because of max exec gas too low",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: signer.acc.GetAddress().String(),
					Msgs: []*codectypes.Any{
						anyMsgExec,
					},
				},
			},
			signers: []Signer{signer},
			setGas:  true,
			err:     true,
		},
		{
			desc: "panic, validate nested smartaccount fail with not enough gas remaining",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: signer.acc.GetAddress().String(),
					Msgs: []*codectypes.Any{
						anyMsgExec,
					},
				},
			},
			signers:      []Signer{signer},
			gasRemaining: true,
			err:          true,
		},
	} {
		if tc.setGas {
			params := typesv1.Params{
				MaxGasExecute: 1000,
			}
			err = s.App.SaKeeper.SetParams(s.ctx, params)
			require.NoError(s.T(), err)
		}

		if tc.gasRemaining {
			gasRemaining := uint64(1000)
			s.App.SaKeeper.SetGasRemaining(s.ctx, gasRemaining)
		}

		sigTx, err := prepareTx(s.ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(s.T(), err)

		saat := smartaccount.NewPostValidateAuthzTxDecorator(s.App.SaKeeper)

		if !tc.err {
			require.NotPanics(s.T(), func() {
				_, err = saat.PostHandle(s.ctx, sigTx, false, true, DefaultPostHandler())
				require.NoError(s.T(), err)
			})
		} else {
			require.Panics(s.T(), func() {
				_, err = saat.PostHandle(s.ctx, sigTx, false, true, DefaultPostHandler())
			})
		}
	}
}

func DefaultPostHandler() sdk.PostHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (newCtx sdk.Context, err error) {
		return ctx, nil
	}
}
