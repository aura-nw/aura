package feegrant

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

var ErrAddressNotAllowed = errorsmod.Register(feegrant.DefaultCodespace, 8, "address not allowed")
