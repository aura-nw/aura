package main

import (
	"github.com/aura-nw/aura/cmd/aurad/cmd"
	"os"

	"github.com/aura-nw/aura/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
