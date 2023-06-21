package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgActivateAccount = "activate_account"

var _ sdk.Msg = &MsgActivateAccount{}

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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
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

	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	return nil
}
