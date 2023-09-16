package cli

import (
	"strconv"

	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRecover() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover [address:str] [pub-key] [credentials:base64]",
		Short: "Recover a smart account public key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pubKey, err := typesv1.PubKeyToAny(clientCtx.Codec, []byte(args[1]))
			if err != nil {
				return err
			}

			msg := &typesv1.MsgRecover{
				Creator:     clientCtx.GetFromAddress().String(),
				Address:     args[0],
				PubKey:      pubKey,
				Credentials: args[2],
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
