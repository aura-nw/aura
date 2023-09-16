package auranw

import (
	"encoding/base64"

	errorsmod "cosmossdk.io/errors"
	smartaccounttypes "github.com/aura-nw/aura/x/smartaccount/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgRecover = "recover"

var (
	_ sdk.Msg                            = &MsgRecover{}
	_ codectypes.UnpackInterfacesMessage = (*MsgRecover)(nil)
)

func (msg *MsgRecover) Route() string {
	return smartaccounttypes.RouterKey
}

func (msg *MsgRecover) Type() string {
	return TypeMsgRecover
}

func (msg *MsgRecover) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRecover) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRecover) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(smartaccounttypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(smartaccounttypes.ErrInvalidAddress, "invalid smart account address (%s)", err)
	}

	_, err = PubKeyDecode(msg.PubKey)
	if err != nil {
		return err
	}

	_, err = base64.StdEncoding.DecodeString(msg.Credentials)
	if err != nil {
		return errorsmod.Wrapf(smartaccounttypes.ErrInvalidCredentials, err.Error())
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgRecover) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.PubKey == nil {
		return nil
	}
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.PubKey, &pubKey)
}
