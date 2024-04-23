KEY="mykey"
CHAINID="aura_9000-1"
MONIKER="localtestnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.aura*

aurad config keyring-backend $KEYRING
aurad config chain-id $CHAINID

# if $KEY exists it should be deleted
aurad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
aurad init $MONIKER --chain-id $CHAINID 

# Allocate genesis accounts (cosmos formatted addresses)
aurad add-genesis-account $KEY 100000000000000000000000000stake --keyring-backend $KEYRING

# Sign genesis transaction
aurad gentx $KEY 1000000000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
aurad collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
aurad validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
aurad start --pruning=nothing  --minimum-gas-prices=0.0001stake
