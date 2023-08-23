package smartaccount_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/types"
)

var (
	mockNextAccountID = types.DefaultSmartAccountId
)

func TestInitGenesis(t *testing.T) {
	ctx, app := helper.SetupGenesisTest()

	params := app.SaKeeper.GetParams(ctx)
	require.Equal(t, helper.GenesisState.Params, params)

	nextAccountID := app.SaKeeper.GetNextAccountID(ctx)
	require.Equal(t, mockNextAccountID, nextAccountID)
}

func TestExportGenesis(t *testing.T) {
	ctx, app := helper.SetupGenesisTest()

	if ctx.IsCheckTx() {
		fmt.Println("go check tx")
	}

	gs := smartaccount.ExportGenesis(ctx, app.SaKeeper)
	require.Equal(t, helper.GenesisState, gs)
}
