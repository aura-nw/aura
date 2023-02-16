package hasura_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/aura-nw/aura/hasura"
)

func TestLoadHasuraConfig(t *testing.T) {
	pathTest := filepath.Join("../")
	hasuraConfig, err := hasura.LoadHasuraConfig(pathTest)
	fmt.Println(hasuraConfig.Node.GRPC.Address)
	fmt.Println(err)
}
