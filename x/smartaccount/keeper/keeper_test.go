package keeper_test

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	helper "github.com/aura-nw/aura/tests/smartaccount"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func (s *KeeperTestSuite) TestIncrementNextAccountID() {
	ctx, app := helper.SetupGenesisTest(s.T())

	keeper := app.SaKeeper

	accID := keeper.GetAndIncrementNextAccountID(ctx)
	require.Equal(s.T(), typesv1.DefaultSmartAccountId, accID)

	newAccID := keeper.GetNextAccountID(ctx)
	require.Equal(s.T(), typesv1.DefaultSmartAccountId+1, newAccID)
}

func (s *KeeperTestSuite) TestGetSetDelSignerAddress() {
	testAddress := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	keeper := s.App.SaKeeper

	acc, err := sdk.AccAddressFromBech32(testAddress)
	require.NoError(s.T(), err)

	keeper.SetSignerAddress(s.ctx, acc)
	getAcc := keeper.GetSignerAddress(s.ctx)
	require.Equal(s.T(), acc, getAcc)

	keeper.DeleteSignerAddress(s.ctx)
	getAcc = keeper.GetSignerAddress(s.ctx)
	require.Equal(s.T(), sdk.AccAddress(nil), getAcc)
}

func (s *KeeperTestSuite) TestGasRemaining() {

	keeper := s.App.SaKeeper

	gasRemaining := uint64(100)

	require.Equal(s.T(), keeper.HasGasRemaining(s.ctx), false)

	keeper.SetGasRemaining(s.ctx, gasRemaining)
	require.Equal(s.T(), keeper.HasGasRemaining(s.ctx), true)

	gas := keeper.GetGasRemaining(s.ctx)
	require.Equal(s.T(), gasRemaining, gas)

	keeper.DeleteGasRemaining(s.ctx)
	require.Equal(s.T(), keeper.HasGasRemaining(s.ctx), false)
}

func (s *KeeperTestSuite) TestValidateActivateSA() {

	keeper := s.App.SaKeeper

	pubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultPubKey)
	require.NoError(s.T(), err)

	// add new base account to chain
	err = helper.AddNewBaseAccount(s.App, s.ctx, helper.UserAddr, nil, uint64(0))
	require.NoError(s.T(), err)

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

		_, err := keeper.ValidateActiveSA(s.ctx, msg)

		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}

}

func (s *KeeperTestSuite) TestPrepareBeforeActive() {

	keeper := s.App.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(s.App, s.ctx, helper.UserAddr, nil, 1)
	require.NoError(s.T(), err)

	accAddr, err := sdk.AccAddressFromBech32(helper.UserAddr)
	require.NoError(s.T(), err)
	acc := s.App.AccountKeeper.GetAccount(s.ctx, accAddr)
	require.Equal(s.T(), uint64(1), acc.GetSequence())

	// PrepareBeforeActive et sequence of account to zero
	err = keeper.PrepareBeforeActive(s.ctx, acc)
	require.NoError(s.T(), err)

	acc = s.App.AccountKeeper.GetAccount(s.ctx, accAddr)
	require.Equal(s.T(), uint64(0), acc.GetSequence())
}

func (s *KeeperTestSuite) TestActiveSmartAccount() {
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
		cachedCtx, _ := s.ctx.CacheContext()

		keeper := s.App.SaKeeper

		acc, pub, err := helper.GenerateInActivateAccount(
			s.App,
			cachedCtx,
			helper.DefaultPubKey,
			uint64(1),
			helper.DefaultSalt,
			helper.DefaultMsg,
		)
		require.NoError(s.T(), err)

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

		_, err = keeper.ActiveSmartAccount(cachedCtx, msg, acc)
		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}
}

func (s *KeeperTestSuite) TestHandleAfterActive() {

	keeper := s.App.SaKeeper

	// add new base account to chain
	err := helper.AddNewBaseAccount(s.App, s.ctx, helper.UserAddr, nil, 1)
	require.NoError(s.T(), err)

	accAddr, err := sdk.AccAddressFromBech32(helper.UserAddr)
	require.NoError(s.T(), err)
	acc := s.App.AccountKeeper.GetAccount(s.ctx, accAddr)
	require.Equal(s.T(), uint64(1), acc.GetSequence())
	require.Equal(s.T(), nil, acc.GetPubKey())

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultPubKey)
	require.NoError(s.T(), err)

	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(s.T(), err)

	// PrepareBeforeActive et sequence of account to zero
	backupSe := uint64(2)
	err = keeper.HandleAfterActive(s.ctx, acc, backupSe, dPubKey)
	require.NoError(s.T(), err)

	acc = s.App.AccountKeeper.GetAccount(s.ctx, accAddr)
	require.Equal(s.T(), uint64(2), acc.GetSequence())
	require.Equal(s.T(), dPubKey, acc.GetPubKey())
}

func (s *KeeperTestSuite) TestValidateRecoverSA() {
	defaultCredentials := ""
	testSAAddress := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"
	testBAAddress := "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3"

	keeper := s.App.SaKeeper

	err := helper.AddNewSmartAccount(s.App, s.ctx, testSAAddress, nil, 0)
	require.NoError(s.T(), err)

	err = helper.AddNewBaseAccount(s.App, s.ctx, testBAAddress, nil, 0)
	require.NoError(s.T(), err)

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
		pubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultPubKey)
		require.NoError(s.T(), err)

		msg := &typesv1.MsgRecover{
			Creator:     helper.UserAddr,
			Address:     tc.AccountAddress,
			PubKey:      pubKey,
			Credentials: defaultCredentials,
		}

		_, err = keeper.ValidateRecoverSA(s.ctx, msg)
		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}
}

