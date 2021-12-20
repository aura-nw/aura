package keeper

import (
	"fmt"
	wasmvm "github.com/CosmWasm/wasmvm"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"path/filepath"

	"github.com/tendermint/tendermint/libs/log"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/aura-nw/aura/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// contractMemoryLimit is the memory limit of each contract execution (in MiB)
// constant value so all nodes run with the same limit.
const contractMemoryLimit = 32

// Option is an extension point to instantiate keeper with non default values
type Option interface {
	apply(*Keeper)
}

// WasmVMQueryHandler is an extension point for custom query handler implementations
type WasmVMQueryHandler interface {
	// HandleQuery executes the requested query
	HandleQuery(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error)
}

type CoinTransferrer interface {
	// TransferCoins sends the coin amounts from the source to the destination with rules applied.
	TransferCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// WasmVMResponseHandler is an extension point to handles the response data returned by a contract call.
type WasmVMResponseHandler interface {
	// Handle processes the data returned by a contract invocation.
	Handle(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		ibcPort string,
		messages []wasmvmtypes.SubMsg,
		origRspData []byte,
	) ([]byte, error)
}

type (
	Keeper struct {
		cdc                   codec.BinaryCodec
		storeKey              sdk.StoreKey
		accountKeeper         types.AccountKeeper
		bank                  CoinTransferrer
		portKeeper            types.PortKeeper
		capabilityKeeper      types.CapabilityKeeper
		wasmVM                types.WasmerEngine
		wasmVMQueryHandler    WasmVMQueryHandler
		wasmVMResponseHandler WasmVMResponseHandler
		messenger             Messenger
		// queryGasLimit is the max wasmvm gas that can be spent on executing a query with a contract
		queryGasLimit uint64
		paramSpace    paramtypes.Subspace
		gasRegister   GasRegister
	}
)

func (k Keeper) reply(ctx sdk.Context, contractAddress sdk.AccAddress, reply wasmvmtypes.Reply) ([]byte, error) {
	panic("implement me")
}

func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *types.ContractInfo {
	panic("implement me")
}

func (k Keeper) QueryRaw(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte {
	panic("implement me")
}

func (k Keeper) QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	panic("implement me")
}

func (k Keeper) IsPinnedCode(ctx sdk.Context, codeID uint64) bool {
	panic("implement me")
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distKeeper types.DistributionKeeper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	capabilityKeeper types.CapabilityKeeper,
	portSource types.ICS20TransferPortSource,
	serviceRouter types.MsgServiceRouter,
	queryRouter GRPCQueryRouter,
	homeDir string,
	wasmConfig types.WasmConfig,
	supportedFeatures string,
	opts ...Option,

) Keeper {
	wasmer, err := wasmvm.NewVM(filepath.Join(homeDir, "wasm"), supportedFeatures, contractMemoryLimit, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	if err != nil {
		panic(err)
	}
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		wasmVM:           wasmer,
		accountKeeper:    accountKeeper,
		bank:             NewBankCoinTransferrer(bankKeeper),
		portKeeper:       portKeeper,
		capabilityKeeper: capabilityKeeper,
		messenger:        NewDefaultMessageHandler(serviceRouter, channelKeeper, capabilityKeeper, bankKeeper, cdc, portSource),
		queryGasLimit:    wasmConfig.SmartQueryGasLimit,
		paramSpace:       paramSpace,
		gasRegister:      NewDefaultWasmGasRegister(),
	}
	keeper.wasmVMQueryHandler = DefaultQueryPlugins(bankKeeper, stakingKeeper, distKeeper, channelKeeper, queryRouter, keeper)
	for _, o := range opts {
		o.apply(&keeper)
	}
	// not updateable, yet
	keeper.wasmVMResponseHandler = NewDefaultWasmVMContractResponseHandler(NewMessageDispatcher(keeper.messenger, keeper))
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return moduleLogger(ctx)
}

func moduleLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// BankCoinTransferrer replicates the cosmos-sdk behaviour as in
// https://github.com/cosmos/cosmos-sdk/blob/v0.41.4/x/bank/keeper/msg_server.go#L26
type BankCoinTransferrer struct {
	keeper types.BankKeeper
}

func NewBankCoinTransferrer(keeper types.BankKeeper) BankCoinTransferrer {
	return BankCoinTransferrer{
		keeper: keeper,
	}
}

// TransferCoins transfers coins from source to destination account when coin send was enabled for them and the recipient
// is not in the blocked address list.
func (c BankCoinTransferrer) TransferCoins(parentCtx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	em := sdk.NewEventManager()
	ctx := parentCtx.WithEventManager(em)
	if err := c.keeper.SendEnabledCoins(ctx, amt...); err != nil {
		return err
	}
	if c.keeper.BlockedAddr(fromAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
	}
	sdkerr := c.keeper.SendCoins(ctx, fromAddr, toAddr, amt)
	if sdkerr != nil {
		return sdkerr
	}
	for _, e := range em.Events() {
		if e.Type == sdk.EventTypeMessage { // skip messages as we talk to the keeper directly
			continue
		}
		parentCtx.EventManager().EmitEvent(e)
	}
	return nil
}

type msgDispatcher interface {
	DispatchSubmessages(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmvmtypes.SubMsg) ([]byte, error)
}

// DefaultWasmVMContractResponseHandler default implementation that first dispatches submessage then normal messages.
// The Submessage execution may include an success/failure response handling by the contract that can overwrite the
// original
type DefaultWasmVMContractResponseHandler struct {
	md msgDispatcher
}

func NewDefaultWasmVMContractResponseHandler(md msgDispatcher) *DefaultWasmVMContractResponseHandler {
	return &DefaultWasmVMContractResponseHandler{md: md}
}

// Handle processes the data returned by a contract invocation.
func (h DefaultWasmVMContractResponseHandler) Handle(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, messages []wasmvmtypes.SubMsg, origRspData []byte) ([]byte, error) {
	result := origRspData
	switch rsp, err := h.md.DispatchSubmessages(ctx, contractAddr, ibcPort, messages); {
	case err != nil:
		return nil, sdkerrors.Wrap(err, "submessages")
	case rsp != nil:
		result = rsp
	}
	return result, nil
}
