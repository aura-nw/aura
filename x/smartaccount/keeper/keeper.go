package keeper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       sdk.StoreKey
		memKey         sdk.StoreKey
		paramstore     paramtypes.Subspace
		WasmKeeper     wasmkeeper.Keeper
		ContractKeeper *wasmkeeper.PermissionedKeeper
		AccountKeeper  types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	wp wasmkeeper.Keeper,
	contractKeeper *wasmkeeper.PermissionedKeeper,
	ak types.AccountKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		WasmKeeper:     wp,
		ContractKeeper: contractKeeper,
		AccountKeeper:  ak,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ------------------------------- NextAccountId -------------------------------

func (k Keeper) GetAndIncrementNextAccountID(ctx sdk.Context) uint64 {
	id := k.GetNextAccountID(ctx)

	k.SetNextAccountID(ctx, id+1)

	return id
}

func (k Keeper) GetNextAccountID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	return sdk.BigEndianToUint64(store.Get(types.KeyPrefix(types.AccountIDKey)))
}

func (k Keeper) SetNextAccountID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.AccountIDKey), sdk.Uint64ToBigEndian(id))
}

func (k Keeper) ValidateActiveSA(ctx sdk.Context, msg *types.MsgActivateAccount) (authtypes.AccountI, error) {
	// validate code id use to init smart account
	if !k.isWhitelistCodeID(ctx, msg.CodeID) {
		k.Logger(ctx).Error("active-sm", "code-id", msg.CodeID)
		return nil, types.ErrInvalidCodeID
	}

	// get smart contract account by address, account must exist in chain before activation
	signer, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		k.Logger(ctx).Error("active-sm", "account-addr", msg.AccountAddress, "err", err.Error())
		return nil, types.ErrInvalidAddress
	}

	return k.IsInactiveAccount(ctx, signer)
}

func (k Keeper) PrepareBeforeActive(ctx sdk.Context, sAccount authtypes.AccountI) error {
	// set sequence of smart account to ZERO
	// we need to set to zero for keep balance and information of sm
	err := sAccount.SetSequence(0)
	if err != nil {
		return err
	}

	k.AccountKeeper.SetAccount(ctx, sAccount)

	return nil
}

func (k Keeper) ActiveSmartAccount(
	ctx sdk.Context,
	msg *types.MsgActivateAccount,
	sAccount authtypes.AccountI,
) (cryptotypes.PubKey, error) {

	pubKey, err := types.PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, err
	}

	// we use pubkey.Address() as owner of this contract
	// remember this account doesn't exist on chain yet if have not received any funds before
	owner := sdk.AccAddress(pubKey.Address())

	// instantiate smartcontract by code_id
	address, _, err := k.ContractKeeper.Instantiate2(
		ctx,
		msg.CodeID,
		owner,                 // owner
		sAccount.GetAddress(), // admin
		msg.InitMsg,           // message
		fmt.Sprintf("%s/%d", types.ModuleName, k.GetAndIncrementNextAccountID(ctx)), // label
		sdk.NewCoins(), // empty funds
		msg.Salt,       // salt
		true,
	)
	if err != nil {
		return nil, err
	}

	contractAddrStr := address.String()

	// make sure the new contract has the same address as predicted
	if contractAddrStr != msg.AccountAddress {
		k.Logger(ctx).Error("active-sm", "contract-addr", contractAddrStr, "input-addr", msg.AccountAddress)
		return nil, types.ErrBadInstantiateMsg
	}

	ctx.Logger().Info(
		"smart account created",
		types.AttributeKeyCreator, msg.AccountAddress,
		types.AttributeKeyCodeID, msg.CodeID,
		types.AttributeKeyContractAddr, contractAddrStr,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSmartAccountActivated,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.AccountAddress),
			sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(msg.CodeID, 10)),
			sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddrStr),
		),
	)

	return pubKey, nil
}

// HandleAfterActive change type of account to smart account and recover seq of this account
func (k Keeper) HandleAfterActive(ctx sdk.Context, sAccount authtypes.AccountI, backupSeq uint64, pubKey cryptotypes.PubKey) error {
	// set sequence of new account to pre-sequence
	err := sAccount.SetSequence(backupSeq)
	if err != nil {
		return err
	}

	// create new smart account type
	smartAccount := types.NewSmartAccountFromAccount(sAccount)

	// set smart account pubkey
	return k.UpdateAccountPubKey(ctx, smartAccount, pubKey)
}

// ValidateRecoverSA check input before recover smart account
func (k Keeper) ValidateRecoverSA(ctx sdk.Context, msg *types.MsgRecover) (authtypes.AccountI, error) {
	saAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		k.Logger(ctx).Error("recover-sa", "decode-err", err.Error())
		return nil, types.ErrInvalidAddress
	}

	// only allow accounts with type SmartAccount to be restored pubkey
	smartAccount := k.AccountKeeper.GetAccount(ctx, saAddr)
	if _, ok := smartAccount.(*types.SmartAccount); !ok {
		return nil, types.ErrAccountNotFoundForAddress
	}

	return smartAccount, nil
}

// CallSMValidate to check logic recover from smart account
func (k Keeper) CallSMValidate(ctx sdk.Context, msg *types.MsgRecover, saAddr sdk.AccAddress, pubKey cryptotypes.PubKey) error {
	// credentials
	credentials, err := base64.StdEncoding.DecodeString(msg.Credentials)
	if err != nil {
		k.Logger(ctx).Error("recover-sa-decodestr", "err", err.Error())
		return types.ErrInvalidCredentials
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
		return err
	}

	// check recover logic in smart contract
	_, err = k.ContractKeeper.Sudo(ctx, saAddr, sudoMsgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) UpdateAccountPubKey(ctx sdk.Context, acc authtypes.AccountI, pubKey cryptotypes.PubKey) error {
	err := acc.SetPubKey(pubKey)
	if err != nil {
		return err
	}

	k.AccountKeeper.SetAccount(ctx, acc)

	return nil
}

// isWhitelistCodeID return true if codeID in the whitelist
// otherwise return false
func (k Keeper) isWhitelistCodeID(ctx sdk.Context, codeID uint64) bool {
	params := k.GetParams(ctx)

	return params.IsAllowedCodeID(codeID)
}

// Inactive smart-account must be base account with empty public key or smart account
// and has not been used for any instantiated contracts
func (k Keeper) IsInactiveAccount(ctx sdk.Context, acc sdk.AccAddress) (authtypes.AccountI, error) {
	sAccount := k.AccountKeeper.GetAccount(ctx, acc)

	// check if account has type base or smart
	_, isBaseAccount := sAccount.(*authtypes.BaseAccount)
	_, isSmartAccount := sAccount.(*types.SmartAccount)
	if !isBaseAccount && !isSmartAccount {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFoundForAddress, acc.String())
	}

	// check if base account already has public key
	if sAccount.GetPubKey() != nil && isBaseAccount {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExists, acc.String())
	}

	// check if contract with account not been instantiated
	if k.WasmKeeper.HasContractInfo(ctx, acc) {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExists, acc.String())
	}

	return sAccount, nil
}

func (k Keeper) GetSmartAccountByAddress(ctx sdk.Context, address sdk.AccAddress) (*types.SmartAccount, error) {
	signerAcc, err := authante.GetSignerAcc(ctx, k.AccountKeeper, address)
	if err != nil {
		return nil, err
	}

	saAcc, ok := signerAcc.(*types.SmartAccount)
	if !ok {
		return nil, nil
	}

	return saAcc, nil
}
