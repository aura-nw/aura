package keeper

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	SmartAccountLabel = "smart account"
)

func instantiateSmartAccount(ctx sdk.Context, wasmKeepper *wasmkeeper.PermissionedKeeper, msg *types.MsgCreateAccount) (sdk.AccAddress, error) {

	creator, aMsg := sdk.AccAddressFromBech32(msg.Creator)
	if aMsg != nil {
		return nil, fmt.Errorf(types.ErrAddressFromBech32, aMsg)
	}

	// instantiate smartcontract by code id
	address, _, iErr := wasmKeepper.Instantiate(
		ctx,
		uint64(msg.CodeId),
		creator,                 // owner
		creator,                 // admin
		[]byte(msg.InitMessage), // message
		SmartAccountLabel,       // label
		msg.Funds,               // funds
	)
	if iErr != nil {
		return nil, fmt.Errorf(types.ErrBadInstantiateMsg, iErr.Error())
	}

	return address, nil
}