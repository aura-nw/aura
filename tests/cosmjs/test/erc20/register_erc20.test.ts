import { SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet, StdFee } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther, getAddress, PublicClient } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { evmos, cosmos } from '@aura-nw/aurajs';

import hre from "hardhat";
import { assert } from 'chai';

import { convertBech32AddressToEthAddress } from '../util/convert_address';
import { USERS, setupClients, localaura, auradev } from '../util/test_setup';


let cosmosAccounts: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClients: WalletClient[];
let publicClient: PublicClient;
let erc20ContractAddr: `0x${string}`;

describe('Should work with ERC20 tokens', () => {
  before(async () => {
    const testClients = await setupClients('auradev');
    cosmosAccounts = testClients.cosmosAccounts;
    cosmosClients = testClients.cosmosClients;
    evmAccounts = testClients.evmAccounts;
    evmClients = testClients.evmClients;
    publicClient = testClients.publicClient;

    const TestErc20Code = await hre.ethers.getContractFactory("TestERC20");
    // console.log(await evmAccounts[0].signMessage({ message: 'hello world' }))
    const TestErc20Abi = JSON.parse(TestErc20Code.interface.formatJson())
    const txHash = await evmClients[0].deployContract({
      abi: TestErc20Abi,
      account: evmAccounts[0],
      bytecode: TestErc20Code.bytecode as `0x${string}`,
      args: ['TestToken', 'TTT', parseEther('1000000')],
      chain: auradev
    })

    const txReceipt = await publicClient.waitForTransactionReceipt({ hash: txHash });
    console.log(txReceipt)

    if (!txReceipt.contractAddress) {
      throw new Error('Contract address not found');
    }
    erc20ContractAddr = txReceipt.contractAddress;
    console.log(erc20ContractAddr)

    // send some token to the a cosmos account
    const [cosmosAccount] = await cosmosAccounts[0].getAccounts();
    const receiver = convertBech32AddressToEthAddress('aura', cosmosAccount.address)
    const sendAmt = 1000000n

    const transferTx = await evmClients[0].writeContract({
      address: erc20ContractAddr,
      abi: TestErc20Abi,
      functionName: 'transfer',
      args: [receiver, sendAmt],
      account: evmAccounts[0],
      chain: auradev
    })

    console.log(transferTx)
  })

  it('can register an ERC20 token', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();

    const registerMsg = evmos.erc20.v1.RegisterERC20Proposal.fromPartial({
      erc20addresses: [
        erc20ContractAddr
      ],
      description: "Register an TestErc20 token",
      title: "Register TestErc20"
    })
    console.log(registerMsg)
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

    console.log(account.address)
    const submitTx = await cosmosClients[0].signAndBroadcast(account.address, [{
      // typeUrl: cosmos.gov.v1.MsgExecLegacyContent.typeUrl,
      typeUrl: cosmos.gov.v1beta1.MsgSubmitProposal.typeUrl,
      value: proposalMsg
    }], fee, 'Register TestErc20')

    console.log(submitTx)

    const proposalId = submitTx?.events.find(
      (event: any) => event.type === 'submit_proposal'
    )?.attributes.find((attr: any) => attr.key === 'proposal_id')?.value;

    if (!proposalId) {
      throw new Error('Proposal ID not found');
    }
    assert.isDefined(proposalId);

    const [validatorAddress] = await cosmosAccounts[4].getAccounts();

    const voteTx = await cosmosClients[4].signAndBroadcast(validatorAddress.address, [{
      typeUrl: cosmos.gov.v1beta1.MsgVote.typeUrl,
      value: {
        option: cosmos.gov.v1beta1.VoteOption.VOTE_OPTION_YES,
        proposalId: proposalId,
        voter: validatorAddress.address
      }
    }], fee, 'Vote for Register TestErc20')

    // wait 10 seconds
    await new Promise((resolve) => setTimeout(resolve, 10000));

    // get the proposal
    const proposalReq = cosmos.gov.v1beta1.QueryProposalRequest.fromJSON({
      proposalId: proposalId
    })

    const queryClient = await cosmos.ClientFactory.createLCDClient({
      restEndpoint: 'http://0.0.0.0:1317'
    })

    const { proposal } = await queryClient.cosmos.gov.v1beta1.proposal(proposalReq);

    // assert passed
    assert.equal(proposal.status.toString(), 'PROPOSAL_STATUS_PASSED');

    cosmosClients[0].registry.register(
      evmos.erc20.v1.MsgConvertERC20.typeUrl,
      evmos.erc20.v1.MsgConvertERC20
    )

    // convert erc20 to coin
    const convertFee = {
      amount: [{ amount: '500000', denom: 'utaura' }],
      gas: '1000000'
    } as StdFee
    const sender = convertBech32AddressToEthAddress('aura', account.address)
    const convertErc20Tx = await cosmosClients[0].signAndBroadcast(account.address, [{
      typeUrl: evmos.erc20.v1.MsgConvertERC20.typeUrl,
      value: {
        contractAddress: erc20ContractAddr,
        receiver: account.address,
        amount: '100000',
        sender
      }
    }], convertFee, 'Convert TestErc20 to coin')

    // console.log(convertErc20Tx)

    //get balances from cosmos
    const erc20CoinBalance = await cosmosClients[0].getBalance(account.address, `erc20/${getAddress(erc20ContractAddr)}`);
    // console.log(erc20CoinBalance)

    assert.equal(erc20CoinBalance.amount, '100000');
  });
});
