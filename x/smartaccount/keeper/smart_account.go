package keeper

import (
	"encoding/hex"
	"fmt"
	"strconv"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/aura-nw/aura/x/smartaccount/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InstantiateSmartAccount(ctx sdk.Context, keeper Keeper, wasmKeepper *wasmkeeper.PermissionedKeeper, msg *types.MsgCreateAccount) (sdk.AccAddress, []byte, error) {

	creator, aMsg := sdk.AccAddressFromBech32(msg.Creator)
	if aMsg != nil {
		return nil, nil, fmt.Errorf(types.ErrAddressFromBech32, aMsg)
	}

	// instantiate smartcontract by code id
	address, data, iErr := wasmKeepper.Instantiate2(
		ctx,
		msg.CodeID,
		creator,     // owner
		creator,     // admin
		msg.InitMsg, // message
		fmt.Sprintf("%s/%d", types.ModuleName, keeper.GetAndIncrementNextAccountID(ctx)), // label
		msg.Funds, // funds
		msg.Salt,  // salt
		true,
	)
	if iErr != nil {
		return nil, nil, fmt.Errorf(types.ErrBadInstantiateMsg, iErr.Error())
	}

	// set the contract's admin to itself
	if err := wasmKeepper.UpdateContractAdmin(ctx, address, creator, address); err != nil {
		return nil, nil, err
	}

	contractAddrStr := address.String()

	ctx.Logger().Info(
		"smart account created",
		types.AttributeKeyCreator, msg.Creator,
		types.AttributeKeyCodeID, msg.CodeID,
		types.AttributeKeyContractAddr, contractAddrStr,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAccountRegistered,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(msg.CodeID, 10)),
			sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddrStr),
		),
	)

	return address, data, nil
}

func PubKeyDecode(raw string) (*secp256k1.PubKey, error) {
	bz, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf(types.ErrBadPublicKey, err.Error())
	}

	// secp25k61 public key
	pubKey := &secp256k1.PubKey{Key: nil}
	keyErr := pubKey.UnmarshalAmino(bz)
	if keyErr != nil {
		return nil, fmt.Errorf(types.ErrBadPublicKey, keyErr.Error())
	}

	return pubKey, nil
}
