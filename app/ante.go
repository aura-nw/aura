package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	smartaccountkeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcante "github.com/cosmos/ibc-go/v4/modules/core/ante"
	keeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"

	// stargazeante "github.com/public-awesome/stargaze/v4/internal/ante"
	smartaccount "github.com/aura-nw/aura/x/smartaccount"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions
	WasmKeeper         wasmkeeper.Keeper
	SmartAccountKeeper smartaccountkeeper.Keeper
	IBCKeeper          *keeper.Keeper
	WasmConfig         *wasmTypes.WasmConfig
	TXCounterStoreKey  sdk.StoreKey
	Codec              codec.BinaryCodec
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
		// stargazeante.NewMinCommissionDecorator(options.Codec),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		smartaccount.NewSetPubKeyDecorator(options.AccountKeeper, options.WasmKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),

		// new ante for account abstraction
		smartaccount.NewSmartAccountDecorator(options.WasmKeeper, options.AccountKeeper),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewAnteDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
