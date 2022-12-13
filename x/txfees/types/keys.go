package types

const (
	// ModuleName defines the module name
	ModuleName = "txfees"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_txfees"

	// FeeCollectorName the module account name for the fee collector account address.
	FeeCollectorName = "fee_collector"

	// TxFeeCollectorName the module account name for the alt fee collector account address.
	TxFeeCollectorName = "tx_fee_collector"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
