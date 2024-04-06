package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

// RegisterLegacyAminoCodec registers the necessary evmutil interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgConvertCoinToERC20{}, "evmutil/MsgConvertCoinToERC20")
	legacy.RegisterAminoMsg(cdc, &MsgConvertERC20ToCoin{}, "evmutil/MsgConvertERC20ToCoin")
	legacy.RegisterAminoMsg(cdc, &MsgConvertCosmosCoinToERC20{}, "evmutil/MsgConvertCosmosCoinToERC20")
	legacy.RegisterAminoMsg(cdc, &MsgConvertCosmosCoinFromERC20{}, "evmutil/MsgConvertCosmosCoinFromERC20")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConvertCoinToERC20{},
		&MsgConvertERC20ToCoin{},
		&MsgConvertCosmosCoinToERC20{},
		&MsgConvertCosmosCoinFromERC20{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}
