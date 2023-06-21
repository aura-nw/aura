package types

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	fmt "fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InstantiateSalt struct {
	Owner   string `json:"owner"`
	CodeID  uint64 `json:"code_id"`
	InitMsg []byte `json:"init_msg"`
	PubKey  []byte `json:"pub_key"`
}

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
		return nil, fmt.Errorf(ErrAddressFromBech32, err)
	}

	salt, err := GenerateSalt(owner, codeId, initMsg, pubKey)
	if err != nil {
		return nil, err
	}

	codeInfo := wasmKeeper.GetCodeInfo(ctx, codeId)
	if codeInfo == nil {
		return nil, fmt.Errorf(ErrNoSuchCodeID, codeId)
	}

	addrGenerator := wasmkeeper.PredicableAddressGenerator(ownerAcc, salt, initMsg, true)
	contractAddress := addrGenerator(ctx, codeId, codeInfo.CodeHash)
	if wasmKeeper.HasContractInfo(ctx, contractAddress) {
		return nil, fmt.Errorf(ErrInstantiateDuplicate)
	}

	return contractAddress, nil
}

func PubKeyDecode(raw string) (*secp256k1.PubKey, error) {
	bz, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf(ErrBadPublicKey, err.Error())
	}

	// secp25k61 public key
	pubKey := &secp256k1.PubKey{Key: nil}
	keyErr := pubKey.UnmarshalAmino(bz)
	if keyErr != nil {
		return nil, fmt.Errorf(ErrBadPublicKey, keyErr.Error())
	}

	return pubKey, nil
}
