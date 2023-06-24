package keeper

import (
	"fmt"
	"strconv"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InstantiateSmartAccount(
	ctx sdk.Context,
	keeper Keeper,
	wasmKeepper *wasmkeeper.PermissionedKeeper,
	msg *types.MsgActivateAccount,
) (sdk.AccAddress, []byte, cryptotypes.PubKey, error) {

	admin, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(types.ErrAddressFromBech32, err)
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(types.ErrAddressFromBech32, err)
	}

	pubKey, err := types.PubKeyDecode(msg.PubKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// generate salt using owner, codeId, initMsg and pubkey
	// make sure the contract that is initiated will have the address according to the pre-configured configuration
	salt, err := types.GenerateSalt(
		msg.Owner,
		msg.CodeID,
		msg.InitMsg,
		pubKey.Bytes(),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// instantiate smartcontract by contract with code_id
	address, data, iErr := wasmKeepper.Instantiate2(
		ctx,
		msg.CodeID,
		owner,       // owner
		admin,       // admin
		msg.InitMsg, // message
		fmt.Sprintf("%s/%d", types.ModuleName, keeper.GetAndIncrementNextAccountID(ctx)), // label
		sdk.NewCoins(), // empty funds
		salt,           // salt
		true,
	)
	if iErr != nil {
		return nil, nil, nil, fmt.Errorf(types.ErrBadInstantiateMsg, iErr.Error())
	}

	contractAddrStr := address.String()

	// make sure the new contract has the same address as predicted
	if contractAddrStr != msg.AccountAddress {
		return nil, nil, nil, fmt.Errorf(types.ErrBadInstantiateMsg, "wrong predicted address")
	}

	ctx.Logger().Info(
		"smart account created",
		types.AttributeKeyCreator, msg.AccountAddress,
		types.AttributeKeyCodeID, msg.CodeID,
		types.AttributeKeyContractAddr, contractAddrStr,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAccountRegistered,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.AccountAddress),
			sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(msg.CodeID, 10)),
			sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddrStr),
		),
	)

	return address, data, pubKey, nil
}

func IsWhitelistCodeID(ctx sdk.Context, keeper Keeper, codeID uint64) bool {
	params := keeper.GetParams(ctx)
	if params.WhitelistCodeID == nil {
		return false
	}

	// code_id must be in whitelist and has activated status
	for _, codeIDAllowed := range params.WhitelistCodeID {
		if codeID == codeIDAllowed.CodeID {
			if codeIDAllowed.Status {
				return true
			} else {
				return false
			}
		}
	}

	return false
}
