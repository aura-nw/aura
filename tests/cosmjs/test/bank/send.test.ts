import { GasPrice, SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet } from '@cosmjs/amino';

import { http, WalletClient, createPublicClient, parseEther } from 'viem'
import { localhost } from 'viem/chains'
import { HDAccount } from 'viem/accounts'

import { convertEthAddressToBech32Address } from '../util/convert_address';
import { setupClients } from '../util/test_setup';


const users = [
  {
    key: "user1",
    mnemonic: "copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom",
    address: "aura1q986wh082dp6wndt58j60hrgsr8kh9wg88awl8"
  },
  {
    key: "user2",
    mnemonic: "maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual",
    address: "aura19uz0dc9j950knkxzdxs2m92q5474sgx335pz6r"
  },
  {
    key: "user3",
    mnemonic: "will wear settle write dance topic tape sea glory hotel oppose rebel client problem era video gossip glide during yard balance cancel file rose",
    address: "aura1xejsnure97tteuqz4wggvl8cla3s9knj65alnu"
  },
  {
    key: "user4",
    mnemonic: "doll midnight silk carpet brush boring pluck office gown inquiry duck chief aim exit gain never tennis crime fragile ship cloud surface exotic patch",
    address: "aura1rfn972g75dhyp586jmyda7vpsuqa4w0syh099q"
  }
];

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

    const client = await SigningStargateClient.connectWithSigner(
      'http://0.0.0.0:26657',
      cosmosAccounts[0],
      {
        gasPrice: GasPrice.fromString('0.025uaura'),
      }
    );

    const recipient = users[1].address;
    const amount = [
      {
        denom: 'uaura',
        amount: '1000',
      },
    ];

    await client.sendTokens(account.address, recipient, amount, 1.5);
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
