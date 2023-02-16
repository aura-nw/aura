package hasura

import (
	"github.com/aura-nw/aura/hasura/utils"
	modulestypes "github.com/aura-nw/aura/x/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/node"
	nodeconfig "github.com/forbole/juno/v4/node/config"
)

const (
	ModuleName = "actions"
)

type HasuraService struct {
	cfg     *Config
	node    node.Node
	sources *modulestypes.Sources
}

func NewHasuraService(homepath string, encodingConfig *params.EncodingConfig) *HasuraService {
	actionsCfg, err := LoadHasuraConfig(homepath)
	if err != nil {
		panic(err)
	}
	nodeCfg := nodeconfig.NewConfig(nodeconfig.TypeRemote, actionsCfg.Node)

	// Build the node
	junoNode, err := utils.BuildNode(nodeCfg, encodingConfig)
	if err != nil {
		panic(err)
	}

	// Build the sources
	sources, err := modulestypes.BuildSources(nodeCfg, encodingConfig)
	if err != nil {
		panic(err)
	}

	return &HasuraService{
		cfg:     actionsCfg,
		node:    junoNode,
		sources: sources,
	}
}
