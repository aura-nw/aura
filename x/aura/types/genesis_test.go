package types_test

import (
	"testing"

	"github.com/aura-nw/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				// this line is used by starport scaffolding # types/genesis/validField
				Params: types.Params{
					MaxSupply: "10000000",
				},
			},
			valid: true,
		},
		{
			desc: "invalid max supply when is small",
			genState: &types.GenesisState{
				Params: types.Params{
					MaxSupply: "5",
				},
			},
			valid: false,
		},

		{
			desc: "limit list exclude address",
			genState: &types.GenesisState{
				Params: types.Params{
					ExcludeCirculatingAddr: []string{"addr1", "addr2", "addr3", "addr4", "addr5", "addr6", "addr7", "addr8", "addr9", "addr10", "addr11"},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate list exclude address",
			genState: &types.GenesisState{
				Params: types.Params{
					ExcludeCirculatingAddr: []string{"addr1", "addr1", "addr2"},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
