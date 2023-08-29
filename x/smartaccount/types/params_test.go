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
		{
			params: &types.Params{
				MaxGasExecute:   types.DefaultMaxGas,
				DisableMsgsList: []string{"/cosmwasm.wasm.v1.MsgExecuteContract"},
			},
			expErr: false,
		},
		{
			params: &types.Params{
				MaxGasExecute:   types.DefaultMaxGas,
				DisableMsgsList: []string{"/cosmwasm.wasm.v1.MsgExecuteContract", "/cosmwasm.wasm.v1.MsgExecuteContract"},
			},
			expErr: true,
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
		disableMsgs    []string
		codeID         uint64
		expAllowed     bool
	}{
		{
			allowedCodeIDs: []*types.CodeID{},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         88888,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     true,
		},
		{
			allowedCodeIDs: []*types.CodeID{{69420, false}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
	} {
		params := types.NewParams(tc.allowedCodeIDs, tc.disableMsgs, types.DefaultMaxGas)

		allowed := params.IsAllowedCodeID(tc.codeID)
		require.Equal(t, tc.expAllowed, allowed)
	}
}
