package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgActivateAccount = "activate_account"

var _ sdk.Msg = &MsgActivateAccount{}

func NewMsgActivateAccount(creator string) *MsgActivateAccount {
	return &MsgActivateAccount{
		Creator: creator,
	}
}

func (msg *MsgActivateAccount) Route() string {
	return RouterKey
}

func (msg *MsgActivateAccount) Type() string {
	return TypeMsgActivateAccount
}

func (msg *MsgActivateAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
