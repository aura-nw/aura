package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-key [address] [pub-key]",
		Short: "Update a smart account public key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bz, err := hex.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf(types.ErrBadPublicKey, err.Error())
			}

			// secp25k61 public key
			pubKey := secp256k1.PubKey{Key: nil}
			keyErr := pubKey.UnmarshalAmino(bz)
			if keyErr != nil {
				return fmt.Errorf(types.ErrBadPublicKey, keyErr.Error())
			}

			msg := &types.MsgUpdateKey{
				Creator: clientCtx.GetFromAddress().String(),
				Address: args[0],
				PubKey:  pubKey,
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
