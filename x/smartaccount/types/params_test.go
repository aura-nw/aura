package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aura-nw/aura/x/smartaccount/types"
)

func TestValidateParams(t *testing.T) {
	for _, tc := range []struct {
		params *types.Params
		expErr bool
	}{
		{
			params: &types.Params{},
			expErr: true,
		},
		{
			params: &types.Params{MaxGasExecute: 0},
			expErr: true,
		},
		{
			params: &types.Params{MaxGasExecute: types.DefaultMaxGas},
			expErr: false,
		},
	} {
		err := tc.params.Validate()

		if tc.expErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestDeterminedAllowedCodeID(t *testing.T) {
	for _, tc := range []struct {
		allowedCodeIDs []*types.CodeID
		codeID         uint64
		expAllowed     bool
	}{
		{
			allowedCodeIDs: []*types.CodeID{},
			codeID:         69420,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, true}},
			codeID:         88888,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, true}},
			codeID:         69420,
			expAllowed:     true,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, false}},
			codeID:         69420,
			expAllowed:     false,
		},
	} {
		params := types.NewParams(tc.allowedCodeIDs, types.DefaultMaxGas)

		allowed := params.IsAllowedCodeID(tc.codeID)
		require.Equal(t, tc.expAllowed, allowed)
	}
}
