package cli

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// ParseAddrFromHexOrBech32 parses a string address that can be either a hex or
// Bech32 string.
func ParseAddrFromHexOrBech32(addrString string) (common.Address, error) {
	if common.IsHexAddress(addrString) {
		return common.HexToAddress(addrString), nil
	}

	cfg := sdk.GetConfig()

	if !strings.HasPrefix(addrString, cfg.GetBech32AccountAddrPrefix()) {
		return common.Address{}, fmt.Errorf("receiver '%s' is not a hex or bech32 address (prefix does not match)", addrString)
	}

	accAddr, err := sdk.AccAddressFromBech32(addrString)
	if err != nil {
		return common.Address{}, fmt.Errorf("receiver '%s' is not a hex or bech32 address (could not parse as bech32 string)", addrString)
	}

	return common.BytesToAddress(accAddr), nil

}
