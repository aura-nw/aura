package keeper_test

import (
	"testing"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		params types.Params
		err    bool
	}{
		{
			desc: "error, duplicate codeID in whitlist",
			params: types.NewParams(
				[]*types.CodeID{{CodeID: 1, Status: true}, {CodeID: 1, Status: false}}, // duplicate codeID
				types.DefaultMaxGas,
			),
			err: true,
		},
		{
			desc: "error, max_gas_execute with zero value",
			params: types.NewParams(
				[]*types.CodeID{{CodeID: 1, Status: true}},
				uint64(0), // zero max gas execute
			),
			err: true,
		},
		{
			desc: "create new params successfully",
			params: types.NewParams(
				[]*types.CodeID{{CodeID: 1, Status: true}},
				types.DefaultMaxGas,
			),
			err: false,
		},
	} {
		ctx, app := helper.SetupGenesisTest()

		keeper := app.SaKeeper

		err := keeper.SetParams(ctx, tc.params)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)

			params := keeper.GetParams(ctx)
			require.Equal(t, tc.params, params)
		}
	}
}
