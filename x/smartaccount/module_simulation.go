package smartaccount

import (
	"math/rand"

	smartaccountsimulation "github.com/aura-nw/aura/x/smartaccount/simulation"
	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = smartaccountsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateAccount = "op_weight_msg_create_account"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAccount int = 100

	opWeightMsgUpdateKey = "op_weight_msg_update_key"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateKey int = 100

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
	smartaccountGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&smartaccountGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateAccount int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateAccount, &weightMsgCreateAccount, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAccount = defaultWeightMsgCreateAccount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAccount,
		smartaccountsimulation.SimulateMsgCreateAccount(am.contractKeeper, am.accountKeeper, am.keeper),
	))

	var weightMsgUpdateKey int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateKey, &weightMsgUpdateKey, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateKey = defaultWeightMsgUpdateKey
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateKey,
		smartaccountsimulation.SimulateMsgUpdateKey(am.contractKeeper, am.accountKeeper, am.keeper),
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
