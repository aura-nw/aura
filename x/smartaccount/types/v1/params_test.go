package v1_test

import (
	"testing"

	v1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	for _, tc := range []struct {
		params *v1.Params
		expErr bool
	}{
		{
			params: &v1.Params{},
			expErr: true,
		},
		{
			params: &v1.Params{MaxGasExecute: 0},
			expErr: true,
		},
		{
			params: &v1.Params{MaxGasExecute: v1.DefaultMaxGas},
			expErr: false,
		},
		{
			params: &v1.Params{
				MaxGasExecute:   v1.DefaultMaxGas,
				DisableMsgsList: []string{"/cosmwasm.wasm.v1.MsgExecuteContract"},
			},
			expErr: false,
		},
		{
			params: &v1.Params{
				MaxGasExecute:   v1.DefaultMaxGas,
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
		allowedCodeIDs []*v1.CodeID
		disableMsgs    []string
		codeID         uint64
		expAllowed     bool
	}{
		{
			allowedCodeIDs: []*v1.CodeID{},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*v1.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         88888,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*v1.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     true,
		},
		{
			allowedCodeIDs: []*v1.CodeID{{69420, false}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
	} {
		params := v1.NewParams(tc.allowedCodeIDs, tc.disableMsgs, v1.DefaultMaxGas)

		allowed := params.IsAllowedCodeID(tc.codeID)
		require.Equal(t, tc.expAllowed, allowed)
	}
}
