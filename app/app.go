package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	appparams "github.com/aura-nw/aura/app/params"
	v703 "github.com/aura-nw/aura/app/upgrades/v0.7.3"

	v500 "github.com/aura-nw/aura/app/upgrades/v0.5.0"
	v501 "github.com/aura-nw/aura/app/upgrades/v0.5.1"
	v600 "github.com/aura-nw/aura/app/upgrades/v0.6.0"
	v601 "github.com/aura-nw/aura/app/upgrades/v0.6.1"
	v700 "github.com/aura-nw/aura/app/upgrades/v0.7.0"
	v701 "github.com/aura-nw/aura/app/upgrades/v0.7.1"
	v702 "github.com/aura-nw/aura/app/upgrades/v0.7.2"

	"github.com/aura-nw/aura/app/internal"

	"github.com/aura-nw/aura/app/utils"
	custombank "github.com/aura-nw/aura/x/bank"
	custombankkeeper "github.com/aura-nw/aura/x/bank/keeper"
	customfeegrantmodule "github.com/aura-nw/aura/x/feegrant/module"
	custommint "github.com/aura-nw/aura/x/mint"
	custommintkeeper "github.com/aura-nw/aura/x/mint/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"

	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/aura-nw/aura/docs"

	auramodule "github.com/aura-nw/aura/x/aura"
	auramodulekeeper "github.com/aura-nw/aura/x/aura/keeper"
	auramoduletypes "github.com/aura-nw/aura/x/aura/types"

	samodule "github.com/aura-nw/aura/x/smartaccount"
	samodulekeeper "github.com/aura-nw/aura/x/smartaccount/keeper"
	samoduletypes "github.com/aura-nw/aura/x/smartaccount/types"

	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	ibc_hooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"

	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/prometheus/client_golang/prometheus"

	v0_3_0 "github.com/aura-nw/aura/app/upgrades/v0.3.0"
	v0_3_1 "github.com/aura-nw/aura/app/upgrades/v0.3.1"
	v0_3_2 "github.com/aura-nw/aura/app/upgrades/v0.3.2"
	v0_3_3 "github.com/aura-nw/aura/app/upgrades/v0.3.3"
	v0_4_0 "github.com/aura-nw/aura/app/upgrades/v0.4.0"
	v0_4_1 "github.com/aura-nw/aura/app/upgrades/v0.4.1"
	v0_4_2 "github.com/aura-nw/aura/app/upgrades/v0.4.2"
	v0_4_4 "github.com/aura-nw/aura/app/upgrades/v0.4.4"
	v0_4_5 "github.com/aura-nw/aura/app/upgrades/v0.4.5"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

const (
	AccountAddressPrefix = "aura"
	Name                 = "aura"
)

// this line is used by starport scaffolding # stargate/wasm/app/enabledProposals

var (
	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""

	EmptyWasmOpts []wasmkeeper.Option

	ChainID = ""
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasmtypes.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasmtypes.EnableAllProposals
		}
		return wasmtypes.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasmtypes.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)
	return govProposalHandlers
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasmkeeper.Option {
	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}
	return wasmOpts
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		customfeegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		consensus.AppModuleBasic{},
		ibc.AppModuleBasic{},
		// https://github.com/cosmos/ibc-go/blob/main/docs/migrations/v6-to-v7.md?plain=1#L46-L67
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		authvesting.AppModuleBasic{},
		auramodule.AppModuleBasic{},
		samodule.AppModuleBasic{},
		wasm.AppModuleBasic{},
		ibc_hooks.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		wasmtypes.ModuleName:           {authtypes.Burner},
		// this line is used by starport scaffolding # stargate/app/maccPerms
	}
)

var (
	_ servertypes.Application = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)

	utils.RegisterDenoms()
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            custombankkeeper.BaseKeeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            custommintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCHooksKeeper        *ibchookskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	WasmKeeper            wasmkeeper.Keeper
	ContractKeeper        *wasmkeeper.PermissionedKeeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	AuraKeeper auramodulekeeper.Keeper

	SaKeeper samodulekeeper.Keeper

	// Middleware wrapper
	Ics20WasmHooks   *ibc_hooks.WasmHooks
	HooksICS4Wrapper ibc_hooks.ICS4Middleware

	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// the module manager
	mm *module.Manager

	// the configurator
	configurator module.Configurator
}

