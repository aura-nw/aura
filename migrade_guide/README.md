## Changelog

### v0.2.2

- Upgrade x/wasm to 0.24.0
- Upgrade cosmos sdk to 0.45.2

### v0.2.1

- Add module authz
- Change Wasm param: allow every body can upload their smart contract

### v0.2

- Remove current custom wasm module
- Upgrade to Cosmos SDK `V0.45.1`
- Upgrade to ibc-go v2
- Add wasm module from CosmWasm
- Change voting_period and max_deposit_period from 10 minutes to 12 hours

## References

- [Upgrade to Cosmos SDK `V0.45.1`](https://github.com/cosmos/cosmos-sdk/tree/master/docs/migrations)
- [Upgrade to ibc v2](https://github.com/cosmos/ibc-go/blob/main/docs/migrations/v1-to-v2.md)
- [Add wasm module from CosmWasm](https://docs.cosmwasm.com/docs/1.0/integration#integrating-wasmd)

## Guide upgrade

- Dowload the latest pre-release version

```
mkdir aura45
cd aura45
git clone https://github.com/aura-nw/aura.git
cd aura
git checkout dev

# install dependency and build new aurad, this file aurad will stay in current directory
make build
```

- Deposit and vote for the proposal

This is the first upgrade proposal for testnet.

```
# query proposal
aurad q gov proposals

# deposit to proposal
aurad tx gov deposit <proposal_id> <deposit_amount> --from <key_name> --yes --fees 20uaura --chain-id aura-testnet

# vote for proposal
aurad tx gov vote <proposal_id> yes --from <key_name> --yes --fees 20uaura --chain-id aura-testnet

```

- After proposal is passed, the network will halt once a pre-defined upgrade block height has been reached.
- Stop current aurad daemon.
- Start new aurad daemon.
