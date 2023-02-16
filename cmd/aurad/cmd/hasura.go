package cmd

import (
	"github.com/aura-nw/aura/hasura"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cobra"
)

var (
	flagHasuraHost = "hasura-host"
	flagHasuraPort = "hasura-port"
	flagHasuraRpc  = "hasura-rpc"
	flagHasuraGrpc = "hasura-grpc"
)

// get cmd to initialize all files for tendermint testnet and application
func Hasura() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hasura",
		Short: "Start hasura entity",
		Long: `hasura will create a hasura handler services, you need to create hasura.yaml in homepath/config
		( example: ~/.aura/config ). It reads host, port to run Mux handler, node config for rpc which it connects
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cdc := codec.NewLegacyAmino()
			interfaceRegistry := types.NewInterfaceRegistry()
			marshaler := codec.NewProtoCodec(interfaceRegistry)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			return hasura.NewHasuraService(clientCtx.HomeDir, &params.EncodingConfig{InterfaceRegistry: interfaceRegistry,
				Marshaler: marshaler,
				TxConfig:  tx.NewTxConfig(marshaler, tx.DefaultSignModes),
				Amino:     cdc}).Start()
		},
	}

	cmd.Flags().String(flagHasuraHost, "127.0.0.1", "Hasura handlers host")
	cmd.Flags().Int(flagHasuraPort, 3000, "Hasura handlers port")
	cmd.Flags().String(flagHasuraRpc, "http://localhost:26657", "URL Rpc which Hasura handlers connect to")
	cmd.Flags().String(flagHasuraGrpc, "http://localhost:9090", "URL Grpc which Hasura handlers connect to")

	return cmd
}
