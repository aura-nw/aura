package main

import (
	"os"

	"github.com/aura-nw/aura/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/tendermint/spm/cosmoscmd"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
	)

	//testnet cmd
	rootCmd.AddCommand(
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
