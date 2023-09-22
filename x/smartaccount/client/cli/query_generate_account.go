package cli

import (
	"strconv"

	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGenerateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-account [code_id:uint64] [salt:string] [init_msg:string] [pub_key]",
		Short: "Query GenerateAccount",
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

			pubKey, err := typesv1.PubKeyToAny(clientCtx.Codec, []byte(args[3]))
			if err != nil {
				return err
			}

			params := &typesv1.QueryGenerateAccountRequest{
				CodeID:  codeID,
				Salt:    []byte(args[1]),
				InitMsg: []byte(args[2]),
				PubKey:  pubKey,
			}

			queryClient := typesv1.NewQueryClient(clientCtx)

			res, err := queryClient.GenerateAccount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
