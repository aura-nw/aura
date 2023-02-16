package hasura_test

import (
	"path/filepath"
	"testing"

	"github.com/aura-nw/aura/hasura"
	"github.com/cosmos/cosmos-sdk/simapp"
)

func TestStartHasura(t *testing.T) {
	codec := simapp.MakeTestEncodingConfig()
	pathTest := filepath.Join("../")
	err := hasura.NewHasuraService(pathTest, &codec).Start()
	if err != nil {
		panic(err)
	}
}
