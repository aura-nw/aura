package smartaccount

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sakeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
)

// Return tx' signer as SmartAccount, if not return nil
func GetSmartAccountTxSigner(ctx sdk.Context, sigTx authsigning.SigVerifiableTx, saKeeper sakeeper.Keeper) (*typesv1.SmartAccount, error) {
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

func GetValidActivateAccountMessage(sigTx authsigning.SigVerifiableTx) (*typesv1.MsgActivateAccount, error) {
	msgs := sigTx.GetMsgs()

	if len(msgs) != 1 {
		// smart account activation message must stand alone
		// it will prevent bundling of multiple messages with the activating message as the first message to avoid AnteHandler check
		for _, msg := range msgs {
			if _, ok := msg.(*typesv1.MsgActivateAccount); ok {
				return nil, errorsmod.Wrap(types.ErrInvalidTx, "smart account activation message must stand alone")
			}
		}

		return nil, nil
	}

	activateMsg, ok := msgs[0].(*typesv1.MsgActivateAccount)
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
		return nil, errorsmod.Wrap(types.ErrInvalidTx, "smart-account activation tx does not allow multiple signers")
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
		return ctx, errorsmod.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(types.ErrInvalidTx, "not a FeeTx")
	}

	activateMsg, err := GetValidActivateAccountMessage(sigTx)
	if err != nil {
		return ctx, err
	}

	if activateMsg == nil {
		err = handleSmartAccountTx(ctx, d.SaKeeper, sigTx, feeTx, simulate)
		if err != nil {
			return ctx, err
		}
	} else {
		err = handleSmartAccountActivate(ctx, d.SaKeeper, activateMsg, simulate)
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

func handleSmartAccountTx(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	sigTx authsigning.SigVerifiableTx,
	feeTx sdk.FeeTx,
	simulate bool,
) error {

	signerAcc, err := GetSmartAccountTxSigner(ctx, sigTx, saKeeper)
	if err != nil {
		return err
	}

	// if is not smartaccount tx type
	if signerAcc == nil {
		// do some thing
		return nil
	}

	// save the account address to the module store. we will need it in the
	// posthandler
	//
	// TODO: a question is that instead of writing to store, can we just put this
	// in memory instead. in practice however, the address is deleted in the post
	// handler, so it's never actually written to disk, meaning the difference in
	// gas consumption should be really small. still worth investigating tho.
	// referenced from Larry0x' abstractaccount AnteHandler:
	// https://github.com/larry0x/abstract-account/blob/b3c6432e593d450e7c58dae94cdf2a95930f8159/x/abstractaccount/ante.go#L81
	saKeeper.SetSignerAddress(ctx, signerAcc.GetAddress())

	msgs := sigTx.GetMsgs()

	// check if tx messages is allowed for smartaccount
	// permitted messages will be determined by the government
	err = saKeeper.CheckAllowedMsgs(ctx, msgs)
	if err != nil {
		return err
	}

	msgsData, err := types.ParseMessagesString(msgs)
	if err != nil {
		return err
	}

	callInfo := types.CallInfo{
		Gas:        feeTx.GetGas(),
		Fee:        feeTx.GetFee(),
		FeePayer:   feeTx.FeePayer().String(),
		FeeGranter: feeTx.FeeGranter().String(),
	}

	emptyAuthzInfo := types.AuthzInfo{
		Grantee: "",
	}

	preExecuteMessage, err := json.Marshal(&types.AccountMsg{
		PreExecuteTx: &types.PreExecuteTx{
			Msgs:      msgsData,
			CallInfo:  callInfo,
			AuthzInfo: emptyAuthzInfo,
		},
	})
	if err != nil {
		return err
	}

	params := saKeeper.GetParams(ctx)

	// execute SA contract for pre-execute transaction with limit gas
	// will using cacheCtx instead of ctx to run contract execution
	gasRemaining, err := sudoWithGasLimit(ctx, saKeeper.ContractKeeper, signerAcc.GetAddress(), preExecuteMessage, params.MaxGasExecute)
	if err != nil {
		return err
	}

	// free gas remaining after validate smartaccount msgs
	saKeeper.SetGasRemaining(ctx, gasRemaining)

	return nil
}

func handleSmartAccountActivate(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	activateMsg *typesv1.MsgActivateAccount,
	simulate bool,
) error {

	// get signer of smart account activation message
	signer := activateMsg.GetSigners()[0]

	// decode string to pubkey
	pubKey, err := typesv1.PubKeyDecode(activateMsg.PubKey)
	if err != nil {
		return err
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
		return err
	}

	// the signer of the activation message must be the address generated by Instantiate2's PredicableAddressGenerator
	if !signer.Equals(predicAddr) {
		return errorsmod.Wrap(types.ErrInvalidAddress, "not the same as predicted")
	}

	if !ctx.IsReCheckTx() && !simulate {
		// get smart contract account by address
		sAccount := saKeeper.AccountKeeper.GetAccount(ctx, signer)
		_, isBase := sAccount.(*authtypes.BaseAccount)
		if !isBase {
			return errorsmod.Wrap(types.ErrAccountNotFoundForAddress, signer.String())
		}

		// remove temporary pubkey for account
		// wasmd instantiate2 method requires base account with nil pubkey and 0 sequence to be considered an initializable address for smart contract
		err = saKeeper.UpdateAccountPubKey(ctx, sAccount, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Call a contract's sudo with a gas limit
//
// using gas for validate SA msgs will not count to tx total used
//
// referenced from Osmosis' protorev posthandler:
// https://github.com/osmosis-labs/osmosis/blob/98025f185ab2ee1b060511ed22679112abcc08fa/x/protorev/keeper/posthandler.go#L42-L43
func sudoWithGasLimit(
	ctx sdk.Context, contractKeeper *wasmkeeper.PermissionedKeeper,
	contractAddr sdk.AccAddress, msg []byte, maxGas sdk.Gas,
) (uint64, error) {
	cacheCtx, write := ctx.CacheContext()
	cacheCtx = cacheCtx.WithGasMeter(sdk.NewGasMeter(maxGas))

	if _, err := contractKeeper.Sudo(
		cacheCtx,
		contractAddr, // contract address
		msg,
	); err != nil {
		return maxGas, err
	}

	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return cacheCtx.GasMeter().GasRemaining(), nil
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
		return ctx, errorsmod.Wrap(types.ErrInvalidTx, "not a SigVerifiableTx")
	}

	activateMsg, err := GetValidActivateAccountMessage(sigTx)
	if err != nil {
		return ctx, err
	}

	// if is smart account activation message
	if activateMsg != nil {
		// get message signer
		signer := activateMsg.GetSigners()[0]

		// get smartaccount by address, must be inactivate account
		sAccount, err := d.saKeeper.IsInactiveAccount(ctx, signer)
		if err != nil {
			return ctx, err
		}

		if !ctx.IsReCheckTx() && !simulate {
			// decode any to pubkey
			pubKey, err := typesv1.PubKeyDecode(activateMsg.PubKey)
			if err != nil {
				return ctx, err
			}

			// set temporary pubkey for account
			// need this for the next ante signature checks
			err = d.saKeeper.UpdateAccountPubKey(ctx, sAccount, pubKey)
			if err != nil {
				return ctx, err
			}
		}

		return next(ctx, tx, simulate)
	}

	signerAcc, err := GetSmartAccountTxSigner(ctx, sigTx, d.saKeeper)
	if err != nil {
		return ctx, err
	}

	// if is smart account tx skip authante NewSetPubKeyDecorator
	// need this to avoid pubkey and address equal check of authante SetPubKeyDecorator
	if signerAcc != nil {
		// if this is smart account tx and not in simulation mode
		// check if pubkey is set
		if !simulate && signerAcc.GetPubKey() == nil {
			return ctx, errorsmod.Wrap(types.ErrNilPubkey, signerAcc.String())
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

// ------------------------- ValidateAuthzTx Decorator ------------------------- \\

type ValidateAuthzTxDecorator struct {
	SaKeeper sakeeper.Keeper
}

func NewValidateAuthzTxDecorator(saKeeper sakeeper.Keeper) *ValidateAuthzTxDecorator {
	return &ValidateAuthzTxDecorator{
		SaKeeper: saKeeper,
	}
}

// If smartaccount messages is executed throught MsgAuthzExec, SmartAccountDecorator cannot detect and validate these messages
// using AuthzTxDecorator AnteHandler to provide validation for nested smartaccount msgs
func (d *ValidateAuthzTxDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {

	params := d.SaKeeper.GetParams(ctx)
	maxGas := params.MaxGasExecute

	if d.SaKeeper.HasGasRemaining(ctx) {
		// if pre ante handlers has used free gas, get the remaining
		maxGas = d.SaKeeper.GetGasRemaining(ctx)
		d.SaKeeper.DeleteGasRemaining(ctx)
	}

	cacheCtx, write := ctx.CacheContext()
	cacheCtx = cacheCtx.WithGasMeter(sdk.NewGasMeter(maxGas))

	// using gas for validate authz will not count to tx total used
	err = validateAuthzTx(cacheCtx, d.SaKeeper, tx, true)
	if err != nil {
		return ctx, err
	}

	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return next(ctx, tx, simulate)
}

func validateAuthzTx(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	tx sdk.Tx,
	isAnte bool,
) error {

	for _, msg := range tx.GetMsgs() {
		if msgExec, ok := msg.(*authz.MsgExec); ok {
			msgs, err := msgExec.GetMessages()
			if err != nil {
				return err
			}

			err = validateNestedSmartAccountMsgs(
				ctx,
				saKeeper,
				msgs,
				msgExec.Grantee,
				isAnte,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateNestedSmartAccountMsgs(
	ctx sdk.Context,
	saKeeper sakeeper.Keeper,
	msgs []sdk.Msg,
	grantee string,
	isAnte bool,
) error {

	for _, msg := range msgs {

		signers := msg.GetSigners()

		if len(signers) == 1 {
			acc, err := saKeeper.GetSmartAccountByAddress(ctx, signers[0])
			if err != nil {
				return err
			}

			if acc != nil {
				msgsData, err := types.ParseMessagesString([]sdk.Msg{msg})
				if err != nil {
					return err
				}

				// call_info is empty if message is executed throught authz exec
				callInfo := types.CallInfo{
					Gas:        0,
					Fee:        sdk.NewCoins(),
					FeePayer:   "",
					FeeGranter: "",
				}

				authzInfo := types.AuthzInfo{
					Grantee: grantee,
				}

				var execMsg []byte
				if isAnte {
					execMsg, err = json.Marshal(&types.AccountMsg{
						PreExecuteTx: &types.PreExecuteTx{
							Msgs:      msgsData,
							CallInfo:  callInfo,
							AuthzInfo: authzInfo,
						},
					})
				} else {
					execMsg, err = json.Marshal(&types.AccountMsg{
						AfterExecuteTx: &types.AfterExecuteTx{
							Msgs:      msgsData,
							CallInfo:  callInfo,
							AuthzInfo: authzInfo,
						},
					})
				}
				if err != nil {
					return err
				}

				if _, err := saKeeper.ContractKeeper.Sudo(
					ctx,
					acc.GetAddress(),
					execMsg,
				); err != nil {
					return err
				}
			}
		}

		if msgExec, ok := msg.(*authz.MsgExec); ok {
			msgs, err := msgExec.GetMessages()
			if err != nil {
				return err
			}

			err = validateNestedSmartAccountMsgs(
				ctx,
				saKeeper,
				msgs,
				msgExec.Grantee,
				isAnte,
			)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
