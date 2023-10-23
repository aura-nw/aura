package types

const (
	// ModuleName defines the module name
	ModuleName = "smartaccount"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_smartaccount"

	AccountIDKey = "smartaccount_id"

	// In the AnteHandler, if the tx only has one sender and this sender is an
	// AbstractAccount, we store its address here. This way, in the PostHandler,
	// we know whether to call the after_tx method.
	SignerAddressKey = "smartaccount_signer"

	// free gas remaining in ante and post handlers
	// each smartaccount tx will have a limited amount of gas that can be used for free
	// the limitation will be determined by the government
	GasRemainingKey = "gas_remaining"
)

var (
	ParamsKey = []byte{0x00}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
