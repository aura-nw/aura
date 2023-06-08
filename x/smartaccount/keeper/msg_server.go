package keeper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type msgServer struct {
	Keeper         Keeper
	ContractKeeper *wasmkeeper.PermissionedKeeper
	AccountKeeper  types.AccountKeeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractKeeper *wasmkeeper.PermissionedKeeper, accountKeeper types.AccountKeeper) types.MsgServer {
	return &msgServer{
		Keeper:         keeper,
		ContractKeeper: contractKeeper,
		AccountKeeper:  accountKeeper,
	}
}

func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	saAddress, data, err := InstantiateSmartAccount(ctx, k.Keeper, k.ContractKeeper, msg)
	if err != nil {
		return nil, err
	}

	saAddressStr := saAddress.String()

	// get smart contract account by address
	scAccount := k.AccountKeeper.GetAccount(ctx, saAddress)
	if _, ok := scAccount.(*authtypes.BaseAccount); !ok {
		return nil, fmt.Errorf(types.ErrAccountNotFoundForAddress, saAddressStr)
	}

	// secp25k61 public key
	pubKey, err := PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, err
	}

	// create new smartaccount type
	smartAccount := types.NewSmartAccountFromAccount(scAccount)

	// set smartaccount pubkey
	err = smartAccount.SetPubKey(pubKey)
	if err != nil {
		return nil, err
	}

	// update smartaccount
	k.AccountKeeper.SetAccount(ctx, smartAccount)

	return &types.MsgCreateAccountResponse{
		Address: saAddressStr,
		Data:    data,
	}, nil
}

func (k msgServer) Recover(goCtx context.Context, msg *types.MsgRecover) (*types.MsgRecoverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	saAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, fmt.Errorf(types.ErrAddressFromBech32, err.Error())
	}

	smartAccount := k.AccountKeeper.GetAccount(ctx, saAddr)
	if _, ok := smartAccount.(*types.SmartAccount); !ok {
		return nil, fmt.Errorf(types.ErrAccountNotFoundForAddress, msg.Address)
	}

	// secp25k61 public key
	pubKey, err := PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, err
	}

	// credentials
	credentials, err := base64.StdEncoding.DecodeString(msg.Credentials)
	if err != nil {
		return nil, err
	}

	sudoMsgBytes, err := json.Marshal(&types.AccountMsg{
		RecoverTx: &types.RecoverTx{
			PubKey:      pubKey.GetKey(),
			Credentials: credentials,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = k.ContractKeeper.Sudo(ctx, saAddr, sudoMsgBytes)
	if err != nil {
		return nil, err
	}

	// set new pubkey for smartaccount
	err = smartAccount.SetPubKey(pubKey)
	if err != nil {
		return nil, err
	}

	// update smartaccount
	k.AccountKeeper.SetAccount(ctx, smartAccount)

	return &types.MsgRecoverResponse{
		Address:   msg.Address,
		NewPubKey: msg.PubKey,
	}, nil
}
