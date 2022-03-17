package main

import (
	"os"

	"github.com/aura-nw/aura/app"
	"github.com/aura-nw/aura/cmd/aurad/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	// "github.com/tendermint/starport/starport/pkg/cosmoscmd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/spm/cosmoscmd"
)

func main() {
	// cmdOptions := cmd.GetWasmCmdOptions()
	// cmdOptions = append(cmdOptions, cosmoscmd.AddSubCmd(tmcmds.RollbackStateCmd))
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
		// cmdOptions...,
	)

	//testnet cmd
	rootCmd.AddCommand(
		cmd.TestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
	)

	// generate genesis vesting accounts cmd
	rootCmd.AddCommand(
		cmd.AddGenesisVestingAccountCmd(app.DefaultNodeHome),
	)

	rootCmd.AddCommand(
		cmd.AddGenesisWasmMsgCmd(app.DefaultNodeHome),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
