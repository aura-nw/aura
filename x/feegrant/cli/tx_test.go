package cli_test

import (
	"github.com/aura-nw/aura/x/feegrant/cli"
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func TestGetTxCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.GetTxCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTxCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCmdFeeGrant(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.NewCmdFeeGrant(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCmdFeeGrant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCmdRevokeFeegrant(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.NewCmdRevokeFeegrant(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCmdRevokeFeegrant() = %v, want %v", got, tt.want)
			}
		})
	}
}
