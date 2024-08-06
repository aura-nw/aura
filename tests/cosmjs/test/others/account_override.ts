import { SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet, StdFee } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseGwei } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { assert } from 'chai';

import { setupClients, auradev } from '../util/test_setup';

let cosmosAccounts: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClients: WalletClient[];
let publicClient = createPublicClient({
  chain: auradev,
  transport: http()
});

// This test is to check if we can override an account when deploying a contract to the same address
// This case happens when a smart account is funded before instantiated as a smart account
describe('Account override', () => {
  before(async () => {
    const testClients = await setupClients('auradev');
    cosmosAccounts = testClients.cosmosAccounts;
    cosmosClients = testClients.cosmosClients;
    evmAccounts = testClients.evmAccounts;
    evmClients = testClients.evmClients;
  })

  it('can override account when deploy contract', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();

    const evmClient = evmClients[0];

    const deployerProxy = '0x4e59b44847b379578588920ca78fbf26c0b4956c' as `0x${string}`
    const sampleData = '6080604052348015600f57600080fd5b5060848061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063c3cafc6f14602d575b600080fd5b6033604f565b604051808260ff1660ff16815260200191505060405180910390f35b6000602a90509056fea165627a7a72305820ab7651cb86b8c1487590004c2444f26ae30077a6b96c6bc62dda37f1328539250029' as `0x${string}`
    // generate random 32bytes as salt
    const salt = '0x' + Array.from({ length: 64 }, () => Math.floor(Math.random() * 16).toString(16)).join('')
    const data = salt + sampleData as `0x${string}`

    // use eth_call to get the deployed contract address
    const contractAddress = await publicClient.request({
      method: 'eth_call',
      params: [{
        to: deployerProxy,
        data
      },
        "latest"
      ],
      id: 1
    })

    const balance = await publicClient.request({
      method: 'eth_getBalance',
      params: [contractAddress, 'latest'],
      id: 1
    })
    assert.equal(balance, '0x0')
    const cosmosAccountRes = await fetch(`https://lcd.dev.aura.network/evmos/evm/v1/cosmos_account/${contractAddress}`)
    const cosmosAddress = (await cosmosAccountRes.json()).cosmos_address

    const fee = {
      amount: [{ amount: '200000', denom: 'utaura' }],
      gas: '200000'
    } as StdFee
    await cosmosClients[0].sendTokens(
      account.address,
      cosmosAddress,
      [{ denom: 'utaura', amount: '1000000' }],
      fee,
      ""
    )

    let accountInfoRes = await fetch(`https://lcd.dev.aura.network/cosmos/auth/v1beta1/accounts/${cosmosAddress}`)
    let accountInfo = await accountInfoRes.json()
    assert(accountInfo.account['@type'] === '/cosmos.auth.v1beta1.BaseAccount', 'account type does not match')

    // deploy the contract by raw transaction
    const txHash = await evmClient.sendTransaction({
      account: evmAccounts[0],
      to: deployerProxy,
      data,
      maxPriorityFeePerGas: parseGwei('35'),
      maxFeePerGas: parseGwei('35'),
      gasLimit: 2000000n,
      value: 0n,
      chain: auradev
    })

    const txReceipt = await publicClient.waitForTransactionReceipt({ hash: txHash });

    // get code from the deployed contract
    const bytecode = await publicClient.request({
      method: 'eth_getCode',
      params: [contractAddress, 'latest'],
      id: 1
    })

    assert(bytecode !== '0x', ' bytecode does not match')

    accountInfoRes = await fetch(`https://lcd.dev.aura.network/cosmos/auth/v1beta1/accounts/${cosmosAddress}`)
    accountInfo = await accountInfoRes.json()
    assert(accountInfo.account['@type'] === '/ethermint.types.v1.EthAccount', 'account type does not match')
  })
});