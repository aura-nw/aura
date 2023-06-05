package types

import (
	"errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	_ authtypes.AccountI                 = (*SmartAccount)(nil)
	_ codectypes.UnpackInterfacesMessage = (*SmartAccount)(nil)
)

// ------------------------------ SmartAccount ------------------------------

func NewSmartAccount(address string, accountNum, seq uint64) *SmartAccount {
	return &SmartAccount{
		Address:       address,
		AccountNumber: accountNum,
		Sequence:      seq,
	}
}

func NewSmartAccountFromAccount(acc authtypes.AccountI) *SmartAccount {
	return NewSmartAccount(acc.GetAddress().String(), acc.GetAccountNumber(), acc.GetSequence())
}

func (acc *SmartAccount) GetAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(acc.Address)
	return addr
}

func (acc *SmartAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override SmartAccount address")
	}

	acc.Address = addr.String()

	return nil
}

func (acc *SmartAccount) GetPubKey() cryptotypes.PubKey {
	if acc.PubKey == nil {
		return nil
	}

	content, ok := acc.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil
	}
	return content
}

func (acc *SmartAccount) SetPubKey(pubKey cryptotypes.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
		return nil
	}
	any, err := codectypes.NewAnyWithValue(pubKey)

	if err == nil {
		acc.PubKey = any
	}
	return err
}

func (acc *SmartAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

func (acc *SmartAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber

	return nil
}

func (acc *SmartAccount) GetSequence() uint64 {
	return acc.Sequence
}

func (acc *SmartAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (acc SmartAccount) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if acc.PubKey == nil {
		return nil
	}
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(acc.PubKey, &pubKey)
}
