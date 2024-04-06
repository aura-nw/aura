import { StdFee, SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet, StdFee } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther, getContract } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { evmos, cosmos, getSigningCosmosClient } from '@aura-nw/aurajs';

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
    const txHash = await evmClients[0].deployContract({
      abi: TestErc20Abi,
      account: evmAccounts[0],
      bytecode: TestErc20Code.bytecode as `0x${string}`,
      args: ['TestToken', 'TTT', parseEther('1000000')],
      chain: undefined
    })

    const txReceipt = await publicClient.waitForTransactionReceipt({ hash: txHash });

    if (txReceipt.contractAddress) {
      erc20Contract = getContract({
        address: txReceipt.contractAddress,
        abi: TestErc20Abi,
        client: publicClient,
      })
    }
  })

  it('can register an ERC20 token', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();

    const registerMsg = evmos.erc20.v1.RegisterERC20Proposal.fromPartial({
      erc20addresses: [
        erc20Contract.address
      ],
      description: "Register an TestErc20 token",
      title: "Register TestErc20"
    })
    const registerMsgRaw = evmos.erc20.v1.RegisterERC20Proposal.encode(registerMsg).finish();

    const proposalMsg = cosmos.gov.v1beta1.MsgSubmitProposal.fromPartial({
      content: {
        typeUrl: evmos.erc20.v1.RegisterERC20Proposal.typeUrl,
        value: registerMsgRaw
      },
      initialDeposit: [
        {
          amount: '1000000',
          denom: 'uaura'
        }
      ],
      proposer: account.address,
      // authority: "aura10d07y265gmmuvt4z0w9aw880jnsr700jp5y852"
    })

    const fee = {
      amount: [{ amount: '200000', denom: 'uaura' }],
      gas: '400000'
    } as StdFee
    // await cosmosClients[0].sendTokens(account.address,  "aura10d07y265gmmuvt4z0w9aw880jnsr700jp5y852", [{ denom: 'uaura', amount: '1000000' }], fee)


    console.log(account);
    console.log(await cosmosClients[0].getAccount(account.address));

    console.log(await cosmosClients[0].getBlock());
    const tx = await cosmosClients[0].signAndBroadcast(account.address, [{
      // typeUrl: cosmos.gov.v1.MsgExecLegacyContent.typeUrl,
      typeUrl: cosmos.gov.v1beta1.MsgSubmitProposal.typeUrl,
      value: proposalMsg
    }], fee, 'Register TestErc20');
    console.log(tx);
    // decode authInfoBytes
    // const authInfo = cosmos.tx.v1beta1.AuthInfo.decode(tx.authInfoBytes);
    // console.log(JSON.stringify(authInfo, null, 2));

    // const body = cosmos.tx.v1beta1.TxBody.decode(tx.bodyBytes);
    // console.log(JSON.stringify(body, null, 2));

    // await cosmosClients[0].sendTokens(account.address, account.address, [{ denom: 'uaura', amount: '1000000' }], fee);
  });
});
