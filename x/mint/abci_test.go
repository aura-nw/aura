package mint_test

import (
	"github.com/aura-nw/aura/x/mint"
	"github.com/aura-nw/aura/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestBeginBlocker(t *testing.T) {
	type args struct {
		ctx types.Context
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mint.BeginBlocker(tt.args.ctx, tt.args.k)
		})
	}
}
