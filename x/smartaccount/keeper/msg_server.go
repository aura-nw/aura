package keeper

import (
	"context"
	"encoding/base64"
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) ActivateAccount(goCtx context.Context, msg *types.MsgActivateAccount) (*types.MsgActivateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if codeID is in whitelist
	if !IsWhitelistCodeID(ctx, k.Keeper, msg.CodeID) {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "codeID not allowed")
	}

	// get smart contract account by address, account must exist in chain before activation
	signer, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid smart account address (%s)", err)
	}

	// Only allow inactive smart account to be activated
	sAccount, err := types.IsInactiveAccount(ctx, signer, k.AccountKeeper, k.Keeper.wasmKeeper)
	if err != nil {
		return nil, err
	}

	// get current sequence of account
	sequence := sAccount.GetSequence()

	// set sequence to 0 so we can instantiate it later
	err = types.UpdateAccountSequence(ctx, k.AccountKeeper, sAccount, 0)
	if err != nil {
		return nil, err
	}

	saAddress, data, pubKey, err := InstantiateSmartAccount(ctx, k.Keeper, k.ContractKeeper, msg)
	if err != nil {
		return nil, err
	}

	saAddressStr := saAddress.String()

	// get smart contract account by address
	scAccount := k.AccountKeeper.GetAccount(ctx, saAddress)
	if _, ok := scAccount.(*authtypes.BaseAccount); !ok {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFoundForAddress, saAddressStr)
	}

	// set sequence of new account to pre-sequence
	err = scAccount.SetSequence(sequence)
	if err != nil {
		return nil, err
	}

	// create new smartaccount type
	smartAccount := types.NewSmartAccountFromAccount(scAccount)

	// set smartaccount pubkey
	err = types.UpdateAccountPubKey(ctx, k.AccountKeeper, smartAccount, pubKey)
	if err != nil {
		return nil, err
	}

	return &types.MsgActivateAccountResponse{
		Address: saAddressStr,
		Data:    data,
	}, nil
}

func (k msgServer) Recover(goCtx context.Context, msg *types.MsgRecover) (*types.MsgRecoverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	saAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid smart account address (%s)", err)
	}

	// only allow accounts with type SmartAccount to be restored pubkey
	smartAccount := k.AccountKeeper.GetAccount(ctx, saAddr)
	if _, ok := smartAccount.(*types.SmartAccount); !ok {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFoundForAddress, msg.Address)
	}

	// public key
	pubKey, err := types.PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, err
	}

	// credentials
	credentials, err := base64.StdEncoding.DecodeString(msg.Credentials)
	if err != nil {
		return nil, err
	}

	// data pass into recover message call
	sudoMsgBytes, err := json.Marshal(&types.AccountMsg{
		RecoverTx: &types.RecoverTx{
			Caller:      msg.Creator,    // caller of message
			PubKey:      pubKey.Bytes(), // new public key
			Credentials: credentials,    // credentials
		},
	})
	if err != nil {
		return nil, err
	}

	// check recover logic in smart acontract
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

	ctx.Logger().Info(
		"smart account recovery",
		types.AttributeKeyCreator, msg.Creator,
		types.AttributeKeyContractAddr, msg.Address,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAccountRecovery,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Address),
		),
	)

	return &types.MsgRecoverResponse{}, nil
}

// this line is used by starport scaffolding # handler/msgServer
