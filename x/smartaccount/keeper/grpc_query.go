package keeper

import (
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
)

var _ typesv1.QueryServer = Keeper{}
