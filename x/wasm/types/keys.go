package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "wasm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_wasm"

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName
)

// nolint
var (
	CodeKeyPrefix                                  = []byte{0x01}
	ContractKeyPrefix                              = []byte{0x02}
	ContractStorePrefix                            = []byte{0x03}
	SequenceKeyPrefix                              = []byte{0x04}
	ContractCodeHistoryElementPrefix               = []byte{0x05}
	ContractByCodeIDAndCreatedSecondaryIndexPrefix = []byte{0x06}
	PinnedCodeIndexPrefix                          = []byte{0x07}
	TXCounterPrefix                                = []byte{0x08}

	KeyLastCodeID     = append(SequenceKeyPrefix, []byte("lastCodeId")...)
	KeyLastInstanceID = append(SequenceKeyPrefix, []byte("lastContractId")...)
)

// GetCodeKey constructs the key for retreiving the ID for the WASM code
func GetCodeKey(codeID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(codeID)
	return append(CodeKeyPrefix, contractIDBz...)
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(ContractKeyPrefix, addr...)
}

// GetContractStorePrefix returns the store prefix for the WASM contract instance
func GetContractStorePrefix(addr sdk.AccAddress) []byte {
	return append(ContractStorePrefix, addr...)
}

// GetPinnedCodeIndexPrefix returns the key prefix for a code id pinned into the wasmvm cache
func GetPinnedCodeIndexPrefix(codeID uint64) []byte {
	prefixLen := len(PinnedCodeIndexPrefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], PinnedCodeIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(codeID))
	return r
}
