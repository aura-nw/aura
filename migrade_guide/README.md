## Changelog

- Remove current custom wasm module
- Upgrade to Cosmos SDK `V0.45.1`
- Upgrade to ibc-go v2
- Add wasm module from CosmWasm

## References

- [Upgrade to Cosmos SDK `V0.45.1`](https://github.com/cosmos/cosmos-sdk/tree/master/docs/migrations)
- [Upgrade to ibc v2](https://github.com/cosmos/ibc-go/blob/main/docs/migrations/v1-to-v2.md)
- [Add wasm module from CosmWasm](https://docs.cosmwasm.com/docs/1.0/integration#integrating-wasmd)

## Guide upgrade

- Clone new source code or checkout current source code to branch feature/add_module_wasm

```
mkdir aura45
cd aura45
git clone https://github.com/aura-nw/aura.git
cd aura
git checkout feature/add_module_wasm

# install dependency and build new aurad, this file aurad will stay in current directory
make build
```

- Create proposal to upgrade

```
# block height upgrade happens
HEIGHT=100

# create proposal for this upgrading (name: v2, height: $HEIGHT)
aurad tx gov submit-proposal software-upgrade v2 --title upgrade-to-0.45 --description upgrade0.45 --upgrade-height $HEIGHT --from Hanoi --yes --fees 20uaura --chain-id aura-testnet

# deposit to proposal
aurad tx gov deposit 1 20000000uaura --from Hanoi --yes --fees 20uaura --chain-id aura-testnet

# vote for proposal until it is passed
aurad tx gov vote 1 yes --from Hanoi --yes --fees 20uaura --chain-id aura-testnet

# query proposal
aurad q gov proposals
```

- After proposal is passed, aurad daemon will stuck at block height $HEIGHT
- Stop current aurad daemon
- Start new aurad daemon
