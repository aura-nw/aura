package main

import (
	"os"

	"github.com/aura-nw/aura/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/tendermint/spm/cosmoscmd"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmcmds "github.com/tendermint/tendermint/cmd/tendermint/commands"
)

func main() {
	cmdOptions := GetWasmCmdOptions()
	cmdOptions = append(cmdOptions, cosmoscmd.AddSubCmd(tmcmds.RollbackStateCmd))
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
		cmdOptions...,
	)

	//testnet cmd
	rootCmd.AddCommand(
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