func (s *KeeperTestSuite) TestCallSMValidate() {
	customMsg := []byte("{\"recover_key\":\"024ab33b4f0808eba493ac4e3ead798c8339e2fd216b20ca110001fd094784c07f\"}")

	keeper := s.App.SaKeeper

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultPubKey)
	require.NoError(s.T(), err)

	rPubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultRPubKery)
	require.NoError(s.T(), err)
	dRPubKey, err := typesv1.PubKeyDecode(rPubKey)
	require.NoError(s.T(), err)

	acc, _, err := helper.GenerateInActivateAccount(
		s.App,
		s.ctx,
		helper.DefaultPubKey,
		uint64(2),
		helper.DefaultSalt,
		customMsg,
	)
	require.NoError(s.T(), err)

	msg := &typesv1.MsgActivateAccount{
		AccountAddress: acc.GetAddress().String(),
		CodeID:         uint64(2),
		Salt:           helper.DefaultSalt,
		InitMsg:        customMsg,
		PubKey:         pubKey,
	}
	_, err = keeper.ActiveSmartAccount(s.ctx, msg, acc)
	require.NoError(s.T(), err)

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

		err = keeper.CallSMValidate(s.ctx, msg, acc.GetAddress(), dRPubKey)
		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}
}

func (s *KeeperTestSuite) TestIsInactiveAccount() {
	testBAAddress1 := "cosmos1kzlrmxw3h2n4uzuv73m33cfw7xt7qjf3hlqx33ulc02e9dhxu46qgfxg9l"
	testBAAddress2 := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	keeper := s.App.SaKeeper

	// prepare pubkey
	pubKey, err := typesv1.PubKeyToAny(s.App.AppCodec(), helper.DefaultPubKey)
	require.NoError(s.T(), err)

	dPubKey, err := typesv1.PubKeyDecode(pubKey)
	require.NoError(s.T(), err)

	err = helper.AddNewBaseAccount(s.App, s.ctx, testBAAddress1, nil, 0)
	require.NoError(s.T(), err)

	err = helper.AddNewBaseAccount(s.App, s.ctx, testBAAddress2, dPubKey, 0)
	require.NoError(s.T(), err)

	acc, _, err := helper.GenerateInActivateAccount(
		s.App,
		s.ctx,
		helper.DefaultPubKey,
		helper.DefaultCodeID,
		helper.DefaultSalt,
		helper.DefaultMsg,
	)
	require.NoError(s.T(), err)

	msg := &typesv1.MsgActivateAccount{
		AccountAddress: acc.GetAddress().String(),
		CodeID:         helper.DefaultCodeID,
		Salt:           helper.DefaultSalt,
		InitMsg:        helper.DefaultMsg,
		PubKey:         pubKey,
	}
	_, err = keeper.ActiveSmartAccount(s.ctx, msg, acc)
	require.NoError(s.T(), err)

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
		require.NoError(s.T(), err)

		_, err = keeper.IsInactiveAccount(s.ctx, acc)
		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}
	}
}

func (s *KeeperTestSuite) TestGetSmartAccountByAddress() {
	testAddress1 := "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3"
	testAddress2 := "cosmos1kzlrmxw3h2n4uzuv73m33cfw7xt7qjf3hlqx33ulc02e9dhxu46qgfxg9l"
	testAddress3 := "cosmos10uxaa5gkxpeungu2c9qswx035v6t3r24w6v2r6dxd858rq2mzknqj8ru28"

	err := helper.AddNewSmartAccount(s.App, s.ctx, testAddress1, nil, 0)
	require.NoError(s.T(), err)

	err = helper.AddNewBaseAccount(s.App, s.ctx, testAddress2, nil, 0)
	require.NoError(s.T(), err)

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
		require.NoError(s.T(), err)

		saAcc, err := s.App.SaKeeper.GetSmartAccountByAddress(s.ctx, acc)

		if tc.err {
			require.Error(s.T(), err)
		} else {
			require.NoError(s.T(), err)
		}

		if tc.exist {
			require.NotEqual(s.T(), (*typesv1.SmartAccount)(nil), saAcc)
		} else {
			require.Equal(s.T(), (*typesv1.SmartAccount)(nil), saAcc)
		}
	}
}

func (s *KeeperTestSuite) TestCheckAllowedMsgs() {

	keeper := s.App.SaKeeper

	params := typesv1.Params{
		DisableMsgsList: []string{
			"/cosmwasm.wasm.v1.MsgUpdateAdmin",
			"/cosmwasm.wasm.v1.MsgClearAdmin",
		},
		MaxGasExecute: 2000000,
	}
	err := keeper.SetParams(s.ctx, params)
	require.NoError(s.T(), err)

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
		err := keeper.CheckAllowedMsgs(s.ctx, tc.msgs)

		if !tc.err {
			require.NoError(s.T(), err)
		} else {
			require.Error(s.T(), err)
		}
	}
}
