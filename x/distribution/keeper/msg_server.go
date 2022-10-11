package keeper

import (
	"context"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	orgdistrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type msgServer struct {
	types.MsgServer
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		MsgServer: orgdistrkeeper.NewMsgServerImpl(keeper.Keeper),
		Keeper:    keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) SetWithdrawAddress(ctx context.Context, address *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	return m.MsgServer.SetWithdrawAddress(ctx, address)
}

func (m msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	amount, err := m.WithdrawDelegationRewards(ctx, delegatorAddress, valAddr)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	)
	return &types.MsgWithdrawDelegatorRewardResponse{}, nil
}

func (m msgServer) WithdrawValidatorCommission(ctx context.Context, commission *types.MsgWithdrawValidatorCommission) (*types.MsgWithdrawValidatorCommissionResponse, error) {
	return m.MsgServer.WithdrawValidatorCommission(ctx, commission)
}

func (m msgServer) FundCommunityPool(ctx context.Context, pool *types.MsgFundCommunityPool) (*types.MsgFundCommunityPoolResponse, error) {
	return m.MsgServer.FundCommunityPool(ctx, pool)
}
