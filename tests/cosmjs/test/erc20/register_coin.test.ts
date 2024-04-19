import { SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet, StdFee } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther, getAddress } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { evmos, cosmos } from '@aura-nw/aurajs';

import { assert } from 'chai';

import { setupClients, auradev } from '../util/test_setup';


let cosmosAccounts: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClients: WalletClient[];
let publicClient = createPublicClient({
  chain: localhost,
  transport: http()
});
let erc20Address: `0x${string}`;
const IBCDenom = "ibc/939F7D594BF0C04D914C711F39DA67073B68D39F9619513A752EA4DBC63CA631"

describe('Should register a new coin', () => {
  before(async () => {
    const testClients = await setupClients('localaura');
    cosmosAccounts = testClients.cosmosAccounts;
    cosmosClients = testClients.cosmosClients;
    evmAccounts = testClients.evmAccounts;
    evmClients = testClients.evmClients;
  })

  it('can register an ERC20 token', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();
    const balances = await cosmosClients[0].getAllBalances(account.address);
    console.log('balances', balances)

    const registerMsg = evmos.erc20.v1.RegisterCoinProposal.fromPartial({
      metadata: [{
        description: "IBC Coin",
        denomUnits: [
          {
            denom: IBCDenom,
            exponent: 0,
            aliases: []
          },
          {
            denom: 'usdc',
            exponent: 6,
            aliases: []
          }
        ],
        symbol: "USDC",
        base: IBCDenom,
        display: "usdc",
        name: "USDC",
        uri: "",
        uriHash: ""
      }],
      description: "Register an Native coin",
      title: "Register USDC"
    })
    const registerMsgRaw = evmos.erc20.v1.RegisterCoinProposal.encode(registerMsg).finish();

    const proposalMsg = cosmos.gov.v1beta1.MsgSubmitProposal.fromPartial({
      content: {
        typeUrl: evmos.erc20.v1.RegisterCoinProposal.typeUrl,
        value: registerMsgRaw
      },
      initialDeposit: [
        {
          amount: '1000000',
          denom: 'uaura'
        }
      ],
      proposer: account.address,
    })

    const fee = {
      amount: [{ amount: '40000', denom: 'uaura' }],
      gas: '2000000'
    } as StdFee

    const submitTx = await cosmosClients[0].signAndBroadcast(account.address, [{
      typeUrl: cosmos.gov.v1beta1.MsgSubmitProposal.typeUrl,
      value: proposalMsg
    }], fee, 'Register USDC')
    console.log("submit proposal:", submitTx)

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

    // get token_pairs infomation
    const tokenPairsRes = await fetch('http://0.0.0.0:1317/evmos/erc20/v1/token_pairs').then(res => res.json());

    // token_pairs should have the new coin
    const tokenPair = tokenPairsRes.token_pairs.find((coin: any) => coin.denom === IBCDenom);
    assert.isDefined(tokenPair);

    erc20Address = tokenPair.erc20_address;
    console.log('erc20Address', erc20Address)
  })

  it('can convert native coin to erc20', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();
    cosmosClients[0].registry.register(
      evmos.erc20.v1.MsgConvertCoin.typeUrl,
      evmos.erc20.v1.MsgConvertCoin
    )

    // const erc20Address = '0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd'

    // convert erc20 to coin
    const convertFee = {
      amount: [{ amount: '500000', denom: 'uaura' }],
      gas: '2000000'
    } as StdFee
    const convertMsg = {
      typeUrl: evmos.erc20.v1.MsgConvertCoin.typeUrl,
      value: {
        coin: {
          denom: IBCDenom,
          amount: '100000'
        },
        receiver: evmAccounts[0].address,
        sender: account.address
      }
    }
    console.log('convertMsg', convertMsg)
    const convertCoinTx = await cosmosClients[0].signAndBroadcast(account.address, [convertMsg], convertFee, 'Convert native coin to erc20')

    console.log(convertCoinTx)

    //get balances from erc20
    const erc20CoinBalance = await publicClient.readContract({
      address: erc20Address,
      abi: [{
        inputs: [{ type: 'address' }],
        name: 'balanceOf',
        outputs: [{ type: 'uint256' }],
        stateMutability: 'view',
        type: 'function'
      }],
      functionName: 'balanceOf',
      args: [evmAccounts[0].address],
    })
    console.log(erc20CoinBalance)

    assert.equal(erc20CoinBalance, 100000n);
  });

  it('cannot convert if not enough balance', async () => {
    const [account] = await cosmosAccounts[1].getAccounts();
    cosmosClients[1].registry.register(
      evmos.erc20.v1.MsgConvertCoin.typeUrl,
      evmos.erc20.v1.MsgConvertCoin
    )

    // convert erc20 to coin
    const convertFee = {
      amount: [{ amount: '500000', denom: 'uaura' }],
      gas: '2000000'
    } as StdFee
    const convertMsg = {
      typeUrl: evmos.erc20.v1.MsgConvertCoin.typeUrl,
      value: {
        coin: {
          denom: IBCDenom,
          amount: '100000'
        },
        receiver: evmAccounts[1].address,
        sender: account.address
      }
    }
    console.log('convertMsg', convertMsg)
    try {
      const tx = await cosmosClients[1].signAndBroadcast(account.address, [convertMsg], convertFee, 'Convert native coin to erc20')
      console.log('E', tx)
    } catch (error) {
      console.log('error', error)
      // assert.equal(error.message, 'rpc error: code = Unknown desc = insufficient funds')
    }
  })
});
