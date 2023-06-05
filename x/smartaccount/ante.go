package smartaccount

import (
	"bytes"
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	smartaccountkeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// ------------------------- SmartAccount Decorator ------------------------- \\

type SmartAccountDecorator struct {
	SmartAccountKeeper smartaccountkeeper.Keeper
	WasmKeeper         wasmkeeper.Keeper
	AccountKeeper      authante.AccountKeeper
}

func NewSmartAccountDecorator(smartAccountKeeper smartaccountkeeper.Keeper, wasmKeeper wasmkeeper.Keeper, accountKeeper authante.AccountKeeper) *SmartAccountDecorator {
	return &SmartAccountDecorator{
		SmartAccountKeeper: smartAccountKeeper,
		WasmKeeper:         wasmKeeper,
		AccountKeeper:      accountKeeper,
	}
}

func GenerateValidateQueryMessage(msg *wasmtypes.MsgExecuteContract, msgs []types.MsgData) ([]byte, error) {
	var accMsg types.AccountMsg
	umErr := json.Unmarshal(msg.GetMsg(), &accMsg)
	if umErr != nil {
		return nil, fmt.Errorf("invalid smart account message: %s", umErr.Error())
	} else if accMsg.AfterExecuteTx == nil {
		return nil, fmt.Errorf("must be AfterExecute message")
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
		return nil, fmt.Errorf("invalid after-execute message data: not compatible with tx.messages")
	}

	validateMessage, err := json.Marshal(&types.AccountMsg{
		ValidateTx: &types.ValidateTx{
			Msgs: msgs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot json marshal validate message: %s", err.Error())
	}

	return validateMessage, nil
}

func IsSmartAccountTx(ctx sdk.Context, tx sdk.Tx, accountKeeper authante.AccountKeeper) (bool, *types.SmartAccount, *txsigning.SignatureV2, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return false, nil, nil, fmt.Errorf(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return false, nil, nil, err
	}

	signerAddrs := sigTx.GetSigners()

	// do not allow multi signer yet
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

// AnteHandle is used for performing basic validity checks on a transaction such that it can be thrown out of the mempool.
func (decorator *SmartAccountDecorator) AnteHandle(
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

	msgs := tx.GetMsgs()

	if len(msgs) < 1 {
		return ctx, fmt.Errorf("tx must contain at least the validate message")
	}

	// validate message must be the last message and must be MsgExecuteContract
	var valMsg *wasmtypes.MsgExecuteContract
	if msg, err := msgs[len(msgs)-1].(*wasmtypes.MsgExecuteContract); err {
		valMsg = msg
	} else {
		return ctx, fmt.Errorf("validate message must be type MsgExecuteContract")
	}

	// get smartaccount address
	saAddress := signerAcc.GetAddress().String()

	// the message must be sent from the signer's address which is also the smart contract address
	if valMsg.Sender != saAddress || valMsg.Contract != saAddress {
		return ctx, fmt.Errorf(
			"invalid validate message: sender address and smart contract must be the same",
		)
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

	// query SA contract for validating transaction
	_, vErr := decorator.WasmKeeper.QuerySmart(ctx, signerAcc.GetAddress(), validateMessage)
	if vErr != nil {
		return ctx, fmt.Errorf("tx validate fail: %s", vErr.Error())
	}

	return next(ctx, tx, simulate)
}

// ------------------------- SetPubKey Decorator ------------------------- \\

type SetPubKeyDecorator struct {
	SmartAccountKeeper smartaccountkeeper.Keeper
	AccountKeeper      authante.AccountKeeper
}

func NewSetPubKeyDecorator(smartAccountKeeper smartaccountkeeper.Keeper, accountKeeper authante.AccountKeeper) *SetPubKeyDecorator {
	return &SetPubKeyDecorator{
		SmartAccountKeeper: smartAccountKeeper,
		AccountKeeper:      accountKeeper,
	}
}

func (decorator *SetPubKeyDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	isSmartAccountTx, signerAcc, sig, err := IsSmartAccountTx(ctx, tx, decorator.AccountKeeper)
	if err != nil {
		return ctx, err
	}

	if !isSmartAccountTx {
		svd := authante.NewSetPubKeyDecorator(decorator.AccountKeeper)
		return svd.AnteHandle(ctx, tx, simulate, next)
	}

	// if this is smart account tx, check if pubkey is set
	if signerAcc.GetPubKey() == nil {
		return ctx, fmt.Errorf("smart-account PublicKey must not be null")
	}

	// Also emit the following events, so that txs can be indexed by these
	var events sdk.Events
	events = append(events, sdk.NewEvent(sdk.EventTypeTx,
		sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signerAcc, sig.Sequence)),
	))

	// maybe need add more event here

	ctx.EventManager().EmitEvents(events)

	return next(ctx, tx, simulate)
}
