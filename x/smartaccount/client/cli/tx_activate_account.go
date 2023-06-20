package cli

import (
	"fmt"
	"strconv"

	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdActivateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-account [creator] [owner] [code_id] [init_msg] --funds [coins,optional]",
		Short: "Broadcast message ActivateAccount",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			fundsStr, err := cmd.Flags().GetString(flagFunds)
			if err != nil {
				return fmt.Errorf("funds: %s", err)
			}

			funds, err := sdk.ParseCoinsNormalized(fundsStr)
			if err != nil {
				return fmt.Errorf("funds: %s", err)
			}

			msg := &types.MsgActivateAccount{
				AccountAddress: args[0],
				Owner:          args[1],
				CodeID:         codeID,
				InitMsg:        []byte(args[3]),
				Funds:          funds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
