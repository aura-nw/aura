package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const TypeQueryGenerateAccountRequest = "query_generate_account"

var (
	_ codectypes.UnpackInterfacesMessage = (*QueryGenerateAccountRequest)(nil)
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg QueryGenerateAccountRequest) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.PubKey == nil {
		return nil
	}
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.PubKey, &pubKey)
}
