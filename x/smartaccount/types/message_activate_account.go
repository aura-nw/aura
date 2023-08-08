package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgActivateAccount = "activate_account"

var (
	_ sdk.Msg                            = &MsgActivateAccount{}
	_ codectypes.UnpackInterfacesMessage = (*MsgActivateAccount)(nil)
)

func (msg *MsgActivateAccount) Route() string {
	return RouterKey
}

func (msg *MsgActivateAccount) Type() string {
	return TypeMsgActivateAccount
}

func (msg *MsgActivateAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgActivateAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgActivateAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid smart account address (%s)", err)
	}

	if len(msg.Salt) > 64 {
		return sdkerrors.ErrInvalidRequest.Wrap("length of salt too long")
	}

	_, err = PubKeyDecode(msg.PubKey)
	if err != nil {
		return err
	}

	if msg.CodeID == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("code id cannot be zero")
	}

	if err := msg.InitMsg.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid init msg: %s", err.Error())
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgActivateAccount) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.PubKey == nil {
		return nil
	}
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.PubKey, &pubKey)
}
