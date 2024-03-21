import { GasPrice, SigningStargateClient } from '@cosmjs/stargate';
import { Secp256k1HdWallet } from '@cosmjs/amino';
import { stringToPath } from '@cosmjs/crypto';

import { createWalletClient, http } from 'viem'
import { localhost } from 'viem/chains'
import { mnemonicToAccount, HDAccount } from 'viem/accounts'


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

let cosmosWallets: Secp256k1HdWallet[];
let cosmosClients: SigningStargateClient[];
let evmAccounts: HDAccount[];
let evmClient;

describe('Bank', () => {
  beforeAll(async () => {
    cosmosWallets = await Promise.all(users.map((user) => {
      return Secp256k1HdWallet.fromMnemonic(user.mnemonic, { hdPaths: [stringToPath("m/44'/118'/0'/0/0")], prefix: 'aura' });
    }))

    cosmosClients = await Promise.all(cosmosWallets.map((wallet) => {
      return SigningStargateClient.connectWithSigner(
        'http://0.0.0.0.:26657',
        wallet,
        {
          gasPrice: GasPrice.fromString('0.025uaura'),
        }
      )
    }));

    evmAccounts = users.map((user) => {
      return mnemonicToAccount(user.mnemonic)
    })

    evmClient = createWalletClient({
      chain: localhost,
      transport: http()
    })
  })

  it('should send tokens from a cosmos address to cosmos address', async () => {
    const [account] = await cosmosWallets[0].getAccounts();
    console.log(account.address)

    const client = await SigningStargateClient.connectWithSigner(
      'http://0.0.0.0:26657',
      cosmosWallets[0],
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

    const result = await client.sendTokens(account.address, recipient, amount, 1.5);
    console.log(result);
  }, 10000);

  it('should send tokens from a cosmos address to evm address', async () => {
    const [account] = await cosmosWallets[0].getAccounts();
    console.log(account.address)

  })
});
