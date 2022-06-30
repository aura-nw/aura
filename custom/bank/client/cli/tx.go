package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

var FlagSplit = "split"

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bank transaction subcommands wrap",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		cli.NewSendTxCmd(),
		NewMultiSendTxCmd(),
	)

	return txCmd
}

// NewMultiSendTxCmd returns a CLI command handler for creating a MsgMultiSend transaction.
// For a better UX this command is limited to send funds from one account to two or more accounts.
func NewMultiSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send [from_key_or_address] [to_address_1, to_address_2, ...] [amount]",
		Short: "Send funds from one account to two or more accounts.",
		Long: `Send funds from one account to two or more accounts.
By default, sends the [amount] to each address of the list.
Using the '--split' flag, the [amount] is split equally between the addresses.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address].
When using '--dry-run' a key name cannot be used, only a bech32 address.
`,
		Args: cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[len(args)-1])
			if err != nil {
				return err
			}

			if coins.IsZero() {
				return fmt.Errorf("must send positive amount")
			}

			split, err := cmd.Flags().GetBool(FlagSplit)
			if err != nil {
				return err
			}

			totalAddrs := sdk.NewInt(int64(len(args) - 2))
			// coins to be received by the addresses
			sendCoins := coins
			if split {
				sendCoins = quoInt(coins, totalAddrs)
			}

			var output []types.Output
			for _, arg := range args[1 : len(args)-1] {
				toAddr, err := sdk.AccAddressFromBech32(arg)
				if err != nil {
					return err
				}

				output = append(output, types.NewOutput(toAddr, sendCoins))
			}

			// amount to be send from the from address
			var amount sdk.Coins
			if split {
				// user input: 1000stake to send to 3 addresses
				// actual: 333stake to each address (=> 999stake actually sent)
				amount = mulInt(sendCoins, totalAddrs)
			} else {
				amount = mulInt(coins, totalAddrs)
			}

			msg := types.NewMsgMultiSend([]types.Input{types.NewInput(clientCtx.FromAddress, amount)}, output)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool(FlagSplit, false, "Send the equally split token amount to each address")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// quoInt performs the scalar division of coins with a `divisor`
// All coins are divided by x and trucated.
// e.g.
// {2A, 30B} / 2 = {1A, 15B}
// {2A} / 2 = {1A}
// {4A} / {8A} = {0A}
// {2A} / 0 = panics
// Note, if IsValid was true on Coins, IsValid stays true,
// unless the `divisor` is greater than the smallest coin amount.
func quoInt(coins sdk.Coins, x sdk.Int) sdk.Coins {
	coins, ok := safeQuoInt(coins, x)
	if !ok {
		panic("dividing by zero is an invalid operation on coins")
	}

	return coins
}

func safeQuoInt(coins sdk.Coins, x sdk.Int) (sdk.Coins, bool) {
	if x.IsZero() {
		return nil, false
	}

	var res sdk.Coins
	for _, coin := range coins {
		coin := coin
		res = append(res, newCoin(coin.Denom, coin.Amount.Quo(x)))
	}

	return res, true
}

// mulInt performs the scalar multiplication of coins with a `multiplier`
// All coins are multipled by x
// e.g.
// {2A, 3B} * 2 = {4A, 6B}
// {2A} * 0 panics
// Note, if IsValid was true on Coins, IsValid stays true.
func mulInt(coins sdk.Coins, x sdk.Int) sdk.Coins {
	coins, ok := safeMulInt(coins, x)
	if !ok {
		panic("multiplying by zero is an invalid operation on coins")
	}

	return coins
}

// safeMulInt performs the same arithmetic as MulInt but returns false
// if the `multiplier` is zero because it makes IsValid return false.
func safeMulInt(coins sdk.Coins, x sdk.Int) (sdk.Coins, bool) {
	if x.IsZero() {
		return nil, false
	}

	res := make(sdk.Coins, len(coins))
	for i, coin := range coins {
		coin := coin
		res[i] = newCoin(coin.Denom, coin.Amount.Mul(x))
	}

	return res, true
}

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative or if the denomination is invalid.
func newCoin(denom string, amount sdk.Int) sdk.Coin {
	coin := sdk.Coin{
		Denom:  denom,
		Amount: amount,
	}

	if err := coin.Validate(); err != nil {
		panic(err)
	}

	return coin
}
