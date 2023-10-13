package smartaccount

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sakeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ------------------------- AfterTx Decorator ------------------------- \\
type AfterTxDecorator struct {
	saKeeper sakeeper.Keeper
}

func NewAfterTxDecorator(saKeeper sakeeper.Keeper) *AfterTxDecorator {
	return &AfterTxDecorator{
		saKeeper: saKeeper,
	}
}

// referenced from Larry0x' abstractaccount posthandler:
// https://github.com/larry0x/abstract-account/blob/b3c6432e593d450e7c58dae94cdf2a95930f8159/x/abstractaccount/ante.go#L152-L185
func (d AfterTxDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(types.ErrInvalidTx, "not a FeeTx")
	}

	// load the signer address, which we determined during the AnteHandler
	//
	// if not found, it means this tx is simply not an AA tx. we skip
	signerAddr := d.saKeeper.GetSignerAddress(ctx)
	if signerAddr == nil {
		return next(ctx, tx, simulate, success)
	}

	d.saKeeper.DeleteSignerAddress(ctx)

	msgsData, err := types.ParseMessagesString(feeTx.GetMsgs())
	if err != nil {
		return ctx, err
	}

	afterExecuteMessage, err := json.Marshal(&types.AccountMsg{
		AfterExecuteTx: &types.AfterExecuteTx{
			Msgs: msgsData,
		},
	})
	if err != nil {
		return ctx, err
	}

	params := d.saKeeper.GetParams(ctx)

	// execute SA contract for after-execute transaction with limit gas
	err = sudoWithGasLimit(ctx, d.saKeeper.ContractKeeper, signerAddr, afterExecuteMessage, params.MaxGasExecute)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate, success)
}
