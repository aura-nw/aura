package cli

import (
    "strconv"
	
	 "github.com/spf13/cast"
	"github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/aura-nw/aura/x/smartaccount/types"
)

var _ = strconv.Itoa(0)

func CmdCreateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account [code-id] [init-message] [public-key]",
		Short: "Broadcast message CreateAccount",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
      		 argCodeId, err := cast.ToInt32E(args[0])
            		if err != nil {
                		return err
            		}
             argInitMessage := args[1]
             argPublicKey := args[2]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAccount(
				clientCtx.GetFromAddress().String(),
				argCodeId,
				argInitMessage,
				argPublicKey,
				
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

    return cmd
}