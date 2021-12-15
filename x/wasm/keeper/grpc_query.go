package keeper

import (
	"github.com/aura-nw/aura/x/wasm/types"
)

var _ types.QueryServer = Keeper{}
