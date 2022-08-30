package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const (
	flagVestingStart  = "vesting-start-time"
	flagVestingPeriod = "period-length"
	flagVestingAmt    = "total-vesting-amount"
	flagVestingTime   = "total-vesting-time"
	flagCliffTime     = "cliff-time"
	flagCliffAmount   = "cliff-amount"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisVestingAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-periodic-vesting-account [address_or_key_name] [coin]",
		Short: "Add a genesis periodic vesting account to genesis.json",
		Long: `Add a genesis periodic vesting account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.Codec
			cdc := depCdc

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			vestingStart, err := cmd.Flags().GetInt64(flagVestingStart)
			if err != nil {
				return err
			}
			periodLength, err := cmd.Flags().GetInt64(flagVestingPeriod)
			if err != nil {
				return err
			}
			vestingTime, err := cmd.Flags().GetInt64(flagVestingTime)
			if err != nil {
				return err
			}
			vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
			if err != nil {
				return err
			}
			cliffTime, err := cmd.Flags().GetInt64(flagCliffTime)
			if err != nil {
				return err
			}
			cliffAmtStr, err := cmd.Flags().GetString(flagCliffAmount)
			if err != nil {
				return err
			}

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}
			cliffAmt, err := sdk.ParseCoinsNormalized(cliffAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse cliff amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingStart+vestingTime)
				fmt.Print(vestingAmt.Sort(), "\n")
				if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
					baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
					return errors.New("vesting amount cannot be greater than total amount")
				}

				if vestingStart != 0 && periodLength != 0 && cliffAmt[0].Amount.LTE(vestingAmt[0].Amount) && cliffTime <= vestingTime {
					vestingTime = vestingTime - cliffTime
					var numPeriod int64 = vestingTime / periodLength

					// Currently, only allow to vest 1 type of coin per account
					// Add 1 period if set cliff
					var totalAmount sdk.Int = vestingAmt[0].Amount.Sub(cliffAmt[0].Amount)
					var periodicAmount sdk.Int = totalAmount.QuoRaw(numPeriod)
					if cliffTime > 0 {
						numPeriod = numPeriod + 1
					}
					periods := caculateVestingPeriods(vestingTime, periodLength, vestingAmtStr, numPeriod, totalAmount, periodicAmount, cliffTime, cliffAmtStr)
					genAccount = authvesting.NewPeriodicVestingAccountRaw(baseVestingAccount, vestingStart, periods)
				} else {
					return errors.New("invalid vesting parameters")
				}
			} else {
				return errors.New("command is only allowed to create periodic vesting account")
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(depCdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Int64(flagVestingStart, 0, "schedule start time (unix epoch in seconds) for vesting accounts")
	cmd.Flags().Int64(flagVestingPeriod, 0, "length of the period (in seconds)")
	cmd.Flags().Int64(flagVestingTime, 0, "total vesting time (in seconds)")
	cmd.Flags().Int64(flagCliffTime, 0, "Cliff time (in seconds)")
	cmd.Flags().String(flagCliffAmount, "0uaura", "Cliff amount")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func caculateVestingPeriods(vestingTime int64, periodLength int64, vestingAmtStr string, numPeriod int64, totalAmount sdk.Int, periodicAmount sdk.Int, cliffTime int64, cliffAmtStr string) authvesting.Periods {
	var counter int = 0
	if vestingTime%periodLength != 0 {
		// indivisible vesting time
		periods := make([]authvesting.Period, numPeriod+1)
		if cliffTime > 0 {
			periods[0].Length = cliffTime
			periods[0].Amount, _ = sdk.ParseCoinsNormalized(cliffAmtStr)
			counter = 1
			numPeriod = numPeriod - 1
		}
		for i := counter; i < int(numPeriod+2); i++ {
			periods[i].Length = periodLength
			periods[i].Amount, _ = sdk.ParseCoinsNormalized(vestingAmtStr)
			periods[i].Amount[0].Amount = periodicAmount
			if !totalAmount.ModRaw(numPeriod).IsZero() && int64(i) == numPeriod {
				periods[i].Length = vestingTime % periodLength
				periods[i].Amount[0].Amount = totalAmount.ModRaw(numPeriod).Add(periodicAmount)
			}
		}
		return periods
	} else {
		// divisible vesting time
		periods := make([]authvesting.Period, numPeriod)
		if cliffTime > 0 {
			periods[0].Length = cliffTime
			periods[0].Amount, _ = sdk.ParseCoinsNormalized(cliffAmtStr)
			counter = 1
			numPeriod = numPeriod - 1
		}
		for i := counter; i < int(numPeriod+1); i++ {
			periods[i].Length = periodLength
			periods[i].Amount, _ = sdk.ParseCoinsNormalized(vestingAmtStr)
			periods[i].Amount[0].Amount = periodicAmount
			if !totalAmount.ModRaw(numPeriod).IsZero() && int64(i) == (numPeriod-1) {
				periods[i].Amount[0].Amount = totalAmount.ModRaw(numPeriod).Add(periodicAmount)
			}
		}
		return periods
	}
}
