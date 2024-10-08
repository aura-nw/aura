# Aura

[![Release](https://github.com/aura-nw/aura/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/aura-nw/aura/actions/workflows/release.yml)

This repository contains source code for Aurad (Aura Daemon). Aurad binary is the official client for Aura Network. Aurad is built using Cosmos SDK

Aura Network is a NFT-centric blockchain platform that provides infrastructure assisting to bring user assets to the crypto market.

## Prerequisite

- Go 1.20

## Install Aura daemon

Using Makefile

```bash
    make
```

The **aurad** bin file is located on **${source_directory}/build/** or **GO_PATH** (default ~/go/bin/)

## Setup a LocalNet

### Initialize the Chain

```bash
# <moniker> is the custom username of the node
# <chain-id> is the identity of the chain
aurad init <moniker> --chain-id <chain-id>
```

This command will initialize the home folder containing necessary components for your chain  
(default: ~/.aura)

### Customize the genesis file

A genesis file is a JSON file which defines the initial state of your blockchain. It can be seen as height 0 of your blockchain. The first block, at height 1, will reference the genesis file as its parent.

The docs about genesis customization: <https://hub.cosmos.network/main/resources/genesis.html>

### Create your validator

Create a local key pair for creating validator:

```bash
aurad keys add <key_name> 
```

Add some tokens to the wallet:

```bash
aurad add-genesis-account <key_name> <amount><denom>
```

Create a validtor generation transaction:

```bash
aurad gentx <key_name> <amount><denom> --chain-id <chain-id>
```

Collect the gentx to genesis file:

```bash
aurad collect-gentxs
```

### Run a node

```bash
aurad start 
```

## Run a local test node

```bash
sh scripts/testnode.sh
```

## Setup testnet using testnetCmd

## Contribution

The Aurad is still in development by the Aura Network team. For more information on how to contribute to this project, please contact us at <support@aura.network>

## License

Aura project source code files are made available under Apache-2.0 License, located in the LICENSE file. Basically, you can do whatever you want as long as you include the original copyright and license notice in any copy of the software/source.

## Acknowledgments

Aura project is built using [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and uses additional modules:
- ```github.com/evmos/evmos/v18``` by Tharsis Labs Ltd.(Evmos). This EVM library is distributed under [ENCL-1.0](https://github.com/evmos/evmos/blob/v16.0.3/LICENSE).

- ```x/evmutil``` by Kava Labs, Inc. This module is distributed under [Apache v2 License](https://github.com/Kava-Labs/kava/blob/master/LICENSE.md).
