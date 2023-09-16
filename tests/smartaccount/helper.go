package smartaccount

import (
	"fmt"
	"os"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	"github.com/aura-nw/aura/app"
	"github.com/aura-nw/aura/tests"
	"github.com/aura-nw/aura/x/smartaccount"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	UserAddr     = "cosmos1lg0g3jpu8luawwezcknamz0l003swknjyw9uch"
	GenesisState = &typesv1.GenesisState{
		Params:         typesv1.NewParams([]*typesv1.CodeID{{CodeID: 1, Status: true}}, []string{"/cosmwasm.wasm.v1.MsgExecuteContract"}, typesv1.DefaultMaxGas),
		SmartAccountId: typesv1.DefaultSmartAccountId,
	}
)

const (
	WasmPath1 = "../../tests/smartaccount/wasm/"
	WasmPath2 = "../../../tests/smartaccount/wasm/"
)

var (
	DefaultSalt     = []byte("test")
	DefaultMsg      = []byte("{}")
	DefaultCodeID   = uint64(1)
	DefaultPubKey   = []byte("{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"AnZfdXVALfIcNjpqgzH/4nWsSpP7l5PiCyZAuAWQRBUz\"}")
	DefaultRPubKery = []byte("{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A/2t0ru/iZ4HoiX0DkTuMy9rC2mMeXmiN6luM3pa+IvT\"}")
)

func SetupGenesisTest() (sdk.Context, *app.App) {
	app := tests.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	smartaccount.InitGenesis(ctx, app.SaKeeper, *GenesisState)

	return ctx, app
}

func StoreCodeID(app *app.App, ctx sdk.Context, creator sdk.AccAddress, path string) (uint64, []byte, error) {
	wasm, err := os.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return 0, nil, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return 0, nil, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	return app.ContractKeeper.Create(ctx, creator, wasm, nil)
}

func GenerateInActivateAccount(
	app *app.App,
	ctx sdk.Context,
	path string,
	dPubKey []byte,
	dCodeID uint64,
	dSalt []byte,
	dMsg []byte,
) (*authtypes.BaseAccount, *codectypes.Any, error) {
	/* ======== store wasm ======== */
	user, err := sdk.AccAddressFromBech32(UserAddr)
	if err != nil {
		return nil, nil, err
	}

	// store code
	codeID, _, err := StoreCodeID(app, ctx, user, path)
	if err != nil {
		return nil, nil, err
	}
	if codeID != dCodeID {
		return nil, nil, fmt.Errorf("invalid codeID")
	}

	queryServer := app.SaKeeper

	/* ======== create inactivate smart account ======== */
	pubKey, err := typesv1.PubKeyToAny(app.AppCodec(), dPubKey)
	if err != nil {
		return nil, nil, err
	}

	queryMsg := &typesv1.QueryGenerateAccountRequest{
		CodeID:  dCodeID,
		PubKey:  pubKey,
		Salt:    dSalt,
		InitMsg: dMsg,
	}

	queryRes, err := queryServer.GenerateAccount(sdk.WrapSDKContext(ctx), queryMsg)
	if err != nil {
		return nil, nil, err
	}

	newAccAddr, err := sdk.AccAddressFromBech32(queryRes.Address)
	if err != nil {
		return nil, nil, err
	}
	newAcc := authtypes.NewBaseAccount(newAccAddr, nil, app.AccountKeeper.NextAccountNumber(ctx), 0)

	app.AccountKeeper.SetAccount(ctx, newAcc)

	return newAcc, pubKey, nil
}

func AddNewBaseAccount(app *app.App, ctx sdk.Context, addr string, pubKey cryptotypes.PubKey, sequence uint64) error {
	newAcc, err := NewBaseAccount(app, ctx, addr, pubKey, sequence)
	if err != nil {
		return err
	}

	app.AccountKeeper.SetAccount(ctx, newAcc)
	return nil
}

func AddNewSmartAccount(app *app.App, ctx sdk.Context, addr string, pubKey cryptotypes.PubKey, sequence uint64) error {

	sdk.GetConfig().SetBech32PrefixForAccount("cosmos", "")

	newAcc, err := NewBaseAccount(app, ctx, addr, nil, sequence)
	if err != nil {
		return err
	}

	// create new smart account type
	smartAccount := typesv1.NewSmartAccountFromAccount(newAcc)

	err = smartAccount.SetPubKey(pubKey)
	if err != nil {
		return err
	}

	app.AccountKeeper.SetAccount(ctx, smartAccount)
	return nil
}

func NewBaseAccount(app *app.App, ctx sdk.Context, addr string, pubKey cryptotypes.PubKey, sequence uint64) (*authtypes.BaseAccount, error) {
	newAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	newAcc := authtypes.NewBaseAccount(newAccAddr, pubKey, app.AccountKeeper.NextAccountNumber(ctx), sequence)
	return newAcc, nil
}
