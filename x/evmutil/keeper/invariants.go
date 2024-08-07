// Derived from https://github.com/Kava-Labs/kava/blob/d500cd12362edd0c64e8065fbc595cd9399b08c2/x/evmutil/keeper/grpc_query.go
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Modifications:
// - Removed BackedCoinInvariant

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/aura-nw/aura/x/evmutil/types"
)

// RegisterInvariants registers the swap module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, bankK types.BankKeeper, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "fully-backed", FullyBackedInvariant(bankK, k))
	ir.RegisterRoute(types.ModuleName, "small-balances", SmallBalancesInvariant(bankK, k))
}

// AllInvariants runs all invariants of the swap module
func AllInvariants(bankK types.BankKeeper, k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		if res, stop := FullyBackedInvariant(bankK, k)(ctx); stop {
			return res, stop
		}
		return SmallBalancesInvariant(bankK, k)(ctx)
	}
}

// FullyBackedInvariant ensures all minor balances are backed by the coins in the module account.
//
// The module balance can be greater than the sum of all minor balances. This can happen in rare cases
// where the evm module burns tokens.
func FullyBackedInvariant(bankK types.BankKeeper, k Keeper) sdk.Invariant {
	broken := false
	message := sdk.FormatInvariant(types.ModuleName, "fully backed broken", "sum of minor balances greater than module account")

	return func(ctx sdk.Context) (string, bool) {
		totalMinorBalances := sdk.ZeroInt()
		k.IterateAllAccounts(ctx, func(acc types.Account) bool {
			totalMinorBalances = totalMinorBalances.Add(acc.Balance)
			return false
		})

		bankAddr := authtypes.NewModuleAddress(types.ModuleName)
		bankBalance := bankK.GetBalance(ctx, bankAddr, CosmosDenom).Amount.Mul(ConversionMultiplier)

		broken = totalMinorBalances.GT(bankBalance)

		return message, broken
	}
}

// SmallBalancesInvariant ensures all minor balances are less than the overflow amount, beyond this they should be converted to the major denom.
func SmallBalancesInvariant(_ types.BankKeeper, k Keeper) sdk.Invariant {
	broken := false
	message := sdk.FormatInvariant(types.ModuleName, "small balances broken", "minor balances not all less than overflow")

	return func(ctx sdk.Context) (string, bool) {
		k.IterateAllAccounts(ctx, func(account types.Account) bool {
			if account.Balance.GTE(ConversionMultiplier) {
				broken = true
				return true
			}
			return false
		})
		return message, broken
	}
}
