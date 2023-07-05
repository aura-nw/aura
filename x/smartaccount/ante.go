package smartaccount

import (
	"bytes"
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sakeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GenerateValidateQueryMessage(msg *wasmtypes.MsgExecuteContract, msgs []types.MsgData) ([]byte, error) {
	var accMsg types.AccountMsg
	umErr := json.Unmarshal(msg.GetMsg(), &accMsg)
	if umErr != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, umErr.Error())
	} else if accMsg.AfterExecuteTx == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "must be AfterExecute message")
	}

	callMsgData, err := json.Marshal(accMsg.AfterExecuteTx.Msgs)
	if err != nil {
		return nil, err
	}

	msgData, err := json.Marshal(&msgs)
	if err != nil {
		return nil, err
	}

	// data in validate message must compatiable to tx.messages
	if !bytes.Equal(msgData, callMsgData) {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "after-execute message data not compatible with tx.messages")
	}

	validateMessage, err := json.Marshal(&types.AccountMsg{
		ValidateTx: &types.ValidateTx{
			Msgs: msgs,
		},
	})
	if err != nil {
		return nil, err
	}

	return validateMessage, nil
}

func IsSmartAccountTx(ctx sdk.Context, tx sdk.Tx, accountKeeper authante.AccountKeeper) (bool, *types.SmartAccount, *txsigning.SignatureV2, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return false, nil, nil, sdkerrors.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return false, nil, nil, err
	}

	signerAddrs := sigTx.GetSigners()

	// do not allow multi signer
	if len(signerAddrs) != 1 || len(sigs) != 1 {
		return false, nil, nil, nil
	}

	signerAcc, err := authante.GetSignerAcc(ctx, accountKeeper, signerAddrs[0])
	if err != nil {
		return false, nil, nil, err
	}

	saAcc, ok := signerAcc.(*types.SmartAccount)
	if !ok {
		return false, nil, nil, nil
	}

	return true, saAcc, &sigs[0], nil
}

func IsActivateAccountMessage(tx sdk.Tx) (bool, *types.MsgActivateAccount, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return false, nil, sdkerrors.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	msgs := sigTx.GetMsgs()

	if len(msgs) != 1 {
		// smart account activation message must stand alone
		for _, msg := range msgs {
			if _, ok := msg.(*types.MsgActivateAccount); ok {
				return false, nil, sdkerrors.Wrap(types.ErrInvalidTx, "smart account activation message must stand alone")
			}
		}

		return false, nil, nil
	}

	activateMsg, ok := msgs[0].(*types.MsgActivateAccount)
	if !ok {
		return false, nil, nil
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return false, nil, err
	}

	signer := sigTx.GetSigners()

	// do not allow multi signer and signature
	if len(signer) != 1 || len(sigs) != 1 {
		return false, nil, sdkerrors.Wrap(types.ErrInvalidTx, "smart-account activation tx does not allow multiple signers")
	}

	return true, activateMsg, nil
}

// ------------------------- SmartAccount Decorator ------------------------- \\

type SmartAccountDecorator struct {
	SaKeeper      sakeeper.Keeper
	WasmKeeper    wasmkeeper.Keeper
	AccountKeeper authante.AccountKeeper
}

func NewSmartAccountDecorator(wasmKeeper wasmkeeper.Keeper, accountKeeper authante.AccountKeeper, saKeeper sakeeper.Keeper) *SmartAccountDecorator {
	return &SmartAccountDecorator{
		SaKeeper:      saKeeper,
		WasmKeeper:    wasmKeeper,
		AccountKeeper: accountKeeper,
	}
}

