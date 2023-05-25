package ante

import (
	"encoding/json"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	smartaccountkeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type SmartAccountDecorator struct {
	SmartAccountKeeper smartaccountkeeper.Keeper
	WasmKeeper         wasmkeeper.Keeper
}

type UserOps struct {
	Messages string `json:"messages"`
}

type ValidateUserOps struct {
	Validate UserOps `json:"validate"`
}

type PreExecuteUserOps struct {
	PreExecute UserOps `json:"pre_execute"`
}

type ValidateUserOpsResponse = bool

func NewSmartAccountDecorator(smartAccountKeeper smartaccountkeeper.Keeper, wasmKeeper wasmkeeper.Keeper) *SmartAccountDecorator {
	return &SmartAccountDecorator{
		SmartAccountKeeper: smartAccountKeeper,
		WasmKeeper:         wasmKeeper,
	}
}

func generateValidateQueryMessage(msg *wasmtypes.MsgExecuteContract, msgs []MsgData) ([]byte, error) {

	var preExecuteUserOps PreExecuteUserOps
	umErr := json.Unmarshal(msg.GetMsg(), &preExecuteUserOps)
	if umErr != nil {
		return nil, fmt.Errorf("invalid pre-execute message data: %s", umErr.Error())
	}

	valMsgData, err := json.Marshal(&msgs)
	if err != nil {
		return nil, fmt.Errorf("cannot json marshal pre-execute message: %s", err.Error())
	}

	// data in validate message must compatiable to tx.messages
	if string(valMsgData) != preExecuteUserOps.PreExecute.Messages {
		return nil, fmt.Errorf("invalid pre-execute message data: not compatible with tx.messages")
	}

	validateUserOps := ValidateUserOps{
		Validate: preExecuteUserOps.PreExecute,
	}
	validateMessage, err := json.Marshal(validateUserOps)
	if err != nil {
		return nil, fmt.Errorf("cannot json marshal validate message: %s", err.Error())
	}

	return validateMessage, nil
}

func (decorator SmartAccountDecorator) isSmartAccountTx(ctx sdk.Context, sigTx authsigning.SigVerifiableTx) bool {
	signerAddrs := sigTx.GetSigners()

	// not support multi signer yet
	if len(signerAddrs) != 1 {
		return false
	}

	signerAddr := signerAddrs[0]

	// get smartaccount value from module keeper
	accountValue := decorator.SmartAccountKeeper.GetSmartAccount(ctx, signerAddr.String())

	// account type must be not empty and is active
	return accountValue.Type != "" && accountValue.Active
}

// AnteHandle is used for performing basic validity checks on a transaction such that it can be thrown out of the mempool.
func (decorator *SmartAccountDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, fmt.Errorf("invalid transaction type")
	}

	// if is not smartaccount tx type
	if !decorator.isSmartAccountTx(ctx, sigTx) {
		// do some thing
		return next(ctx, tx, simulate)
	}

	msgs := sigTx.GetMsgs()

	if len(msgs) < 1 {
		return ctx, fmt.Errorf("tx must contain at least the validate message")
	}

	// validate message must be the first message and must be MsgExecuteContract
	var valMsg *wasmtypes.MsgExecuteContract
	if msg, err := msgs[0].(*wasmtypes.MsgExecuteContract); err {
		valMsg = msg
	} else {
		return ctx, fmt.Errorf("validate message must be type MsgExecuteContract")
	}

	// get smartaccount address
	smartAccount := sigTx.GetSigners()[0]
	saAddress := smartAccount.String()

	// the message must be sent from the signer's address which is also the smart contract address
	if valMsg.Sender != saAddress || valMsg.Contract != saAddress {
		return ctx, fmt.Errorf(
			"invalid validate message: expected (%s),(%s) - got (%s),(%s)",
			saAddress, saAddress,
			valMsg.Sender, valMsg.Contract,
		)
	}

	// parse messages in tx to list of string
	valMsgData, err := parseMessagesString(msgs[1:])
	if err != nil {
		return ctx, err
	}

	// create message for SA contract query
	validateMessage, err := generateValidateQueryMessage(valMsg, valMsgData)
	if err != nil {
		return ctx, err
	}

	// get contract address
	contractAddr, err := sdk.AccAddressFromBech32(valMsg.Contract)
	if err != nil {
		return ctx, fmt.Errorf("cannot convert bech32 to account address: %s", err.Error())
	}

	// query SA contract for validating transaction
	rawQueryResponse, err := decorator.WasmKeeper.QuerySmart(ctx, contractAddr, validateMessage)
	if err != nil {
		return ctx, fmt.Errorf("transaction validate fail: %s", err.Error())
	}

	var queryResponse ValidateUserOpsResponse
	rErr := json.Unmarshal(rawQueryResponse, &queryResponse)
	if rErr != nil {
		return ctx, fmt.Errorf("cannot json marshal validate response message: %s", rErr.Error())
	}

	// pre-validate response
	if !queryResponse {
		return ctx, fmt.Errorf("transaction validate fail")
	}

	return next(ctx, tx, simulate)
}
