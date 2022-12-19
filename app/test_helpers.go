package app

import (
	"encoding/json"

	"github.com/cosmos/ibc-go/v3/testing/simapp"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var defaultGenesisBz []byte

func getDefaultGenesisStateBytes() []byte {
	if len(defaultGenesisBz) == 0 {
		encConfig := MakeEncodingConfig()
		genesisState := NewDefaultGenesisState(encConfig.Codec)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		defaultGenesisBz = stateBytes
	}
	return defaultGenesisBz
}

func Setup(isCheckTx bool) *App {
	db := db.NewMemDB()
	app := New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, cosmoscmd.MakeEncodingConfig(ModuleBasics), simapp.EmptyAppOptions{})

	if !isCheckTx {
		stateBytes := getDefaultGenesisStateBytes()

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app.(*App)
}
