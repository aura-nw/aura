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
		Use:   "activate-account [account-address:str] [owner:str] [code_id:uint64] [pub_key] [init_msg:str] --funds [coins,optional]",
		Short: "Broadcast message ActivateAccount",
		Args:  cobra.ExactArgs(5),
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

			pubKey, err := types.PubKeyToAny(clientCtx.Codec, []byte(args[3]))
			if err != nil {
				return err
			}

			msg := &types.MsgActivateAccount{
				AccountAddress: args[0],
				Owner:          args[1],
				CodeID:         codeID,
				PubKey:         pubKey,
				InitMsg:        []byte(args[4]),
				Funds:          funds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(flagFunds, "", "Coins to send to the account during instantiation")

	return cmd
}
