package ibc_middleware

import (
	// external libraries
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	// ibc-go
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"
)

var _ porttypes.ICS4Wrapper = &ICS4Middleware{}

type ICS4Middleware struct {
	channel porttypes.ICS4Wrapper
	wasm    *WasmHooks
}

func NewICS4Middleware(channel porttypes.ICS4Wrapper, wasm *WasmHooks) ICS4Middleware {
	return ICS4Middleware{
		channel: channel,
		wasm:    wasm,
	}
}

// SendPacket implements the ICS4 Wrapper interface
func (i ICS4Middleware) SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {

	// Override SendPacket method for processing the callback contract
	return i.wasm.SendPacketOverride(i, ctx, channelCap, packet)
}

// WriteAcknowledgement implements the ICS4 Wrapper interface
func (i ICS4Middleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, ack ibcexported.Acknowledgement) error {

	return i.channel.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

// GetAppVersion implements the ICS4 Wrapper interface
func (i ICS4Middleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {

	version, err := i.channel.GetAppVersion(ctx, portID, channelID)

	return version, err
}
