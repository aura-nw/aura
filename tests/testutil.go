package tests

import (
	"encoding/json"
	"github.com/aura-nw/aura/app"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var defaultGenesisBz []byte

func getDefaultGenesisStateBytes() []byte {
	if len(defaultGenesisBz) == 0 {
		encConfig := app.MakeEncodingConfig()
		genesisState := app.NewDefaultGenesisState(encConfig.Codec)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		defaultGenesisBz = stateBytes
	}
	return defaultGenesisBz
}

func Setup(isCheckTx bool) *app.App {
	db := db.NewMemDB()
	appObj := app.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		cosmoscmd.MakeEncodingConfig(app.ModuleBasics),
		simapp.EmptyAppOptions{})

	if !isCheckTx {
		stateBytes := getDefaultGenesisStateBytes()

		appObj.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return appObj.(*app.App)
}
