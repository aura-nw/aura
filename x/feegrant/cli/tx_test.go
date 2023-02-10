package cli_test

import (
	"github.com/aura-nw/aura/x/feegrant/cli"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetTxCmd(t *testing.T) {
	txCmd := cli.GetTxCmd()
	require.NotNil(t, txCmd)
}
