package wasm_test

import (
	"testing"

	keepertest "github.com/aura-nw/aura/testutil/keeper"
	"github.com/aura-nw/aura/x/wasm"
	"github.com/aura-nw/aura/x/wasm/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.WasmKeeper(t)
	wasm.InitGenesis(ctx, *k, genesisState)
	got := wasm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
