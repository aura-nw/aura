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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

const flagFunds = "funds"

func CmdCreateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account [code-id] [init-msg] [public-key] [salt] --funds [coins,optional]",
		Short: "Create a smart account",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
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

			bz, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf(types.ErrBadPublicKey, err.Error())
			}

			// secp25k61 public key
			pubKey := secp256k1.PubKey{Key: nil}
			keyErr := pubKey.UnmarshalAmino(bz)
			if keyErr != nil {
				return fmt.Errorf(types.ErrBadPublicKey, keyErr.Error())
			}

			msg := &types.MsgCreateAccount{
				Creator: clientCtx.GetFromAddress().String(),
				CodeID:  codeID,
				InitMsg: []byte(args[1]),
				PubKey:  pubKey,
				Funds:   funds,
				Salt:    []byte(args[3]),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
		SilenceUsage: true,
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(flagFunds, "", "Coins to send to the account during instantiation")

	return cmd
}
