package keeper_test

import (
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	helper "github.com/aura-nw/aura/tests/smartaccount"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIncrementNextAccountID(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	accID := keeper.GetAndIncrementNextAccountID(ctx)
	require.Equal(t, typesv1.DefaultSmartAccountId, accID)

	newAccID := keeper.GetNextAccountID(ctx)
	require.Equal(t, typesv1.DefaultSmartAccountId+1, newAccID)
}

func TestGetSetDelSignerAddress(t *testing.T) {
	testAddress := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	acc, err := sdk.AccAddressFromBech32(testAddress)
	require.NoError(t, err)

	keeper.SetSignerAddress(ctx, acc)
	getAcc := keeper.GetSignerAddress(ctx)
	require.Equal(t, acc, getAcc)

	keeper.DeleteSignerAddress(ctx)
	getAcc = keeper.GetSignerAddress(ctx)
	require.Equal(t, sdk.AccAddress(nil), getAcc)
}

func TestGasRemaining(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	gasRemaining := uint64(100)

	require.Equal(t, keeper.HasGasRemaining(ctx), false)

	keeper.SetGasRemaining(ctx, gasRemaining)
	require.Equal(t, keeper.HasGasRemaining(ctx), true)

	gas := keeper.GetGasRemaining(ctx)
	require.Equal(t, gasRemaining, gas)

	keeper.DeleteGasRemaining(ctx)
	require.Equal(t, keeper.HasGasRemaining(ctx), false)
}

func TestValidateActivateSA(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	// add new base account to chain
	err = helper.AddNewBaseAccount(app, ctx, helper.UserAddr, nil, uint64(0))
	require.NoError(t, err)

	for _, tc := range []struct {
		desc           string
		accountAddress string
		codeID         uint64
		err            bool
	}{
		{
			desc:           "error, codeID not in whitelist",
			accountAddress: helper.UserAddr,
			codeID:         2, // not whitelist codeID
			err:            true,
		},
		{
			desc:           "error, invalid bench32 string for account address",
			accountAddress: "", // invalid bench32
			codeID:         1,
			err:            true,
		},
		{
			desc:           "validate activate smartaccount successfully",
			accountAddress: helper.UserAddr,
			codeID:         1,
			err:            false,
		},
	} {

		msg := &typesv1.MsgActivateAccount{
			AccountAddress: tc.accountAddress,
			CodeID:         tc.codeID,
			Salt:           helper.DefaultSalt,
			InitMsg:        helper.DefaultMsg,
			PubKey:         pubKey,
		}

		_, err := keeper.ValidateActiveSA(ctx, msg)

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

}

func TestPrepareBeforeActive(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(app, ctx, helper.UserAddr, nil, 1)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(helper.UserAddr)
	require.NoError(t, err)
	acc := app.AccountKeeper.GetAccount(ctx, accAddr)
	require.Equal(t, uint64(1), acc.GetSequence())

	// PrepareBeforeActive et sequence of account to zero
	err = keeper.PrepareBeforeActive(ctx, acc)
	require.NoError(t, err)

	acc = app.AccountKeeper.GetAccount(ctx, accAddr)
	require.Equal(t, uint64(0), acc.GetSequence())
}

func TestActiveSmartAccount(t *testing.T) {
	for _, tc := range []struct {
		desc           string
		AccountAddress string
		codeID         uint64
		err            bool
	}{
		{
			desc:           "error, codeID not exist on chain",
			codeID:         2, // not existed codeID
			AccountAddress: "",
			err:            true,
		},
		{
			desc:           "error, smartaccount with address not found",
			AccountAddress: "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28", // not smartaccount address
			codeID:         1,
			err:            true,
		},
		{
			desc:           "active smartaccount successfully",
			codeID:         1,
			AccountAddress: "",
			err:            false,
		},
	} {
		ctx, app := helper.SetupGenesisTest(t)

		keeper := app.SaKeeper

		acc, pub, err := helper.GenerateInActivateAccount(
			app,
			ctx,
			helper.WasmPath2+"base.wasm",
			helper.DefaultPubKey,
			uint64(1),
			helper.DefaultSalt,
			helper.DefaultMsg,
		)
		require.NoError(t, err)

		accAddr := acc.GetAddress().String()
		if tc.AccountAddress != "" {
			accAddr = tc.AccountAddress
		}

		msg := &typesv1.MsgActivateAccount{
			AccountAddress: accAddr,
			CodeID:         tc.codeID,
			Salt:           helper.DefaultSalt,
			InitMsg:        helper.DefaultMsg,
			PubKey:         pub,
		}

		_, err = keeper.ActiveSmartAccount(ctx, msg, acc)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestHandleAfterActive(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(app, ctx, helper.UserAddr, nil, 1)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(helper.UserAddr)
	require.NoError(t, err)
	acc := app.AccountKeeper.GetAccount(ctx, accAddr)
	require.Equal(t, uint64(1), acc.GetSequence())
	require.Equal(t, nil, acc.GetPubKey())

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(t, err)

	// PrepareBeforeActive et sequence of account to zero
	backupSe := uint64(2)
	err = keeper.HandleAfterActive(ctx, acc, backupSe, dPubKey)
	require.NoError(t, err)

	acc = app.AccountKeeper.GetAccount(ctx, accAddr)
	require.Equal(t, uint64(2), acc.GetSequence())
	require.Equal(t, dPubKey, acc.GetPubKey())
}

func TestValidateRecoverSA(t *testing.T) {
	defaultCredentials := ""
	testSAAddress := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"
	testBAAddress := "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3"

	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	err := helper.AddNewSmartAccount(app, ctx, testSAAddress, nil, 0)
	require.NoError(t, err)

	err = helper.AddNewBaseAccount(app, ctx, testBAAddress, nil, 0)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc           string
		AccountAddress string
		err            bool
	}{
		{
			desc:           "error, account address invalid bench32 string",
			AccountAddress: "", // invalide bench32 string
			err:            true,
		},
		{
			desc:           "error, smartaccount with address not found",
			AccountAddress: testBAAddress, // not smartaccount address
			err:            true,
		},
		{
			desc:           "validate recover smartaccount successfully",
			AccountAddress: testSAAddress,
			err:            false,
		},
	} {
		pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
		require.NoError(t, err)

		msg := &typesv1.MsgRecover{
			Creator:     helper.UserAddr,
			Address:     tc.AccountAddress,
			PubKey:      pubKey,
			Credentials: defaultCredentials,
		}

		_, err = keeper.ValidateRecoverSA(ctx, msg)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestCallSMValidate(t *testing.T) {
	customMsg := []byte("{\"recover_key\":\"024ab33b4f0808eba493ac4e3ead798c8339e2fd216b20ca110001fd094784c07f\"}")

	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	rPubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultRPubKery)
	require.NoError(t, err)
	dRPubKey, err := typesv1.PubKeyDecode(rPubKey)
	require.NoError(t, err)

	acc, _, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath2+"recovery.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		customMsg,
	)
	require.NoError(t, err)

	msg := &typesv1.MsgActivateAccount{
		AccountAddress: acc.GetAddress().String(),
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        customMsg,
		PubKey:         pubKey,
	}
	_, err = keeper.ActiveSmartAccount(ctx, msg, acc)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc        string
		credentials string
		err         bool
	}{
		{
			desc:        "error, credentials invalid base64 string",
			credentials: "abc6432", // invalid base64
			err:         true,
		},
		{
			desc:        "error, wrong credentials",
			credentials: "eyJ0ZXN0IjoidGVzdCJ9", // wrong credentials
			err:         true,
		},
		{
			desc:        "validate for smartaccount recover successfully",
			credentials: "eyJzaWduYXR1cmUiOls4LDI0NywxOTksMTM4LDIzOCwxOTQsMTI5LDI1NCwyNTEsMTMxLDIzNywyNDEsMzMsODcsMTAzLDQyLDEzOCwyMjcsMjM3LDEyMyw5MiwyMjYsNjMsMTc0LDIwMSw2OCwyMSwzMiw5OSwxMzEsMjM1LDIzMSwyOCwxNzAsMjAzLDE4MCwxMTEsMiwyMjAsMTI2LDE0NCwxNzQsMTYxLDkyLDI1LDIwMiw2MiwxODEsMjUyLDE3OCwxNjMsNDAsMTc3LDIxMCwxNzYsNSwxNDUsMjAwLDU0LDE5MiwxMDgsMyw3Nyw2MV19",
			err:         false,
		},
	} {

		msg := &typesv1.MsgRecover{
			Creator:     helper.UserAddr,
			Address:     acc.GetAddress().String(),
			PubKey:      rPubKey,
			Credentials: tc.credentials,
		}

		err = keeper.CallSMValidate(ctx, msg, acc.GetAddress(), dRPubKey)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestIsInactiveAccount(t *testing.T) {
	testBAAddress1 := "cosmos1kzlrmxw3h2n4uzuv73m33cfw7xt7qjf3hlqx33ulc02e9dhxu46qgfxg9l"
	testBAAddress2 := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(t, err)

	err = helper.AddNewBaseAccount(app, ctx, testBAAddress1, nil, 0)
	require.NoError(t, err)

	err = helper.AddNewBaseAccount(app, ctx, testBAAddress2, dPubKey, 0)
	require.NoError(t, err)

	acc, _, err := helper.GenerateInActivateAccount(
		app,
		ctx,
		helper.WasmPath2+"base.wasm",
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(t, err)

	msg := &typesv1.MsgActivateAccount{
		AccountAddress: acc.GetAddress().String(),
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey,
	}
	_, err = keeper.ActiveSmartAccount(ctx, msg, acc)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc           string
		AccountAddress string
		err            bool
	}{
		{
			desc:           "error, smart account with address not found",
			AccountAddress: "cosmos1st3fng2vjcpz5lhg46un94zg0vn3nj658wc0chc57z29hx8zqeys6mvxdd", // not existed account
			err:            true,
		},
		{
			desc:           "is inactive smartaccount, baseaccount without pubkey",
			AccountAddress: testBAAddress1,
			err:            false,
		},
		{
			desc:           "error, baseaccount with pubkey",
			AccountAddress: testBAAddress2, // base account has pubkey
			err:            true,
		},
		{
			desc:           "error, smartaccount with linked smartcontract",
			AccountAddress: acc.GetAddress().String(), // smart contract created
			err:            true,
		},
	} {
		acc, err := sdk.AccAddressFromBech32(tc.AccountAddress)
		require.NoError(t, err)

		_, err = keeper.IsInactiveAccount(ctx, acc)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestGetSmartAccountByAddress(t *testing.T) {
	testAddress1 := "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3"
	testAddress2 := "cosmos1kzlrmxw3h2n4uzuv73m33cfw7xt7qjf3hlqx33ulc02e9dhxu46qgfxg9l"
	testAddress3 := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	ctx, app := helper.SetupGenesisTest(t)

	err := helper.AddNewSmartAccount(app, ctx, testAddress1, nil, 0)
	require.NoError(t, err)

	err = helper.AddNewBaseAccount(app, ctx, testAddress2, nil, 0)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc           string
		AccountAddress string
		err            bool
		exist          bool
	}{
		{
			desc:           "get smartaccount sucessfully",
			AccountAddress: testAddress1,
			err:            false,
			exist:          true,
		},
		{
			desc:           "is baseaccount, got nil",
			AccountAddress: testAddress2,
			err:            false,
		},
		{
			desc:           "error, account with address not found",
			AccountAddress: testAddress3,
			err:            true,
		},
	} {
		acc, err := sdk.AccAddressFromBech32(tc.AccountAddress)
		require.NoError(t, err)

		saAcc, err := app.SaKeeper.GetSmartAccountByAddress(ctx, acc)

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		if tc.exist {
			require.NotEqual(t, (*typesv1.SmartAccount)(nil), saAcc)
		} else {
			require.Equal(t, (*typesv1.SmartAccount)(nil), saAcc)
		}
	}
}

func TestCheckAllowedMsgs(t *testing.T) {
	ctx, app := helper.SetupGenesisTest(t)

	keeper := app.SaKeeper

	params := typesv1.Params{
		DisableMsgsList: []string{
			"/cosmwasm.wasm.v1.MsgUpdateAdmin",
			"/cosmwasm.wasm.v1.MsgClearAdmin",
		},
		MaxGasExecute: 2000000,
	}
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc string
		msgs []sdk.Msg
		err  bool
	}{
		{
			desc: "allowed msgs",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{},
			},
			err: false,
		},
		{
			desc: "note allowed msgs",
			msgs: []sdk.Msg{
				&typesv1.MsgActivateAccount{},
				&wasmtypes.MsgUpdateAdmin{},
			},
			err: true,
		},
	} {
		err := keeper.CheckAllowedMsgs(ctx, tc.msgs)

		if !tc.err {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
