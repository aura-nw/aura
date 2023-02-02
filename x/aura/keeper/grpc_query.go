package keeper

import (
	"github.com/aura-nw/aura/x/aura/types"
)

var _ types.QueryServer = Keeper{}
