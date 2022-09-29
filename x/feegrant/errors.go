package feegrant

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

var ErrAddressNotAllowed = sdkerrors.Register(feegrant.DefaultCodespace, 8, "address not allowed")