// New returns a reference to an initialized Gaia.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	ChainID = GetChainID(appOpts)

	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey, crisistypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, consensusparamtypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		auramoduletypes.StoreKey,
		samoduletypes.StoreKey,
		authzkeeper.StoreKey,
		wasmtypes.StoreKey,
		ibchookstypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], govModAddress)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	app.CapabilityKeeper.Seal()
	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, AccountAddressPrefix,
		govModAddress,
	)

	app.AuraKeeper = auramodulekeeper.NewKeeper(
		appCodec,
		keys[auramoduletypes.StoreKey],
		keys[auramoduletypes.MemStoreKey],
		app.GetSubspace(auramoduletypes.ModuleName),
	)

	app.BankKeeper = custombankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, internal.MergeExcludeAddrs(app.ModuleAccountAddrs()), govModAddress, app.AuraKeeper,
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, govModAddress,
	)

	app.MintKeeper = custommintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], stakingKeeper,
		app.AccountKeeper, app.BankKeeper, app.AuraKeeper, authtypes.FeeCollectorName,
		govModAddress,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		stakingKeeper, authtypes.FeeCollectorName, govModAddress,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, cdc, keys[slashingtypes.StoreKey], stakingKeeper, govModAddress,
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec, keys[crisistypes.StoreKey], invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName, govModAddress,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp, govModAddress)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
		app.AccountKeeper,
	)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)
	app.StakingKeeper = stakingKeeper

	// ... other modules keepers

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibcexported.StoreKey], app.GetSubspace(ibcexported.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govv1beta.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		app.keys[ibchookstypes.StoreKey],
	)
	app.IBCHooksKeeper = &hooksKeeper

	auraPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	wasmHooks := ibc_hooks.NewWasmHooks(app.IBCHooksKeeper, &app.WasmKeeper, auraPrefix) // The contract keeper needs to be set later
	app.Ics20WasmHooks = &wasmHooks
	app.HooksICS4Wrapper = ibc_hooks.NewICS4Middleware(
		app.IBCKeeper.ChannelKeeper,
		app.Ics20WasmHooks,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		// The ICS4Wrapper is replaced by the ICS4Wrapper middleware instead of the channel
		app.HooksICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)
	middlewareTransferModule := ibc_hooks.NewIBCMiddleware(
		transferIBCModule,
		&app.HooksICS4Wrapper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	auraModule := auramodule.NewAppModule(appCodec, app.AuraKeeper)
	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, middlewareTransferModule)

	// ------ CosmWasm setup ------
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	availableCapabilities := strings.Join(wasmapp.AllCapabilities(), ",")
	wasmOpts := GetWasmOpts(appOpts)
	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.BaseApp.MsgServiceRouter(),
		app.BaseApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		govModAddress,
		wasmOpts...,
	)
	app.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(&app.WasmKeeper)

	app.SaKeeper = samodulekeeper.NewKeeper(
		appCodec,
		keys[samoduletypes.StoreKey],
		keys[samoduletypes.MemStoreKey],
		app.GetSubspace(samoduletypes.ModuleName),
		app.WasmKeeper,
		app.ContractKeeper,
		app.AccountKeeper,
	)

	// sa module
	saModule := samodule.NewAppModule(appCodec, app.SaKeeper, app.ContractKeeper, app.AccountKeeper)

	// Pass the contract keeper to ICS4Wrappers for ibc middlewares
	app.Ics20WasmHooks.ContractKeeper = &app.WasmKeeper

	// The gov proposal types can be individually enabled
	enabledProposals := GetEnabledProposals()
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasmtypes.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals)) //nolint:staticcheck // still use same ver 0.41.0 of wasmd
	}

	customGovConfig := govtypes.DefaultConfig()
	customGovConfig.MaxMetadataLen = 2000
	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		stakingKeeper, app.BaseApp.MsgServiceRouter(), customGovConfig, govModAddress,
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	app.GovKeeper.SetLegacyRouter(govRouter)

	// Add wasm module route to the ibc router, then set and seal it
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper))
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		authvesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		custombank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		customfeegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		custommint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, minttypes.DefaultInflationCalculationFn, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		// IBC modules
		ibc_hooks.NewAppModule(app.AccountKeeper),
		auraModule,
		saModule,
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		// this line is used by starport scaffolding # stargate/app/appModule
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		// additional non simd modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		auramoduletypes.ModuleName,
		wasmtypes.ModuleName,
		samoduletypes.ModuleName,
		ibchookstypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		// additional non simd modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		auramoduletypes.ModuleName,
		wasmtypes.ModuleName,
		samoduletypes.ModuleName,
		ibchookstypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		// samodule must occur before genutil so that DeliverGenTx can successfully pass the smart account ante handler
		samoduletypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		auramoduletypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		upgradetypes.ModuleName,
		wasmtypes.ModuleName,
		ibchookstypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
	)

	app.mm.SetOrderMigrations(
		authtypes.ModuleName,
		auramoduletypes.ModuleName,
		banktypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		genutiltypes.ModuleName,
		govtypes.ModuleName,
		ibcexported.ModuleName,
		minttypes.ModuleName,
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		wasmtypes.ModuleName,
		samoduletypes.ModuleName,
		crisistypes.ModuleName,
		ibchookstypes.ModuleName,
	)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// SDK v47 - since we do not use dep inject, this gives us access to newer gRPC services.
	// use follow suggestion of migration docs
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeeGrantKeeper,
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer},
			WasmKeeper:         app.WasmKeeper,
			SmartAccountKeeper: app.SaKeeper,
			IBCKeeper:          app.IBCKeeper,
			WasmConfig:         &wasmConfig,
			TXCounterStoreKey:  keys[wasmtypes.StoreKey],
			Codec:              app.appCodec,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)

	// In v0.46, the SDK introduces _postHandlers_. PostHandlers are like
	// antehandlers, but are run _after_ the `runMsgs` execution. They are also
	// defined as a chain, and have the same signature as antehandlers.
	//
	// In baseapp, postHandlers are run in the same store branch as `runMsgs`,
	// meaning that both `runMsgs` and `postHandler` state will be committed if
	// both are successful, and both will be reverted if any of the two fails.
	//
	// The SDK exposes a default postHandlers chain, which comprises of only
	// one decorator: the Transaction Tips decorator. However, some chains do
	// not need it by default, so feel free to comment the next line if you do
	// not need tips.
	// To read more about tips:
	// https://docs.cosmos.network/main/core/tips.html
	//
	// Please note that changing any of the anteHandler or postHandler chain is
	// likely to be a state-machine breaking change, which needs a coordinated
	// upgrade.
	postHandler, err := NewPostHandler(
		PostHandlerOptions{
			HandlerOptions:     posthandler.HandlerOptions{},
			SmartAccountKeeper: app.SaKeeper,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	app.SetEndBlocker(app.EndBlocker)
	app.setupUpgradeHandlers()
	//app.RegisterUpgradeHandlers(configurator)

	// add wasm snapshot
	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper
	// this line is used by starport scaffolding # stargate/app/beforeInitReturn
	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Gaia's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Gaia's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		docs.RegisterOpenAPIService(Name, apiSvr.Router)
		//apiSvr.Router.Handle("/static/openapi.yml", http.FileServer(http.FS(docs.Docs)))
		//apiSvr.Router.HandleFunc("/", openapiconsole.Handler(Name, "/static/openapi.yml"))
	}

}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

