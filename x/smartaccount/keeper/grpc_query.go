package keeper

import (
	"github.com/aura-nw/aura/x/smartaccount/types"
)

var _ types.QueryServer = Keeper{}
