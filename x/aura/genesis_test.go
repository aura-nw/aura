package aura_test

import (
	"testing"

	keepertest "github.com/aura-nw/aura/testutil/keeper"
	"github.com/aura-nw/aura/x/aura"
	"github.com/aura-nw/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AuraKeeper(t)
	aura.InitGenesis(ctx, *k, genesisState)
	got := aura.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