// AnteHandle is used for performing basic validity checks on a transaction such that it can be thrown out of the mempool.
func (decorator *SmartAccountDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	isActivateAccount, activateMsg, err := IsActivateAccountMessage(tx)
	if err != nil {
		return ctx, err
	}

	// if is not activate account message, next to SmartAccountTxDecorator
	if !isActivateAccount {
		satd := NewSmartAccountTxDecorator(decorator.WasmKeeper, decorator.AccountKeeper, decorator.SaKeeper)
		return satd.AnteHandle(ctx, tx, simulate, next)
	}

	// in ReCheckTx mode, below check may not be necessary

	// get signer of smart account activation message
	signer := activateMsg.GetSigners()[0]

	// decode string to pubkey
	pubKey, err := types.PubKeyDecode(activateMsg.PubKey)
	if err != nil {
		return ctx, err
	}

	// generate predictable address using Instantiate2's PredicableAddressGenerator
	predicAddr, err := types.Instantiate2Address(
		ctx,
		decorator.WasmKeeper,
		activateMsg.CodeID,
		activateMsg.InitMsg,
		activateMsg.Salt,
		pubKey,
	)
	if err != nil {
		return ctx, err
	}

	// the signer of the activation message must be the address generated by Instantiate2's PredicableAddressGenerator
	if !signer.Equals(predicAddr) {
		return ctx, sdkerrors.Wrap(types.ErrInvalidAddress, "not the same as predicted")
	}

	// if in delivery mode, remove temporary pubkey from account
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() && !simulate {
		// get smart contract account by address
		sAccount := decorator.AccountKeeper.GetAccount(ctx, signer)
		if _, ok := sAccount.(*authtypes.BaseAccount); !ok {
			return ctx, sdkerrors.Wrap(types.ErrAccountNotFoundForAddress, signer.String())
		}

		// remove temporary pubkey for account
		err = decorator.SaKeeper.UpdateAccountPubKey(ctx, sAccount, nil)
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// ------------------------- SmartAccountTx Decorator ------------------------- \\

type SmartAccountTxDecorator struct {
	SaKeeper      sakeeper.Keeper
	WasmKeeper    wasmkeeper.Keeper
	AccountKeeper authante.AccountKeeper
}

func NewSmartAccountTxDecorator(wasmKeeper wasmkeeper.Keeper, accountKeeper authante.AccountKeeper, saKeeper sakeeper.Keeper) *SmartAccountTxDecorator {
	return &SmartAccountTxDecorator{
		SaKeeper:      saKeeper,
		WasmKeeper:    wasmKeeper,
		AccountKeeper: accountKeeper,
	}
}

// AnteHandle is used for performing basic validity checks on a transaction such that it can be thrown out of the mempool.
func (decorator *SmartAccountTxDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	isSmartAccountTx, signerAcc, _, err := IsSmartAccountTx(ctx, tx, decorator.AccountKeeper)
	if err != nil {
		return ctx, err
	}

	// if is not smartaccount tx type
	if !isSmartAccountTx {
		// do some thing
		return next(ctx, tx, simulate)
	}

	// not support smartaccount tx simulation yet
	if simulate {
		return ctx, sdkerrors.Wrap(types.ErrNotSupported, "Simulation of SmartAccount txs isn't supported yet")
	}

	msgs := tx.GetMsgs()

	// validate message must be the last message and must be MsgExecuteContract
	var valMsg *wasmtypes.MsgExecuteContract
	if msg, err := msgs[len(msgs)-1].(*wasmtypes.MsgExecuteContract); err {
		valMsg = msg
	} else {
		return ctx, sdkerrors.Wrap(types.ErrInvalidMsg, "validate message must be type MsgExecuteContract")
	}

	// get smartaccount address
	saAddress := signerAcc.GetAddress().String()

	// the message must be sent from the signer's address which is also the smart contract address
	if valMsg.Sender != saAddress || valMsg.Contract != saAddress {
		return ctx, sdkerrors.Wrap(types.ErrInvalidMsg, "sender address and contract address must be the same")
	}

	// parse messages in tx to list of string
	valMsgData, err := types.ParseMessagesString(msgs[:len(msgs)-1])
	if err != nil {
		return ctx, err
	}

	// create message for SA contract query
	validateMessage, err := GenerateValidateQueryMessage(valMsg, valMsgData)
	if err != nil {
		return ctx, err
	}

	params := decorator.SaKeeper.GetParams(ctx)

	// query SA contract for validating transaction with limit gas
	err = querySmartWithGasLimit(ctx, decorator.WasmKeeper, signerAcc.GetAddress(), validateMessage, params.MaxGasQuery)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// Call a contract's query smart with a gas limit
// referenced from Osmosis' protorev posthandler:
// https://github.com/osmosis-labs/osmosis/blob/98025f185ab2ee1b060511ed22679112abcc08fa/x/protorev/keeper/posthandler.go#L42-L43
func querySmartWithGasLimit(
	ctx sdk.Context, wasmKeeper wasmkeeper.Keeper,
	contractAddr sdk.AccAddress, msg []byte, maxGas sdk.Gas,
) error {
	cacheCtx, write := ctx.CacheContext()
	cacheCtx = cacheCtx.WithGasMeter(sdk.NewGasMeter(maxGas))

	if _, err := wasmKeeper.QuerySmart(cacheCtx, contractAddr, msg); err != nil {
		return err
	}

	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return nil
}

// ------------------------- SetPubKey Decorator ------------------------- \\

type SetPubKeyDecorator struct {
	AccountKeeper authante.AccountKeeper
	WasmKeeper    wasmkeeper.Keeper
	saKeeper      sakeeper.Keeper
}

func NewSetPubKeyDecorator(accountKeeper authante.AccountKeeper, wasmKeeper wasmkeeper.Keeper, saKeeper sakeeper.Keeper) *SetPubKeyDecorator {
	return &SetPubKeyDecorator{
		AccountKeeper: accountKeeper,
		WasmKeeper:    wasmKeeper,
		saKeeper:      saKeeper,
	}
}

func (decorator *SetPubKeyDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	isActivateAccountMsg, activateMsg, err := IsActivateAccountMessage(tx)
	if err != nil {
		return ctx, err
	}

	// if is smart account activation message
	if isActivateAccountMsg {
		// get message signer
		signer := activateMsg.GetSigners()[0]

		// get smart contract account by address, account must be inactivate smart account
		sAccount, err := decorator.saKeeper.IsInactiveAccount(ctx, signer)
		if err != nil {
			return ctx, err
		}

		// decode any to pubkey
		pubKey, err := types.PubKeyDecode(activateMsg.PubKey)
		if err != nil {
			return ctx, err
		}

		// set temporary pubkey for account
		// need this for the next ante signature checks
		err = decorator.saKeeper.UpdateAccountPubKey(ctx, sAccount, pubKey)
		if err != nil {
			return ctx, err
		}

		return next(ctx, tx, simulate)
	}

	isSmartAccountTx, signerAcc, sig, err := IsSmartAccountTx(ctx, tx, decorator.AccountKeeper)
	if err != nil {
		return ctx, err
	}

	// if is smart account tx skip authante NewSetPubKeyDecorator
	// need this to avoid pubkey and address equal check of above decorator
	if isSmartAccountTx {
		// if this is smart account tx, check if pubkey is set
		if signerAcc.GetPubKey() == nil {
			return ctx, sdkerrors.Wrap(types.ErrNilPubkey, signerAcc.String())
		}

		// Also emit the following events, so that txs can be indexed by these
		var events sdk.Events
		events = append(events, sdk.NewEvent(types.EventTypeSmartAccountTx,
			sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signerAcc, sig.Sequence)),
		))

		// maybe need add more event here

		ctx.EventManager().EmitEvents(events)

		return next(ctx, tx, simulate)
	}

	// default authant SetPubKeyDecorator
	svd := authante.NewSetPubKeyDecorator(decorator.AccountKeeper)
	return svd.AnteHandle(ctx, tx, simulate, next)
}
