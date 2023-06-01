package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const TypeMsgCreateAccount = "create_account"

var _ sdk.Msg = &MsgCreateAccount{}

func (msg *MsgCreateAccount) Route() string {
	return RouterKey
}

func (msg *MsgCreateAccount) Type() string {
	return TypeMsgCreateAccount
}

func (msg *MsgCreateAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
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

	if err := wasmtypes.ValidateSalt(msg.Salt); err != nil {
		return err
	}

	return nil
}
