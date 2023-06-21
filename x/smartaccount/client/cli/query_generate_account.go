package cli

import (
	"strconv"

	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGenerateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-account [code_id:uint64] [owner:string] [init_msg:string] [pub_key:hex]",
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

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGenerateAccountRequest{
				CodeID:  codeID,
				Owner:   args[1],
				InitMsg: []byte(args[2]),
				PubKey:  args[3],
			}

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
