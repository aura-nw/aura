import { GasPrice, SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { convertEthAddressToBech32Address } from '../util/convert_address';
import { USERS, setupClients } from '../util/test_setup';


let cosmosAccounts: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClients: WalletClient[];
let publicClient = createPublicClient({
  chain: localhost,
  transport: http()
});

describe('Bank', () => {
  beforeAll(async () => {
    const testClients = await setupClients();
    cosmosAccounts = testClients.cosmosAccounts;
    cosmosClients = testClients.cosmosClients;
    evmAccounts = testClients.evmAccounts;
    evmClients = testClients.evmClients;
  })

  it('should send tokens from a cosmos address to cosmos address', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();

    const recipient = USERS[1].address;
    const amount = [
      {
        denom: 'uaura',
        amount: '1000',
      },
    ];

    await cosmosClients[0].sendTokens(account.address, recipient, amount, 1.5);
  }, 10000);

  it('should send tokens from a cosmos address to evm address', async () => {
    const [account] = await cosmosAccounts[0].getAccounts();
    console.log(account.address)

    const evmAccount = evmAccounts[1].address;

    const recipient = convertEthAddressToBech32Address('aura', evmAccount);
    console.log("EVM Recipient: ", evmAccount);
    console.log("Recipient: ", recipient);

    const prevBalance = await publicClient.getBalance({
      address: evmAccount,
    });

    // 1 Aura, should see 1 eAura in the EVM account
    const amount = [
      {
        denom: 'uaura',
        amount: '1000000',
      },
    ];

    const res = await cosmosClients[0].sendTokens(account.address, recipient, amount, 1.5);

    const balance = await publicClient.getBalance({
      address: evmAccount,
      blockNumber: BigInt(res.height)
    });
    console.log(balance);
    expect(balance).toEqual(prevBalance + parseEther('1'));
  })
});
