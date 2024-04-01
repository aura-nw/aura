import { GasPrice, SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther, getContract } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { evmos } from '@aura-nw/aurajs';

import hre from "hardhat";
import { assert } from 'chai';

import { convertEthAddressToBech32Address } from '../util/convert_address';
import { USERS, setupClients } from '../util/test_setup';
import { Test } from 'mocha';
import { deployContract } from 'viem/_types/actions/wallet/deployContract';


let cosmosAccounts: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClients: WalletClient[];
let publicClient = createPublicClient({
  chain: localhost,
  transport: http()
});
let erc20Contract: any;

describe('Should work with ERC20 tokens', () => {
  before(async () => {
    const testClients = await setupClients();
    cosmosAccounts = testClients.cosmosAccounts;
    cosmosClients = testClients.cosmosClients;
    evmAccounts = testClients.evmAccounts;
    evmClients = testClients.evmClients;

    const TestErc20Code = await hre.ethers.getContractFactory("TestERC20");
    // console.log(await evmAccounts[0].signMessage({ message: 'hello world' }))
    const TestErc20Abi = JSON.parse(TestErc20Code.interface.formatJson()),
    const TestErc20Address = await evmClients[0].deployContract({
      abi: TestErc20Abi,
      account: evmAccounts[0],
      bytecode: TestErc20Code.bytecode as `0x${string}`,
      args: ['TestToken', 'TTT', parseEther('1000000')],
      chain: undefined
    })

    erc20Contract = getContract({
      address: TestErc20Address,
      abi: TestErc20Abi,
      client: publicClient,
    })
  })

  it('can register an ERC20 token', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();

    console.log(aurajs)
    const registerMsg = aurajs.evmos.erc20.v1.RegisterERC20Proposal.fromPartial({
      erc20addresses: [
      ],
      description: "Register an TestErc20 token",
      title: "Register TestErc20"
    })

  });

});
