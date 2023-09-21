package cli_test

import (
	"github.com/aura-nw/aura/x/aura/client/cli"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetQueryCmd(t *testing.T) {
	queryCmd := cli.GetQueryCmd("")
	require.NotNil(t, queryCmd)
}
