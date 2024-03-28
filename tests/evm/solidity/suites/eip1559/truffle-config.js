const HDWalletProvider = require('@truffle/hdwallet-provider')
var mnemonic = 'copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom'

module.exports = {
  networks: {
    // Development network is just left as truffle's default settings
    evmos: {
      host: '127.0.0.1', // Localhost (default: none)
      port: 8545, // Standard Ethereum port (default: none)
      network_id: '*', // Any network (default: none)
      gas: 5000000, // Gas sent with each transaction
      gasPrice: 1000000000, // 1 gwei (in wei)
      provider: function () {
        return new HDWalletProvider(mnemonic, 'http://127.0.0.1:8545')
      }
    },
  },
  compilers: {
    solc: {
      version: '0.8.18'
    }
  }
}