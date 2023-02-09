package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aura-nw/aura/app"
)

func TestLoadIndexerConfig(t *testing.T) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome := filepath.Join(userHomeDir, ".aura")
	indexerConfig, err := app.LoadIndexerConfig(DefaultNodeHome)
	fmt.Println(indexerConfig, err)
}
