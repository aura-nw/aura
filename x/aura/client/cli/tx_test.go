package cli_test

import (
	"github.com/aura-nw/aura/x/aura/client/cli"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetTxCmd(t *testing.T) {
	txCmd := cli.GetTxCmd()
	require.Equal(t, len(txCmd.Commands()), 2)
}
