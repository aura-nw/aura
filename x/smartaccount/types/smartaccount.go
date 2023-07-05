package types

import (
	"strconv"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// generate predictable contract address
func Instantiate2Address(
	ctx sdk.Context,
	wasmKeeper wasmkeeper.Keeper,
	codeId uint64,
	initMsg []byte,
	salt []byte,
	pubKey cryptotypes.PubKey,
) (sdk.AccAddress, error) {

	// we use pubkey.Address() as owner of this contract
	// remember this account doesn't exist on chain yet if have not received any funds before
	ownerAcc := sdk.AccAddress(pubKey.Address())

	codeInfo := wasmKeeper.GetCodeInfo(ctx, codeId)
	if codeInfo == nil {
		return nil, sdkerrors.Wrap(ErrNoSuchCodeID, strconv.FormatUint(codeId, 10))
	}

	addrGenerator := wasmkeeper.PredicableAddressGenerator(ownerAcc, salt, initMsg, true)
	contractAddress := addrGenerator(ctx, codeId, codeInfo.CodeHash)
	if wasmKeeper.HasContractInfo(ctx, contractAddress) {
		return nil, sdkerrors.Wrap(ErrInstantiateDuplicate, contractAddress.String())
	}

	return contractAddress, nil
}
