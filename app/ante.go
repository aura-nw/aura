package app

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"

	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	// ethermint ante
	sdkvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmosante "github.com/evmos/evmos/v16/app/ante"
	cosmosante "github.com/evmos/evmos/v16/app/ante/cosmos"
	"github.com/evmos/evmos/v16/app/ante/evm"
	evmante "github.com/evmos/evmos/v16/app/ante/evm"
	evmostypes "github.com/evmos/evmos/v16/types"
	evmtypes "github.com/evmos/evmos/v16/x/evm/types"

	smartaccount "github.com/aura-nw/aura/x/smartaccount"
	smartaccountkeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions
	WasmKeeper         wasmkeeper.Keeper
	SmartAccountKeeper smartaccountkeeper.Keeper
	IBCKeeper          *ibckeeper.Keeper
	WasmConfig         *wasmTypes.WasmConfig
	TXCounterStoreKey  storetypes.StoreKey
	Codec              codec.BinaryCodec
	EvmKeeper          evm.EVMKeeper
	FeeMarketKeeper    evm.FeeMarketKeeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func (app *App) NewAnteHandler(txConfig client.TxConfig, wasmConfig wasmTypes.WasmConfig, wasmKey storetypes.StoreKey) (sdk.AnteHandler, error) {
	// return auraAnteHandler(options)
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		ext, ok := tx.(authante.HasExtensionOptionsTx)
		opts := ext.GetExtensionOptions()
		if ok && len(opts) > 0 {
			// TODO: add config for max gas wanted
			ctx.Logger().Info(fmt.Sprintf("tx: %v", tx))
			maxGasWanted := cast.ToUint64(100000000)
			var evmosoptions = evmosante.HandlerOptions{
				Cdc:                    app.appCodec,
				AccountKeeper:          app.AccountKeeper,
				BankKeeper:             app.BankKeeper,
				ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
				EvmKeeper:              app.EvmKeeper,
				StakingKeeper:          app.StakingKeeper,
				FeegrantKeeper:         app.FeeGrantKeeper,
				DistributionKeeper:     app.DistrKeeper, // what is this
				IBCKeeper:              app.IBCKeeper,
				FeeMarketKeeper:        app.FeeMarketKeeper,
				SignModeHandler:        txConfig.SignModeHandler(),
				SigGasConsumer:         evmosante.SigVerificationGasConsumer,
				MaxTxGasWanted:         maxGasWanted,
				TxFeeChecker:           evmante.NewDynamicFeeChecker(app.EvmKeeper),
			}
			anteHandler = evmosante.NewAnteHandler(evmosoptions)
			return anteHandler(ctx, tx, sim)
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			var options = HandlerOptions{
				HandlerOptions: ante.HandlerOptions{
					AccountKeeper:   app.AccountKeeper,
					BankKeeper:      app.BankKeeper,
					SignModeHandler: txConfig.SignModeHandler(),
					FeegrantKeeper:  app.FeeGrantKeeper,
					SigGasConsumer:  ante.DefaultSigVerificationGasConsumer},
				WasmKeeper:         app.WasmKeeper,
				SmartAccountKeeper: app.SaKeeper,
				IBCKeeper:          app.IBCKeeper,
				WasmConfig:         &wasmConfig,
				TXCounterStoreKey:  wasmKey,
				Codec:              app.appCodec,
			}
			anteHandler = auraAnteHandler(options)
		default:
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}, nil
}

func auraAnteHandler(options HandlerOptions) sdk.AnteHandler {
	// if options.AccountKeeper == nil {
	// 	return nil, errorsmod.Wkrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	// }

	// if options.BankKeeper == nil {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	// }

	// if options.SignModeHandler == nil {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	// }

	// if options.WasmConfig == nil {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	// }

	// if options.TXCounterStoreKey == nil {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	// }

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		cosmosante.RejectMessagesDecorator{}, // reject MsgEthereumTx
		cosmosante.NewAuthzLimiterDecorator( // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&sdkvesting.MsgCreateVestingAccount{}),
		),

		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		// limit simulation gas
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		// stargazeante.NewMinCommissionDecorator(options.Codec),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		smartaccount.NewSetPubKeyDecorator(options.SmartAccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),

		// new ante for account abstraction
		smartaccount.NewSmartAccountDecorator(options.SmartAccountKeeper),
		smartaccount.NewValidateAuthzTxDecorator(options.SmartAccountKeeper),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),

		// ethermint ante
		// evmante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...)
}