func GetChainID(appOpts servertypes.AppOptions) string {
	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	return chainID
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable()) //nolint:staticcheck
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(auramoduletypes.ModuleName)
	paramsKeeper.Subspace(samoduletypes.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	return paramsKeeper
}

func (app *App) setupUpgradeHandlers() {
	// v0.3.0 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_3_0.UpgradeName,
		v0_3_0.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.3.1 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_3_1.UpgradeName,
		v0_3_1.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.3.2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_3_2.UpgradeName,
		v0_3_2.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.3.3 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_3_3.UpgradeName,
		v0_3_3.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.4.0 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_4_0.UpgradeName,
		v0_4_0.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.4.1 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_4_1.UpgradeName,
		v0_4_1.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.4.2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_4_2.UpgradeName,
		v0_4_2.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.4.4 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_4_4.UpgradeName,
		v0_4_4.CreateUpgradeHandler(app.mm, app.configurator),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v0_4_5.UpgradeName,
		v0_4_5.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// v0.5.0 upgrade handler add new module
	app.UpgradeKeeper.SetUpgradeHandler(
		v500.UpgradeName,
		v500.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v501.UpgradeName,
		v501.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v600.UpgradeName,
		v600.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v601.UpgradeName,
		v601.CreateUpgradeHandler(app.mm, app.SaKeeper, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v700.UpgradeName,
		v700.CreateUpgradeHandler(
			app.mm, app.configurator,
			app.SaKeeper,
			app.ParamsKeeper,
			app.ConsensusParamsKeeper,
			*app.IBCKeeper,
			app.AccountKeeper,
		),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v701.UpgradeName,
		v701.CreateUpgradeHandler(
			app.mm, app.configurator,
			app.SaKeeper,
			app.ParamsKeeper,
			app.ConsensusParamsKeeper,
			*app.IBCKeeper,
			app.AccountKeeper,
		),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v702.UpgradeName,
		v702.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v703.UpgradeName,
		v703.CreateUpgradeHandler(app.mm, app.configurator),
	)

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v0_3_0.UpgradeName:
		// no store upgrades in v0.3.0

	case v0_3_1.UpgradeName:
		// no store upgrades in v0.3.1

	case v0_3_2.UpgradeName:
		// no store upgrades in v0.3.2

	case v0_3_3.UpgradeName:
		// no store upgrades in v0.3.3

	case v0_4_0.UpgradeName:
		// no store upgrades in v0.4.0

	case v0_4_1.UpgradeName:
		// no store upgrades in v0.4.1

	case v0_4_2.UpgradeName:
		// no store upgrades in v0.4.2

	case v0_4_4.UpgradeName:
	// no store upgrades in v0.4.4

	case v500.UpgradeName:
	// no store upgrades in v0.5.0

	case v501.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			//Added: []string{ibcmiddlewaretypes.StoreKey},
		}
	case v600.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{samoduletypes.StoreKey},
		}

	case v601.UpgradeName:
		// no store upgrades in v0.6.

	case v700.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				consensusparamtypes.StoreKey,
				crisistypes.StoreKey,
			},
		}

	case v701.UpgradeName:
		if ChainID == "xstaxy-1" {
			storeUpgrades = &storetypes.StoreUpgrades{
				Added: []string{
					ibchookstypes.StoreKey,
					samoduletypes.StoreKey,
					consensusparamtypes.StoreKey,
					crisistypes.StoreKey,
				},
			}
		}

	case v702.UpgradeName:
	// no store upgrades in v0.7.2
	case v703.UpgradeName:
		// no store upgrades in v0.7.3

	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
