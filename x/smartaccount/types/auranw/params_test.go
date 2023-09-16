package auranw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aura-nw/aura/x/smartaccount/types/auranw"
)

func TestValidateParams(t *testing.T) {
	for _, tc := range []struct {
		params *auranw.Params
		expErr bool
	}{
		{
			params: &auranw.Params{},
			expErr: true,
		},
		{
			params: &auranw.Params{MaxGasExecute: 0},
			expErr: true,
		},
		{
			params: &auranw.Params{MaxGasExecute: auranw.DefaultMaxGas},
			expErr: false,
		},
		{
			params: &auranw.Params{
				MaxGasExecute:   auranw.DefaultMaxGas,
				DisableMsgsList: []string{"/cosmwasm.wasm.v1.MsgExecuteContract"},
			},
			expErr: false,
		},
		{
			params: &auranw.Params{
				MaxGasExecute:   auranw.DefaultMaxGas,
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
		allowedCodeIDs []*auranw.CodeID
		disableMsgs    []string
		codeID         uint64
		expAllowed     bool
	}{
		{
			allowedCodeIDs: []*auranw.CodeID{},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*auranw.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         88888,
			expAllowed:     false,
		},
		{
			allowedCodeIDs: []*auranw.CodeID{{69420, true}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     true,
		},
		{
			allowedCodeIDs: []*auranw.CodeID{{69420, false}},
			disableMsgs:    []string{},
			codeID:         69420,
			expAllowed:     false,
		},
	} {
		params := auranw.NewParams(tc.allowedCodeIDs, tc.disableMsgs, auranw.DefaultMaxGas)

		allowed := params.IsAllowedCodeID(tc.codeID)
		require.Equal(t, tc.expAllowed, allowed)
	}
}
