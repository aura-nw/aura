package smartaccount_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tests "github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/keeper"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/authz"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestGetSmartAccountTxSigner(t *testing.T) {
	var (
		app     = tests.Setup(t, false)
		ctx     = app.NewContext(false, tmproto.Header{})
		keybase = keyring.NewInMemory(app.AppCodec())
	)

	acc1, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc1)

	acc2Mock, err := makeMockAccount(keybase, "test2")
	require.NoError(t, err)
	acc2 := typesv1.NewSmartAccountFromAccount(acc2Mock)
	err = acc2.SetPubKey(acc2Mock.GetPubKey())
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc2)
	require.NoError(t, err)

	acc3, err := makeMockAccount(keybase, "test3")
	require.NoError(t, err)

	signer1 := Signer{
		keyName:        "test1",
		acc:            acc1,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}
	signer2 := Signer{
		keyName:        "test2",
		acc:            acc2,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}
	signer3 := Signer{
		keyName:        "test3",
		acc:            acc3,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	for _, tc := range []struct {
		desc    string
		msgs    []sdk.Msg
		signers []Signer
		expIs   bool
		err     bool
	}{
		{
			desc: "tx has one signer and it is an SmartAccount",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc2.GetAddress(), acc1.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer2},
			expIs:   true,
		},
		{
			desc: "tx has one signer but it's not an SmartAccount",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1},
			expIs:   false,
		},
		{
			desc: "tx has a signer but it doesn't exist on the chain yet",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc3.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer3},
			expIs:   false,
			err:     true,
		},
		{
			desc: "tx has more than one signers",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				banktypes.NewMsgSend(acc2.GetAddress(), acc1.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1, signer2},
			expIs:   false,
		},
	} {
		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		signerAcc, err := smartaccount.GetSmartAccountTxSigner(ctx, sigTx, app.SaKeeper)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		if tc.expIs {
			require.NotEqual(t, (*typesv1.SmartAccount)(nil), signerAcc)
		} else {
			require.Equal(t, (*typesv1.SmartAccount)(nil), signerAcc)
		}
	}
}

func TestGetValidActivateAccountMessage(t *testing.T) {
	var (
		app     = tests.Setup(t, false)
		ctx     = app.NewContext(false, tmproto.Header{})
		keybase = keyring.NewInMemory(app.AppCodec())
	)

	acc1, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc1)

	acc2, err := makeMockAccount(keybase, "test2")
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc2)

	signer1 := Signer{
		keyName:        "test1",
		acc:            acc1,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}
	signer2 := Signer{
		keyName:        "test2",
		acc:            acc2,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	for _, tc := range []struct {
		desc    string
		msgs    []sdk.Msg
		signers []Signer
		expIs   bool
		err     bool
	}{
		{
			desc: "tx has one signer and it is an SmartAccount",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{
					AccountAddress: acc2.GetAddress().String(),
				},
			},
			signers: []Signer{signer2},
			expIs:   true,
		},
		{
			desc: "tx has one signer but it's not an SmartAccount",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1},
			expIs:   false,
		},
		{
			desc: "tx has more than one signers",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				banktypes.NewMsgSend(acc2.GetAddress(), acc1.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1, signer2},
			expIs:   false,
		},
		{
			desc: "tx has more than one message and contain activate message",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{AccountAddress: acc1.GetAddress().String()},
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1},
			expIs:   false,
			err:     true,
		},
		{
			desc: "tx has more than one signers and contain activate message",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{AccountAddress: acc1.GetAddress().String()},
			},
			signers: []Signer{signer1, signer2},
			expIs:   false,
			err:     true,
		},
	} {
		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		aaMsg, err := smartaccount.GetValidActivateAccountMessage(sigTx)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		if tc.expIs {
			require.NotEqual(t, (*typesv1.MsgActivateAccount)(nil), aaMsg)
		} else {
			require.Equal(t, (*typesv1.MsgActivateAccount)(nil), aaMsg)
		}
	}
}

