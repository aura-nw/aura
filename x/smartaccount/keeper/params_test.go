package keeper_test

import (
	"testing"

	helper "github.com/aura-nw/aura/tests/smartaccount"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		params typesv1.Params
		err    bool
	}{
		{
			desc: "error, duplicate codeID in whitlist",
			params: typesv1.NewParams(
				[]*typesv1.CodeID{{CodeID: 1, Status: true}, {CodeID: 1, Status: false}}, // duplicate codeID
				[]string{},
				typesv1.DefaultMaxGas,
			),
			err: true,
		},
		{
			desc: "error, duplicate msg",
			params: typesv1.NewParams(
				[]*typesv1.CodeID{{CodeID: 1, Status: true}},
				[]string{"/cosmwasm.wasm.v1.MsgExecuteContract", "/cosmwasm.wasm.v1.MsgExecuteContract"}, // duplicate msg
				typesv1.DefaultMaxGas,
			),
			err: true,
		},
		{
			desc: "error, max_gas_execute with zero value",
			params: typesv1.NewParams(
				[]*typesv1.CodeID{{CodeID: 1, Status: true}},
				[]string{},
				uint64(0), // zero max gas execute
			),
			err: true,
		},
		{
			desc: "create new params successfully",
			params: typesv1.NewParams(
				[]*typesv1.CodeID{{CodeID: 1, Status: true}},
				[]string{"/cosmwasm.wasm.v1.MsgExecuteContract"},
				typesv1.DefaultMaxGas,
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
