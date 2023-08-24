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
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetSmartAccountTxSigner(ctx sdk.Context, sigTx authsigning.SigVerifiableTx, saKeeper sakeeper.Keeper) (*types.SmartAccount, error) {
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return nil, err
	}

	signerAddrs := sigTx.GetSigners()

	// signer of smartaccount tx must stand alone
	if len(signerAddrs) != 1 || len(sigs) != 1 {
		return nil, nil
	}

	saAcc, err := saKeeper.GetSmartAccountByAddress(ctx, signerAddrs[0])
	if err != nil {
		return nil, err
	}

	return saAcc, nil
}

func GetValidActivateAccountMessage(sigTx authsigning.SigVerifiableTx) (*types.MsgActivateAccount, error) {
	msgs := sigTx.GetMsgs()

	if len(msgs) != 1 {
		// smart account activation message must stand alone
		for _, msg := range msgs {
			if _, ok := msg.(*types.MsgActivateAccount); ok {
				return nil, sdkerrors.Wrap(types.ErrInvalidTx, "smart account activation message must stand alone")
			}
		}

		return nil, nil
	}

	activateMsg, ok := msgs[0].(*types.MsgActivateAccount)
	if !ok {
		return nil, nil
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return nil, err
	}

	signer := sigTx.GetSigners()

	// do not allow multi signer and signature
	if len(signer) != 1 || len(sigs) != 1 {
		return nil, sdkerrors.Wrap(types.ErrInvalidTx, "smart-account activation tx does not allow multiple signers")
	}

	return activateMsg, nil
}

// ------------------------- SmartAccount Decorator ------------------------- \\

type SmartAccountDecorator struct {
	SaKeeper sakeeper.Keeper
}

func NewSmartAccountDecorator(saKeeper sakeeper.Keeper) *SmartAccountDecorator {
	return &SmartAccountDecorator{
		SaKeeper: saKeeper,
	}
}

// AnteHandle is used for performing basic validity checks on a transaction such that it can be thrown out of the mempool.
func (d *SmartAccountDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	activateMsg, err := GetValidActivateAccountMessage(sigTx)
	if err != nil {
		return ctx, err
	}

	if activateMsg == nil {
		newCtx, err = HandleSmartAccountTx(ctx, d.SaKeeper, sigTx, simulate)
		if err != nil {
			return newCtx, err
		}
	} else {
		newCtx, err = HandleSmartAccountActivate(ctx, d.SaKeeper, activateMsg, simulate)
		if err != nil {
			return newCtx, err
		}
	}

	return next(newCtx, tx, simulate)
}

func HandleSmartAccountTx(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	sigTx authsigning.SigVerifiableTx,
	simulate bool,
) (sdk.Context, error) {

	signerAcc, err := GetSmartAccountTxSigner(ctx, sigTx, saKeeper)
	if err != nil {
		return ctx, err
	}

	// if is not smartaccount tx type
	if signerAcc == nil {
		// do some thing
		return ctx, nil
	}

	// not support smartaccount tx simulation yet
	if simulate {
		return ctx, sdkerrors.Wrap(types.ErrNotSupported, "Simulation of SmartAccount txs isn't supported yet")
	}

	msgs := sigTx.GetMsgs()

	execMsg, err := ValidateAndGetAfterExecMessage(msgs, signerAcc)
	if err != nil {
		return ctx, err
	}

	// parse messages in tx to list of string
	execMsgData, err := types.ParseMessagesString(msgs[:len(msgs)-1])
	if err != nil {
		return ctx, err
	}

	// create message for SA contract pre-exeucte
	validateMessage, err := GeneratePreExecuteMessage(execMsg, execMsgData)
	if err != nil {
		return ctx, err
	}

	params := saKeeper.GetParams(ctx)

	// execute SA contract for pre-execute transaction with limit gas
	err = executeWithGasLimit(ctx, saKeeper.ContractKeeper, signerAcc.GetAddress(), validateMessage, params.MaxGasExecute)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func HandleSmartAccountActivate(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	activateMsg *types.MsgActivateAccount,
	simulate bool,
) (sdk.Context, error) {
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
		saKeeper.WasmKeeper,
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
		sAccount := saKeeper.AccountKeeper.GetAccount(ctx, signer)
		_, isBase := sAccount.(*authtypes.BaseAccount)
		_, isSa := sAccount.(*types.SmartAccount)
		if !isBase && !isSa {
			return ctx, sdkerrors.Wrap(types.ErrAccountNotFoundForAddress, signer.String())
		}

		// remove temporary pubkey for account
		err = saKeeper.UpdateAccountPubKey(ctx, sAccount, nil)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func ValidateAndGetAfterExecMessage(msgs []sdk.Msg, signerAcc *types.SmartAccount) (*wasmtypes.MsgExecuteContract, error) {
	// after-execute message must be the last message and must be MsgExecuteContract
	var afterExecMsg *wasmtypes.MsgExecuteContract
	if msg, err := msgs[len(msgs)-1].(*wasmtypes.MsgExecuteContract); err {
		afterExecMsg = msg
	} else {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "after-execute message must be type MsgExecuteContract")
	}

	// get smartaccount address
	saAddress := signerAcc.GetAddress().String()

	// the message must be sent from the signer's address which is also the smart contract address
	if afterExecMsg.Sender != saAddress || afterExecMsg.Contract != saAddress {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "sender address and contract address must be the same")
	}

	return afterExecMsg, nil
}

