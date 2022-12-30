package distribution

import (
	"encoding/json"
	customkeeper "github.com/aura-nw/aura/x/distribution/keeper"
	"github.com/aura-nw/aura/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func (a AppModuleBasic) Name() string {
	return disttypes.ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
}

func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

}

func (a AppModuleBasic) DefaultGenesis(jsonCodec codec.JSONCodec) json.RawMessage {
	return jsonCodec.MustMarshalJSON(disttypes.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(jsonCodec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

type AppModule struct {
	distribution.AppModule

	keeper customkeeper.Keeper
}

func NewAppModule(cdc codec.Codec, customKeeper customkeeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, stackingKeeper types.StakingKeeper) AppModule {
	return AppModule{
		AppModule: distribution.NewAppModule(cdc, customKeeper.Keeper, accountKeeper, bankKeeper, stackingKeeper),
		keeper:    customKeeper,
	}
}
