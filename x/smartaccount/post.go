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

func (d AfterTxDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(types.ErrInvalidTx, "not a FeeTx")
	}

	// skip checkTx and re-checkTx
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		return next(ctx, tx, simulate, success)
	}

	// load the signer address, which we determined during the AnteHandler
	//
	// if not found, it means this tx is simply not an AA tx. we skip
	// referenced from Larry0x' abstractaccount posthandler:
	// https://github.com/larry0x/abstract-account/blob/b3c6432e593d450e7c58dae94cdf2a95930f8159/x/abstractaccount/ante.go#L153-L161
	signerAddr := d.saKeeper.GetSignerAddress(ctx)
	if signerAddr == nil {
		return next(ctx, tx, simulate, success)
	}

	d.saKeeper.DeleteSignerAddress(ctx)

	msgsData, err := types.ParseMessagesString(feeTx.GetMsgs())
	if err != nil {
		return ctx, err
	}

	callInfo := types.CallInfo{
		Gas:        feeTx.GetGas(),
		Fee:        feeTx.GetFee(),
		FeePayer:   feeTx.FeePayer().String(),
		FeeGranter: feeTx.FeeGranter().String(),
	}

	afterExecuteMessage, err := json.Marshal(&types.AccountMsg{
		AfterExecuteTx: &types.AfterExecuteTx{
			Msgs:     msgsData,
			CallInfo: callInfo,
			IsAuthz:  false,
		},
	})
	if err != nil {
		return ctx, err
	}

	// execute SA contract for after-execute transaction
	if _, err := d.saKeeper.ContractKeeper.Sudo(
		ctx,
		signerAddr, // contract address
		afterExecuteMessage,
	); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate, success)
}

// ------------------------- PostValidateAuthzTx Decorator ------------------------- \\

type PostValidateAuthzTxDecorator struct {
	SaKeeper sakeeper.Keeper
}

func NewPostValidateAuthzTxDecorator(saKeeper sakeeper.Keeper) *PostValidateAuthzTxDecorator {
	return &PostValidateAuthzTxDecorator{
		SaKeeper: saKeeper,
	}
}

func (d *PostValidateAuthzTxDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {

	// skip checkTx and re-checkTx
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		return next(ctx, tx, simulate, success)
	}

	err = validateAuthzTx(ctx, d.SaKeeper, tx, false)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate, success)
}
