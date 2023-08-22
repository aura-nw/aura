package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
)

// ------------------------------ ActivateAccount ------------------------------
func TestActivateAccount(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	for _, tc := range []struct {
		accountAddress string
		codeID         uint64
		err            bool
	}{
		{
			// error msg
			accountAddress: "",
			codeID:         2, // not whitelist codeID
			err:            true,
		},
		{
			// error msg
			accountAddress: "aura1dkgyvk8zfn5vqg40qw0rhk972ugjppaeenqclwa6f0nsvzmx8mmsnggzpx", // not inactivate smartaccount address
			codeID:         1,
			err:            true,
		},
		{
			// activate succeed
			accountAddress: "",
			codeID:         1,
			err:            false,
		},
	} {
		ctx, app := helper.SetupGenesisTest()

		newAcc, pubKey, err := helper.GenerateInActivateAccount(
			app,
			ctx,
			helper.WasmPath2+"base.wasm",
			helper.DefaultPubKey,
			helper.DefaultCodeID,
			helper.DefaultSalt,
			helper.DefaultMsg,
		)
		require.NoError(t, err)

		/* ======== activate smart account ======== */
		msgServer := keeper.NewMsgServerImpl(app.SaKeeper)

		// smartaccount address
		accAddr := newAcc.Address
		if tc.accountAddress != "" {
			accAddr = tc.accountAddress
		}

		msg := &types.MsgActivateAccount{
			AccountAddress: accAddr,
			CodeID:         tc.codeID,
			Salt:           helper.DefaultSalt,
			InitMsg:        helper.DefaultMsg,
			PubKey:         pubKey,
		}

		// activate account
		res, err := msgServer.ActivateAccount(sdk.WrapSDKContext(ctx), msg)

		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, newAcc.String(), res.Address)

			// must be smartaccount type
			saAccAddr, err := sdk.AccAddressFromBech32(accAddr)
			require.NoError(t, err)

			saAccount := app.AccountKeeper.GetAccount(ctx, saAccAddr)

			_, ok := saAccount.(*types.SmartAccount)
			require.Equal(t, true, ok)
		}
	}
}

// ------------------------------ RecoverAccount ------------------------------
func TestRecoverAccount(t *testing.T) {
	customMsg := []byte("{\"recover_key\":\"024ab33b4f0808eba493ac4e3ead798c8339e2fd216b20ca110001fd094784c07f\"}")

	sdk.GetConfig().SetBech32PrefixForAccount(helper.AccountAddressPrefix, "")

	for _, tc := range []struct {
		desc           string
		accountAddress string
		credentials    string
		err            bool
	}{
		{
			desc:           "error, invalid credentials",
			accountAddress: "",
			credentials:    "eyJ0ZXN0IjoidGVzdCJ9", // invalid credentials
			err:            true,
		},
		{
			desc:           "error, smartaccount not activated",
			accountAddress: "aura1dkgyvk8zfn5vqg40qw0rhk972ugjppaeenqclwa6f0nsvzmx8mmsnggzpx", // not activated smartaccount address
			err:            true,
		},
		{
			desc:           "recover smartaccount successfully",
			accountAddress: "",
			credentials:    "eyJzaWduYXR1cmUiOls4LDI0NywxOTksMTM4LDIzOCwxOTQsMTI5LDI1NCwyNTEsMTMxLDIzNywyNDEsMzMsODcsMTAzLDQyLDEzOCwyMjcsMjM3LDEyMyw5MiwyMjYsNjMsMTc0LDIwMSw2OCwyMSwzMiw5OSwxMzEsMjM1LDIzMSwyOCwxNzAsMjAzLDE4MCwxMTEsMiwyMjAsMTI2LDE0NCwxNzQsMTYxLDkyLDI1LDIwMiw2MiwxODEsMjUyLDE3OCwxNjMsNDAsMTc3LDIxMCwxNzYsNSwxNDUsMjAwLDU0LDE5MiwxMDgsMyw3Nyw2MV19",
			err:            false,
		},
	} {
		ctx, app := helper.SetupGenesisTest()

		newAcc, pubKey, err := helper.GenerateInActivateAccount(
			app,
			ctx,
			helper.WasmPath2+"recovery.wasm",
			helper.DefaultPubKey,
			helper.DefaultCodeID,
			helper.DefaultSalt,
			customMsg,
		)
		require.NoError(t, err)

		msgServer := keeper.NewMsgServerImpl(app.SaKeeper)

		/* ======== activate smart account ======== */
		msg := &types.MsgActivateAccount{
			AccountAddress: newAcc.Address,
			CodeID:         helper.DefaultCodeID,
			Salt:           helper.DefaultSalt,
			InitMsg:        customMsg,
			PubKey:         pubKey,
		}

		// activate account
		res, err := msgServer.ActivateAccount(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.Equal(t, newAcc.String(), res.Address)

		/* ======== Recover ======== */
		accAddr := newAcc.Address
		if tc.accountAddress != "" {
			accAddr = tc.accountAddress
		}

		rPubKey, err := types.PubKeyToAny(app.AppCodec(), helper.DefaultRPubKery)
		require.NoError(t, err)

		recoverMsg := &types.MsgRecover{
			Creator:     userAddr,
			Address:     accAddr,
			PubKey:      rPubKey,
			Credentials: tc.credentials,
		}

		_, rErr := msgServer.Recover(sdk.WrapSDKContext(ctx), recoverMsg)

		saAccAddr, err := sdk.AccAddressFromBech32(newAcc.Address)
		require.NoError(t, err)

		saAccount := app.AccountKeeper.GetAccount(ctx, saAccAddr)

		_, ok := saAccount.(*types.SmartAccount)
		require.Equal(t, true, ok)

		if tc.err {
			dPubKey, err := types.PubKeyDecode(pubKey)
			require.NoError(t, err)

			require.Equal(t, saAccount.GetPubKey(), dPubKey)
			require.Error(t, rErr)
		} else {
			rPubKey, err := types.PubKeyDecode(rPubKey)
			require.NoError(t, err)

			require.Equal(t, saAccount.GetPubKey(), rPubKey)
			require.NoError(t, rErr)

		}
	}
}
