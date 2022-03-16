package app

import (
	wasmkeeper "github.com/aura-nw/aura/x/wasm/keeper"
)

const (
	// DefaultAuraInstanceCost is initially set the same as in wasmd
	DefaultAuraInstanceCost uint64 = 60_000
	// DefaultAuraCompileCost set to a large number for testing
	DefaultAuraCompileCost uint64 = 100
)

// AuraGasRegisterConfig is defaults plus a custom compile amount
func AuraGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultAuraInstanceCost
	gasConfig.CompileCost = DefaultAuraCompileCost

	return gasConfig
}

func NewAuraWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(AuraGasRegisterConfig())
}
