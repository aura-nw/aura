package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	channelkeeper "github.com/cosmos/ibc-go/v2/modules/core/04-channel/keeper"
	ibcante "github.com/cosmos/ibc-go/v2/modules/core/ante"
)

// // NewAnteHandler returns an AnteHandler that checks and increments sequence
// // numbers, checks signatures & account numbers, and deducts fees from the first
// // signer.
// func NewAnteHandler(
// 	ak ante.AccountKeeper,
// 	bankKeeper types.BankKeeper,
// 	sigGasConsumer ante.SignatureVerificationGasConsumer,
// 	signModeHandler signing.SignModeHandler,
// 	txCounterStoreKey sdk.StoreKey,
// 	channelKeeper channelkeeper.Keeper,
// 	fk ante.FeegrantKeeper,
// 	// wasmConfig wasmTypes.WasmConfig,
// ) sdk.AnteHandler {
// 	// copied sdk https://github.com/cosmos/cosmos-sdk/blob/v0.42.9/x/auth/ante/ante.go

// 	return sdk.ChainAnteDecorators(
// 		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
// 		// wasmkeeper.NewLimitSimulationGasDecorator(wasmConfig.SimulationGasLimit), // after setup context to enforce limits early
// 		// wasmkeeper.NewCountTXDecorator(txCounterStoreKey),
// 		ante.NewRejectExtensionOptionsDecorator(),
// 		ante.NewMempoolFeeDecorator(),
// 		ante.NewValidateBasicDecorator(),
// 		ante.TxTimeoutHeightDecorator{},
// 		ante.NewValidateMemoDecorator(ak),
// 		ante.NewConsumeGasForTxSizeDecorator(ak),
// 		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
// 		ante.NewValidateSigCountDecorator(ak),
// 		ante.NewDeductFeeDecorator(ak, bankKeeper, fk),
// 		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
// 		ante.NewSigVerificationDecorator(ak, signModeHandler),
// 		ante.NewIncrementSequenceDecorator(ak),
// 		ibcante.NewAnteDecorator(channelKeeper),
// 	)
// }

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions
	IBCChannelkeeper  channelkeeper.Keeper
	WasmConfig        *wasmTypes.WasmConfig
	TXCounterStoreKey sdk.StoreKey
	Codec             codec.BinaryCodec
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.WasmConfig == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}

	if options.TXCounterStoreKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		// limit simulation gas
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewAnteDecorator(options.IBCChannelkeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
