package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aura-nw/aura/app"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
)

func TestAppStart(t *testing.T) {
	_, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	_ = app.New(logger, db, nil, true, map[int64]bool{}, filepath.Join(".", ".testapp"), 0, cosmoscmd.EncodingConfig(simapp.MakeTestEncodingConfig()), simapp.EmptyAppOptions{}, func(bapp *baseapp.BaseApp) {
		bapp.SetFauxMerkleMode()
	})
}