// Call a contract's execute with a gas limit
// referenced from Osmosis' protorev posthandler:
// https://github.com/osmosis-labs/osmosis/blob/98025f185ab2ee1b060511ed22679112abcc08fa/x/protorev/keeper/posthandler.go#L42-L43
func executeWithGasLimit(
	ctx sdk.Context, contractKeeper *wasmkeeper.PermissionedKeeper,
	contractAddr sdk.AccAddress, msg []byte, maxGas sdk.Gas,
) error {
	cacheCtx, write := ctx.CacheContext()
	cacheCtx = cacheCtx.WithGasMeter(sdk.NewGasMeter(maxGas))

	if _, err := contractKeeper.Execute(
		cacheCtx,
		contractAddr, // contract address
		contractAddr, // signer, the smart account has the same address as the contract linked with it
		msg,
		sdk.NewCoins(), // empty funds
	); err != nil {
		return err
	}

	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return nil
}

func GeneratePreExecuteMessage(msg *wasmtypes.MsgExecuteContract, msgs []types.MsgData) ([]byte, error) {
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

	executeMessage, err := json.Marshal(&types.AccountMsg{
		PreExecuteTx: &types.PreExecuteTx{
			Msgs: msgs,
		},
	})
	if err != nil {
		return nil, err
	}

	return executeMessage, nil
}

// ------------------------- SetPubKey Decorator ------------------------- \\

type SetPubKeyDecorator struct {
	saKeeper sakeeper.Keeper
}

func NewSetPubKeyDecorator(saKeeper sakeeper.Keeper) *SetPubKeyDecorator {
	return &SetPubKeyDecorator{
		saKeeper: saKeeper,
	}
}

func (d *SetPubKeyDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	activateMsg, err := GetValidActivateAccountMessage(sigTx)
	if err != nil {
		return ctx, err
	}

	// if is smart account activation message
	if activateMsg != nil {
		// get message signer
		signer := activateMsg.GetSigners()[0]

		// get smart contract account by address, account must be inactivate smart account
		sAccount, err := d.saKeeper.IsInactiveAccount(ctx, signer)
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
		err = d.saKeeper.UpdateAccountPubKey(ctx, sAccount, pubKey)
		if err != nil {
			return ctx, err
		}

		return next(ctx, tx, simulate)
	}

	signerAcc, err := GetSmartAccountTxSigner(ctx, sigTx, d.saKeeper)
	if err != nil {
		return ctx, err
	}

	// if is smart account tx skip authante NewSetPubKeyDecorator
	// need this to avoid pubkey and address equal check of above decorator
	if signerAcc != nil {
		// if this is smart account tx, check if pubkey is set
		if signerAcc.GetPubKey() == nil {
			return ctx, sdkerrors.Wrap(types.ErrNilPubkey, signerAcc.String())
		}

		// Also emit the following events, so that txs can be indexed by these
		var events sdk.Events
		events = append(events, sdk.NewEvent(types.EventTypeSmartAccountTx,
			sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signerAcc.GetAddress(), signerAcc.GetSequence())),
		))

		// maybe need add more event here

		ctx.EventManager().EmitEvents(events)

		return next(ctx, tx, simulate)
	}

	// default authant SetPubKeyDecorator
	svd := authante.NewSetPubKeyDecorator(d.saKeeper.AccountKeeper)
	return svd.AnteHandle(ctx, tx, simulate, next)
}
