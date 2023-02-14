package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aura-nw/aura/app"
	"github.com/aura-nw/aura/database"
	"github.com/cosmos/cosmos-sdk/simapp"
	junodb "github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/logging"
)

func TestLoadIndexerConfig(t *testing.T) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome := filepath.Join(userHomeDir, ".aurasql")
	indexerConfig, err := app.LoadIndexerConfig(DefaultNodeHome)
	fmt.Println(indexerConfig, err)
}

func TestIndexerConnection(t *testing.T) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	codec := simapp.MakeTestEncodingConfig()
	DefaultNodeHome := filepath.Join(userHomeDir, ".aurasql")
	indexerConfig, err := app.LoadIndexerConfig(DefaultNodeHome)
	db, err := database.Builder(junodb.NewContext(*indexerConfig, &codec, logging.DefaultLogger()))
	bigDipperDb, ok := (db).(*database.Db)
	if !ok {
		t.Fatal(`Error 1`)
	}
	// Delete the public schema
	_, err = bigDipperDb.Sqlx.Exec(`DROP SCHEMA public CASCADE;`)
	if err != nil {
		t.Fatal("Couldn't connect to psql")
	}
}
