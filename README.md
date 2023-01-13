# Aura
[![go-test](https://github.com/aura-nw/aura/actions/workflows/test.yml/badge.svg)](https://github.com/aura-nw/aura/actions/workflows/test.yml)
[![golangci-lint](https://github.com/aura-nw/aura/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/aura-nw/aura/actions/workflows/golangci-lint.yml)

This repository contains source code for Aurad (Aura Daemon). Aurad binary is the official client for Aura Network. Aurad is built using Cosmos SDK

Aura Network is a NFT-centric blockchain platform that provides infrastructure assisting to bring user assets to the crypto market.

## Prerequisite
- Go 1.18

## Install Aura daemon
Using Makefile
```
    make
```
The **aurad** bin file is located on **${source_directory}/build/** or **GO_PATH** (default ~/go/bin/) 

## Setup a LocalNet

### Initialize the Chain
```
# <moniker> is the custom username of the node
# <chain-id> is the identity of the chain
aurad init <moniker> --chain-id <chain-id>
```
This command will initialize the home folder containing necessary components for your chain  
(default: ~/.aura)

### Customize the genesis file
A genesis file is a JSON file which defines the initial state of your blockchain. It can be seen as height 0 of your blockchain. The first block, at height 1, will reference the genesis file as its parent.

The docs about genesis customization: https://hub.cosmos.network/main/resources/genesis.html

### Create your validator
Create a local key pair for creating validator:
```
aurad keys add <key_name> 
```
Add some tokens to the wallet:
```
aurad add-genesis-account <key_name> <amount><denom>
```
Create a validtor generation transaction:
```
aurad gentx <key_name> <amount><denom> --chain-id <chain-id>
```
Collect the gentx to genesis file:
```
aurad collect-gentxs
```

### Run a node
```
aurad start 
```
## Setup testnet using testnetCmd

## Contribution
The Aurad is still in development by the Aura Network team. For more information on how to contribute to this project, please contact us at support@aura.network

## License
Aurad project source code files are made available under Apache-2.0 License, located in the LICENSE file. Basically, you can do whatever you want as long as you include the original copyright and license notice in any copy of the software/source.
