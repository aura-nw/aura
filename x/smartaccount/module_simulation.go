package smartaccount

import (
	"math/rand"

	smartaccountsimulation "github.com/aura-nw/aura/x/smartaccount/simulation"
	"github.com/aura-nw/aura/x/smartaccount/types"
	typesv1 "github.com/aura-nw/aura/x/smartaccount/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = smartaccountsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgRecover = "op_weight_msg_recover"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRecover int = 100

	opWeightMsgActivateAccount = "op_weight_msg_activate_account"
	// TODO: Determine the simulation weight value
	defaultWeightMsgActivateAccount int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	smartaccountGenesis := typesv1.GenesisState{
		Params: typesv1.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&smartaccountGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRecover int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRecover, &weightMsgRecover, nil,
		func(_ *rand.Rand) {
			weightMsgRecover = defaultWeightMsgRecover
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRecover,
		smartaccountsimulation.SimulateMsgRecover(am.contractKeeper, am.accountKeeper, am.keeper),
	))

	var weightMsgActivateAccount int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgActivateAccount, &weightMsgActivateAccount, nil,
		func(_ *rand.Rand) {
			weightMsgActivateAccount = defaultWeightMsgActivateAccount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgActivateAccount,
		smartaccountsimulation.SimulateMsgActivateAccount(am.contractKeeper, am.accountKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
