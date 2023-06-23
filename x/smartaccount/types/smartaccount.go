package types

import (
	"crypto/sha512"
	"encoding/json"
	fmt "fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

func PubKeyDecode(pubKey *codectypes.Any) (cryptotypes.PubKey, error) {
	pkAny := pubKey.GetCachedValue()
	pk, ok := pkAny.(cryptotypes.PubKey)
	if ok {
		return pk, nil
	} else {
		return nil, fmt.Errorf("expecting PubKey, got: %T", pkAny)
	}
}

// Inactivate smart-account must be base account with empty public key
func IsInactivateAccount(ctx sdk.Context, acc sdk.AccAddress, acc_str string, accountKeeper AccountKeeper) (authtypes.AccountI, error) {
	sAccount := accountKeeper.GetAccount(ctx, acc)

	// check if account has type base
	if _, ok := sAccount.(*authtypes.BaseAccount); !ok {
		return nil, fmt.Errorf(ErrAccountNotFoundForAddress, acc_str)
	}

	// check if account already has public key
	if sAccount.GetPubKey() != nil {
		return nil, fmt.Errorf(ErrAccountAlreadyExists)
	}

	return sAccount, nil
}

// Convert pubkey string to *Any
func PubKeyToAny(cdc codec.Codec, raw []byte) (*codectypes.Any, error) {
	var pubKey cryptotypes.PubKey
	err := cdc.UnmarshalInterfaceJSON(raw, &pubKey)
	if err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return nil, err
	}

	return any, nil
}
