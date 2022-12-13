package keeper

import (
	"github.com/aura-nw/aura/x/txfees/types"
)

var _ types.QueryServer = Keeper{}