func TestSetPubKeyDecorator(t *testing.T) {
	var (
		app     = tests.Setup(t, false)
		ctx     = app.NewContext(false, tmproto.Header{})
		keybase = keyring.NewInMemory(app.AppCodec())
	)

	acc, pubKey, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(t, err)

	acc1, _, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		2,
		[]byte("test 2"),
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	err = helper.AddNewSmartAccount(app, ctx, acc1.GetAddress().String(), nil, 0)
	require.NoError(t, err)

	acc1Signer, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)
	err = acc1Signer.SetPubKey(dPubKey)
	require.NoError(t, err)

	acc2, err := makeMockAccount(keybase, "test2")
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc2)

	signer1 := Signer{
		keyName:        "test1",
		acc:            acc1Signer,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	signer2 := Signer{
		keyName:        "test2",
		acc:            acc2,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	for _, tc := range []struct {
		desc    string
		msgs    []sdk.Msg
		signers []Signer
		err     bool
		isSa    bool
	}{
		{
			desc: "is ActivateAccount tx",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{
					AccountAddress: acc.GetAddress().String(),
					CodeID:         helper.DefaultCodeID,
					Salt:           helper.DefaultSalt,
					InitMsg:        helper.DefaultMsg,
					PubKey:         pubKey,
				},
			},
			signers: []Signer{signer1},
			err:     false,
		},
		{
			desc: "is SmartAccount tx",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgExecuteContract{
					Sender:   acc.GetAddress().String(),
					Contract: acc.GetAddress().String(),
					Msg:      []byte("{\"after_execute\":{\"msgs\":[{\"type_url\":\"/cosmos.bank.v1beta1.MsgSend\",\"value\":\"{\\\"from_address\\\":\\\"" + acc.GetAddress().String() + "\\\",\\\"to_address\\\":\\\"" + acc2.GetAddress().String() + "\\\",\\\"amount\\\":[]}\"}]}}"),
				},
			},
			signers: []Signer{signer1},
			err:     false,
			isSa:    true,
		},
		{
			desc: "not ActivateAccount nor SmartAccount tx, just normal tx",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc2.GetAddress(), acc.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer2},
			err:     false,
		},
		{
			desc: "error, is SmartAccount tx but PubKey not yet set",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgExecuteContract{
					Sender:   acc1.GetAddress().String(),
					Contract: acc1.GetAddress().String(),
					Msg:      []byte("{\"after_execute\":{\"msgs\":[{\"type_url\":\"/cosmos.bank.v1beta1.MsgSend\",\"value\":\"{\\\"from_address\\\":\\\"" + acc1.GetAddress().String() + "\\\",\\\"to_address\\\":\\\"" + acc2.GetAddress().String() + "\\\",\\\"amount\\\":[]}\"}]}}"),
				},
			},
			signers: []Signer{signer1},
			err:     true,
		},
	} {
		if tc.isSa {
			err = helper.AddNewSmartAccount(app, ctx, acc.GetAddress().String(), dPubKey, 0)
			require.NoError(t, err)
		}

		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		sad := smartaccount.NewSetPubKeyDecorator(app.SaKeeper)
		_, err = sad.AnteHandle(ctx, sigTx, false, DefaultAnteHandler())

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestSmartAccountDecoratorForTx(t *testing.T) {
	var (
		ctx, app = helper.SetupGenesisTest(t)
		keybase  = keyring.NewInMemory(app.AppCodec())
	)

	params := typesv1.Params{
		WhitelistCodeID: []*typesv1.CodeID{
			{CodeID: 1, Status: true},
		},
		DisableMsgsList: []string{
			"/cosmwasm.wasm.v1.MsgUpdateAdmin",
			"/cosmwasm.wasm.v1.MsgClearAdmin",
		},
		MaxGasExecute: 2000000,
	}
	err := app.SaKeeper.SetParams(ctx, params)
	require.NoError(t, err)

	// base smartaccount
	acc1, pubKey1, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	dPubKey1, err := typesv1.PubKeyDecode(pubKey1)
	require.NoError(t, err)

	acc1Signer, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)
	err = acc1Signer.SetPubKey(dPubKey1)
	require.NoError(t, err)

	msg := &typesv1.MsgActivateAccount{
		AccountAddress: acc1.GetAddress().String(),
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey1,
	}

	msgServer := keeper.NewMsgServerImpl(app.SaKeeper)
	// activate account
	_, err = msgServer.ActivateAccount(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	acc2, err := makeMockAccount(keybase, "test2")
	require.NoError(t, err)

	signer1 := Signer{
		keyName:        "test1",
		acc:            acc1Signer,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	signer2 := Signer{
		keyName:        "test2",
		acc:            acc2,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	for _, tc := range []struct {
		desc     string
		msgs     []sdk.Msg
		signers  []Signer
		simulate bool
		err      bool
	}{
		{
			desc: "is SmartAccount tx, and validate success",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgExecuteContract{
					Sender:   acc1.GetAddress().String(),
					Contract: acc1.GetAddress().String(),
					Msg:      []byte("{\"after_execute\":{\"msgs\":[{\"type_url\":\"/cosmos.bank.v1beta1.MsgSend\",\"value\":\"{\\\"from_address\\\":\\\"" + acc1.GetAddress().String() + "\\\",\\\"to_address\\\":\\\"" + acc2.GetAddress().String() + "\\\",\\\"amount\\\":[]}\"}]}}"),
				},
			},
			signers: []Signer{signer1},
			err:     false,
		},
		{
			desc: "not SmartAccount tx, too many signer",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
			},
			signers: []Signer{signer1, signer2},
			err:     false,
		},
		{
			desc: "SmartAccount tx support simulate",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgExecuteContract{
					Sender:   acc1.GetAddress().String(),
					Contract: acc1.GetAddress().String(),
					Msg:      []byte("{\"after_execute\":{\"msgs\":[{\"type_url\":\"/cosmos.bank.v1beta1.MsgSend\",\"value\":\"{\\\"from_address\\\":\\\"" + acc1.GetAddress().String() + "\\\",\\\"to_address\\\":\\\"" + acc2.GetAddress().String() + "\\\",\\\"amount\\\":[]}\"}]}}"),
				},
			},
			signers:  []Signer{signer1},
			simulate: true,
			err:      false,
		},
		{
			desc: "error, not allowed msgs",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgUpdateAdmin{
					Sender:   acc1.Address,
					NewAdmin: acc1.Address,
					Contract: "cosmos1ztwdgj227nzrkgv0gxt0d3fx5q905ltjzwv5t9",
				},
			},
			signers: []Signer{signer1},
			err:     true,
		},
	} {
		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		satd := smartaccount.NewSmartAccountDecorator(app.SaKeeper)
		_, err = satd.AnteHandle(ctx, sigTx, tc.simulate, DefaultAnteHandler())

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestSmartAccountDecoratorForActivation(t *testing.T) {

	/* =================== test activate account message flow =================== */

	var (
		ctx, app = helper.SetupGenesisTest(t)
		keybase  = keyring.NewInMemory(app.AppCodec())
	)

	// base smartaccount
	acc1, pubKey1, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	dPubKey1, err := typesv1.PubKeyDecode(pubKey1)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc1)

	acc1Signer, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)
	err = acc1Signer.SetPubKey(dPubKey1)
	require.NoError(t, err)

	acc2, err := makeMockAccount(keybase, "test2")
	require.NoError(t, err)

	// setup module account
	acc3, pubKey3, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		2,
		[]byte("account3"),
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	dPubKey3, err := typesv1.PubKeyDecode(pubKey3)
	require.NoError(t, err)

	acc3Signer, err := makeMockAccount(keybase, "test3")
	require.NoError(t, err)
	err = acc3Signer.SetPubKey(dPubKey3)
	require.NoError(t, err)
	moduleAcc3 := authtypes.NewModuleAccount(acc3, "test", "hello")
	app.AccountKeeper.SetAccount(ctx, moduleAcc3)
	require.NoError(t, err)

	signer1 := Signer{
		keyName:        "test1",
		acc:            acc1Signer,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	signer3 := Signer{
		keyName:        "test3",
		acc:            acc3Signer,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	for _, tc := range []struct {
		desc     string
		msgs     []sdk.Msg
		signers  []Signer
		simulate bool
		err      bool
	}{
		{
			desc: "not is ActivateAccount message, just normal SmartAccount tx",
			msgs: []sdk.Msg{
				banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sdk.NewCoins()),
				&wasmtypes.MsgExecuteContract{
					Sender:   acc1.GetAddress().String(),
					Contract: acc1.GetAddress().String(),
					Msg:      []byte("{\"after_execute\":{\"msgs\":[{\"type_url\":\"/cosmos.bank.v1beta1.MsgSend\",\"value\":\"{\\\"from_address\\\":\\\"" + acc1.GetAddress().String() + "\\\",\\\"to_address\\\":\\\"" + acc2.GetAddress().String() + "\\\",\\\"amount\\\":[]}\"}]}}"),
				},
			},
			signers: []Signer{signer1},
			err:     false,
		},
		{
			desc: "is ActivateAccount message",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{
					AccountAddress: acc1.GetAddress().String(),
					CodeID:         helper.DefaultCodeID,
					Salt:           helper.DefaultSalt,
					InitMsg:        helper.DefaultMsg,
					PubKey:         pubKey1,
				},
			},
			signers: []Signer{signer1},
			err:     false,
		},
		{
			desc: "error, is ActivateAccount message but invalid signer",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{
					AccountAddress: acc3.GetAddress().String(),
					CodeID:         2,
					Salt:           []byte("account3"),
					InitMsg:        helper.DefaultMsg,
					PubKey:         pubKey1,
				},
			},
			signers: []Signer{signer3},
			err:     true,
		},
		{
			desc: "error, smartaccount address not the same as predicted",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{
					AccountAddress: acc1.GetAddress().String(),
					CodeID:         helper.DefaultCodeID,
					Salt:           []byte("custom salt"),
					InitMsg:        helper.DefaultMsg,
					PubKey:         pubKey1,
				},
			},
			signers: []Signer{signer1},
			err:     true,
		},
	} {
		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		satd := smartaccount.NewSmartAccountDecorator(app.SaKeeper)
		_, err = satd.AnteHandle(ctx, sigTx, tc.simulate, DefaultAnteHandler())

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestValidateAuthzTxDecorator(t *testing.T) {
	var (
		ctx, app = helper.SetupGenesisTest(t)
		keybase  = keyring.NewInMemory(app.AppCodec())
	)

	newAcc, pubKey, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath1+"base.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(t, err)
	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(t, err)

	msgServer := keeper.NewMsgServerImpl(app.SaKeeper)

	/* ======== activate smart account ======== */
	msg := &typesv1.MsgActivateAccount{
		AccountAddress: newAcc.Address,
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey,
	}

	// activate account
	res, err := msgServer.ActivateAccount(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.Equal(t, newAcc.String(), res.Address)

	accSigner, err := makeMockAccount(keybase, "test")
	require.NoError(t, err)
	err = accSigner.SetPubKey(dPubKey)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, accSigner)

	signer := Signer{
		keyName:        "test",
		acc:            accSigner,
		overrideAccNum: nil,
		overrideSeq:    nil,
	}

	acc1, err := makeMockAccount(keybase, "test1")
	require.NoError(t, err)

	anyBankSend, err := codectypes.NewAnyWithValue(banktypes.NewMsgSend(newAcc.GetAddress(), acc1.GetAddress(), sdk.Coins{}))
	require.NoError(t, err)

	anyMsgExec, err := codectypes.NewAnyWithValue(&authz.MsgExec{
		Grantee: signer.acc.GetAddress().String(),
		Msgs: []*codectypes.Any{
			anyBankSend,
		},
	})
	require.NoError(t, err)

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
			desc: "panic, validate nested smartaccount fail with out of gas",
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
			err = app.SaKeeper.SetParams(ctx, params)
			require.NoError(t, err)
		}

		if tc.gasRemaining {
			gasRemaining := uint64(1000)
			app.SaKeeper.SetGasRemaining(ctx, gasRemaining)
		}

		sigTx, err := prepareTx(ctx, keybase, tc.msgs, tc.signers, mockChainID, true)
		require.NoError(t, err)

		saat := smartaccount.NewValidateAuthzTxDecorator(app.SaKeeper)

		if !tc.err {
			require.NotPanics(t, func() {
				_, err = saat.AnteHandle(ctx, sigTx, false, DefaultAnteHandler())
				require.NoError(t, err)
			})
		} else {
			require.Panics(t, func() {
				_, err = saat.AnteHandle(ctx, sigTx, false, DefaultAnteHandler())
			})
		}
	}
}

func DefaultAnteHandler() sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		return ctx, nil
	}
}
