package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

const DefaultSmartAccountId uint64 = 0

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:         DefaultParams(),
		SmartAccountId: DefaultSmartAccountId,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
