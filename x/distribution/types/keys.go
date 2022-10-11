package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var DelegatorClaimAddrPrefix = []byte{0x09}

func GetDelegatorClaimAddrKey(delAddr sdk.AccAddress) []byte {
	return append(DelegatorClaimAddrPrefix, address.MustLengthPrefix(delAddr.Bytes())...)
}
