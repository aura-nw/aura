package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	amount := sdk.NewIntFromUint64(1_000_000)
	fmt.Println(amount.ToDec())

	fmt.Println(amount.ToDec().QuoInt(sdk.NewIntFromUint64(50_000)))
}
