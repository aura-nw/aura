package feegrant

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/gogo/protobuf/proto"
)

const (
	gasCostPerIteration = uint64(10)
)

var (
	_ feegrant.FeeAllowanceI        = (*AllowedContractAllowance)(nil)
	_ types.UnpackInterfacesMessage = (*AllowedContractAllowance)(nil)
)

// NewAllowedContractAllowance creates new filtered fee allowance.
func NewAllowedContractAllowance(allowance feegrant.FeeAllowanceI, allowedAddress []string) (*AllowedContractAllowance, error) {
	msg, ok := allowance.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &AllowedContractAllowance{
		Allowance:      any,
		AllowedAddress: allowedAddress,
	}, nil
}

func (a *AllowedContractAllowance) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var allowance feegrant.FeeAllowanceI
	return unpacker.UnpackAny(a.Allowance, &allowance)
}

func (a *AllowedContractAllowance) Accept(ctx sdk.Context, fee sdk.Coins, msgs []sdk.Msg) (remove bool, err error) {
	if !a.allContractAllowed(ctx, msgs) {
		return false, sdkerrors.Wrap(ErrAddressNotAllowed, "address does not exist in allowed addresses")
	}

	allowance, err := a.GetAllowance()
	if err != nil {
		return false, err
	}

	return allowance.Accept(ctx, fee, msgs)
}

func (a *AllowedContractAllowance) ValidateBasic() error {
	if a.Allowance == nil {
		return sdkerrors.Wrap(feegrant.ErrNoAllowance, "allowance should not be empty")
	}
	if len(a.AllowedAddress) == 0 {
		return sdkerrors.Wrap(feegrant.ErrNoMessages, "allowed address shouldn't be empty")
	}

	allowance, err := a.GetAllowance()
	if err != nil {
		return err
	}

	return allowance.ValidateBasic()
}

// GetAllowance returns allowed fee allowance.
func (a *AllowedContractAllowance) GetAllowance() (feegrant.FeeAllowanceI, error) {
	allowance, ok := a.Allowance.GetCachedValue().(feegrant.FeeAllowanceI)
	if !ok {
		return nil, sdkerrors.Wrap(feegrant.ErrNoAllowance, "failed to get allowance")
	}

	return allowance, nil
}

func (a *AllowedContractAllowance) allContractAllowed(ctx sdk.Context, msgs []sdk.Msg) bool {
	addressesToMap := a.allowedAddressesToMap(ctx)

	for _, msg := range msgs {
		ctx.GasMeter().ConsumeGas(gasCostPerIteration, "check contract address")
		switch msg := msg.(type) {
		case *wasmtypes.MsgExecuteContract:
			if !addressesToMap[msg.Contract] {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func (a *AllowedContractAllowance) allowedAddressesToMap(ctx sdk.Context) map[string]bool {
	addressesMap := make(map[string]bool, len(a.AllowedAddress))
	for _, msg := range a.AllowedAddress {
		ctx.GasMeter().ConsumeGas(gasCostPerIteration, "check contract address")
		addressesMap[msg] = true
	}

	return addressesMap
}
