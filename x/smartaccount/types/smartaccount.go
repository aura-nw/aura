package types

import (
	"crypto/sha512"
	"encoding/json"
	"strconv"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type InstantiateSalt struct {
	Owner   string `json:"owner"`
	CodeID  uint64 `json:"code_id"`
	InitMsg []byte `json:"init_msg"`
	PubKey  []byte `json:"pub_key"`
}

// generate salt for contract Instantiate2
func GenerateSalt(owner string, codeId uint64, initMsg []byte, pubKey []byte) ([]byte, error) {
	salt := InstantiateSalt{
		Owner:   owner,
		CodeID:  codeId,
		InitMsg: initMsg,
		PubKey:  pubKey,
	}

	salt_bytes, err := json.Marshal(salt)
	if err != nil {
		return nil, err
	}

	// instantiate2 salt max length is 64 bytes, so need hash here
	salt_hashed := sha512.Sum512(salt_bytes)

	return salt_hashed[:], nil
}

// generate predictable contract address
func Instantiate2Address(
	ctx sdk.Context,
	wasmKeeper wasmkeeper.Keeper,
	owner string,
	codeId uint64,
	initMsg []byte,
	pubKey []byte,
) (sdk.AccAddress, error) {

	ownerAcc, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	salt, err := GenerateSalt(owner, codeId, initMsg, pubKey)
	if err != nil {
		return nil, err
	}

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
