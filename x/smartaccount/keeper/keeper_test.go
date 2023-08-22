package keeper_test

import (
	"testing"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	userAddr = "aura17n858c4urvenkf7edjs4uumudej3ekyv432e34"
)

func TestIncrementNextAccountID(t *testing.T) {
	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	accID := keeper.GetAndIncrementNextAccountID(ctx)
	require.Equal(t, types.DefaultSmartAccountId, accID)

	newAccID := keeper.GetNextAccountID(ctx)
	require.Equal(t, types.DefaultSmartAccountId+1, newAccID)
}

func TestValidateActivateSA(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	// add new base account to chain
	err = helper.AddNewBaseAccount(app, ctx, userAddr, nil, uint64(0))
	require.NoError(t, err)

	for _, tc := range []struct {
		desc           string
		accountAddress string
		codeID         uint64
		err            bool
	}{
		{
			desc:           "error, codeID not in whitelist",
			accountAddress: userAddr,
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
			accountAddress: userAddr,
			codeID:         1,
			err:            false,
		},
	} {

		msg := &types.MsgActivateAccount{
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
	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(app, ctx, userAddr, nil, 1)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(userAddr)
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
	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

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
			AccountAddress: "aura1dkgyvk8zfn5vqg40qw0rhk972ugjppaeenqclwa6f0nsvzmx8mmsnggzpx", // not smartaccount address
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
		ctx, app := helper.SetupGenesisTest()

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

		msg := &types.MsgActivateAccount{
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
	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(app, ctx, userAddr, nil, 1)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(userAddr)
	require.NoError(t, err)
	acc := app.AccountKeeper.GetAccount(ctx, accAddr)
	require.Equal(t, uint64(1), acc.GetSequence())
	require.Equal(t, nil, acc.GetPubKey())

	// prepare pubkey
	pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	dPubKey, err := types.PubKeyDecode(pubKey)
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
	testSAAddress := "aura19nywm4al2pc2sdj834gtdm6tvcn5kqpghlwd022tvld0hek4jfeslshhzj"
	testBAAddress := "aura1dkgyvk8zfn5vqg40qw0rhk972ugjppaeenqclwa6f0nsvzmx8mmsnggzpx"

	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

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
		pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
		require.NoError(t, err)

		msg := &types.MsgRecover{
			Creator:     userAddr,
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

	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	// prepare pubkey
	pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	rPubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultRPubKery)
	require.NoError(t, err)
	dRPubKey, err := types.PubKeyDecode(rPubKey)
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

	msg := &types.MsgActivateAccount{
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

		msg := &types.MsgRecover{
			Creator:     userAddr,
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
	testSAAddress1 := "aura19nywm4al2pc2sdj834gtdm6tvcn5kqpghlwd022tvld0hek4jfeslshhzj"
	testBAAddress1 := "aura19nywm4al2pc2sdj834gtdm6tvcn5kqpghlwd022tvld0hek4jfeslshhzj"
	testBAAddress2 := "aura1dkgyvk8zfn5vqg40qw0rhk972ugjppaeenqclwa6f0nsvzmx8mmsnggzpx"

	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	ctx, app := helper.SetupGenesisTest()

	keeper := app.SaKeeper

	// prepare pubkey
	pubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultPubKey)
	require.NoError(t, err)

	dPubKey, err := types.PubKeyDecode(pubKey)
	require.NoError(t, err)

	err = helper.AddNewSmartAccount(app, ctx, testSAAddress1, dPubKey, 0)
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

	msg := &types.MsgActivateAccount{
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
			AccountAddress: "aura1qvw635lxp7y5dgnc5k0zhpxr3cmwlpfmk6rj8ku7s0dz6597r5zq9h7wh9", // not existed account
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
			desc:           "is inactive smartaccount, smartaccount without linked smartcontract",
			AccountAddress: testSAAddress1,
			err:            false,
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
