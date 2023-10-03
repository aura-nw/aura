package auranw

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const DefaultMaxGas = 2_000_000

// Parameter store keys
var (
	WhitelistCodeID = []byte("WhitelistCodeID")
	DisableMsgsList = []byte("DisableMsgsList")
	MaxGasExecute   = []byte("MaxGasExecute")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(whitelist []*CodeID, disableList []string, limitGas uint64) Params {
	return Params{
		WhitelistCodeID: whitelist,
		DisableMsgsList: disableList,
		MaxGasExecute:   limitGas,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	emptyWhitelistCodeID := make([]*CodeID, 0)
	emptyDisableMsgs := make([]string, 0)

	return NewParams(emptyWhitelistCodeID, emptyDisableMsgs, DefaultMaxGas)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// code_id whitelist indicates which contract can be initialized as smart account
		// using gov proposal for updates
		paramtypes.NewParamSetPair(WhitelistCodeID, &p.WhitelistCodeID, validateWhitelistCodeID),
		// list of diable messages for smartaccount
		paramtypes.NewParamSetPair(DisableMsgsList, &p.DisableMsgsList, validateDisableMsgsList),
		// max_gas_query limits the amount of gas that the validation query can use
		paramtypes.NewParamSetPair(MaxGasExecute, &p.MaxGasExecute, validateMaxGasExecute),
	}
}

func validateWhitelistCodeID(i interface{}) error {
	v, ok := i.([]*CodeID)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// not allowed duplicate code_id in whitelist
	visited := make(map[uint64]bool, 0)
	for _, codeID := range v {
		if visited[codeID.CodeID] {
			return fmt.Errorf("duplicate code_id %d in whitelist_code_id", codeID.CodeID)
		} else {
			visited[codeID.CodeID] = true
		}
	}

	return nil
}

func validateMaxGasExecute(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("zero max gas execute")
	}

	return nil
}

func validateDisableMsgsList(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// not allowed duplicate code_id in whitelist
	visited := make(map[string]bool, 0)
	for _, url := range v {
		if visited[url] {
			return fmt.Errorf("duplicate messages %s in disable_msgs_list", url)
		} else {
			visited[url] = true
		}
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	// validate whitelist_code_id param
	err := validateWhitelistCodeID(p.WhitelistCodeID)
	if err != nil {
		return err
	}

	//validate disable_msgs_list param
	err = validateDisableMsgsList(p.DisableMsgsList)
	if err != nil {
		return err
	}

	// validate max gas execute
	err = validateMaxGasExecute(p.MaxGasExecute)
	if err != nil {
		return err
	}

	return nil
}

func (p Params) IsAllowedCodeID(codeID uint64) bool {
	if p.WhitelistCodeID == nil {
		return false
	}

	// code_id must be in whitelist and has activated status
	for _, codeIDAllowed := range p.WhitelistCodeID {
		if codeID == codeIDAllowed.CodeID && codeIDAllowed.Status {
			return true
		}
	}

	return false
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
