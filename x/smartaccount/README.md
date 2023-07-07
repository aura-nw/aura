# SMART-ACCOUNT

a [smart account][4] solution for [CosmWasm][1]-enabled chains.

In our context, `smart account` is a contract account associated with a public key, so it can be considered a programmable EOA. The difference is that unlike an EOA account where **Address** and **PubKey** must be the same, the **PubKey** of a `smart account` can be any key that the account owner set for it. 

Our new account will have **SmartAccount** type instead of [BaseAccount][5] or other existing types.

</br>

## Activation Account 
Like EOA, users can create a local `smart account` and decide when to actually use it. This achived by using cosmwasm [Instantiate2][3] method which will generate an account with predictable address. 
- `Instantiate2` method params 
    - **sender**: actor that signerd the messages. 
    - **admin**: optional address that can execute migrations. 
    - **code_id**: reference to the stored WASM code. 
    - **label**: optional metadata to be stored with a contract instance. 
    - **msg**: json encoded message to be passed to the contract on instantiation. 
    - **funds**: coins that are transferred to the contract on instantiation. 
    - **salt**: an arbitrary value provided by the sender. Size can be 1 to 64. 
    - **fix_msg**: include the msg value into the hash for the predictable address. Default is false.  

</br>

### Query `generate-account` 
Allows users to create smart account addresses based on optional configuration
```Go
type QueryGenerateAccount struct{
    // reference to the stored WASM code, must be in whitelist
    code_id    uint64

    // an arbitrary value provided by the sender. Size can be 1 to 64.
    salt       []byte

    // json encoded message to be passed to the contract on instantiation
    init_msg   []byte

    // public key of this account, must be cosmos supported schemas
    public_key Any
}
``` 

Internally a address is built by `Instantiate2` containing:
```
(len(checksum) | checksum | len(sender_address) | sender_address | len(salt) | salt| len(initMsg) | initMsg)
```
When create a new EOA, users can generate their private key locally and claim their account without sending any transactions. In our smart account case, using `public_key` as `sender_address`, a smart account can be claimed to be owned by the user who has configured the parameters to generate the account address.

</br>

### Message `activate-account`
Allows users to activate their smart account using a pre-generated one. This message will take account with type **BaseAccount** and convert it to **SmartAccount** type with preconfigured public key.
```Go
type MessageActivateAccount struct {
    // AccountAddress is the actor who signs the message
    account_address string

    // reference to the stored WASM code, must be in whitelist
    code_id         uint64

    // an arbitrary value provided by the sender. Size can be 1 to 64.
    salt            []byte

    // json encoded message to be passed to the contract on instantiation
    init_msg        []byte

    // public key of this account, must be cosmos supported schemas
    public_key      Any
}
```

This message is signed by the user's private key, and the signer's address is a pre-generated one. The module will recalculate the address based on the user input then check if it is equal to the signer's address, so other parameters must be the same as the configuration used to generate the account address and the signer must have enough funds to pay for the transaction.

</br>

To illustrate this in a graph:

```plain
             tx
              ↓
  ┌───── Antehandler ─────┐
  │   ┌───────────────┐   │
  │   │  decorator 0  │   │
  │   └───────────────┘   │
  │   ┌───────────────┐   │ Set temporary PubKey      ┌───────────────┐
  │   │  SA SetPubKey │---│-------------------------->|  auth module  |
  │   └───────────────┘   │                           └───────────────┘
  |      ...........      |
  │   ┌───────────────┐   │
  │   │SigVerification│   │
  │   └───────────────┘   │
  |      ...........      |
  │   ┌───────────────┐   │ Remove temporary Pubkey   ┌───────────────┐
  │   │  SA decorator │---│-------------------------->|  auth module  |
  │   └───────────────┘   │                           └───────────────┘
  │   ┌───────────────┐   │
  │   │  decorator 2  │   │
  │   └───────────────┘   │
  └───────────────────────┘
              ↓
      ┌────────────────┐
      │  msg activate  │ 
      └────────────────┘
              ↓
  ┌───────  module ───────┐
  │   ┌───────────────┐   │   instantiate contract   ┌───────────────┐
  │   │   SA module   │---│------------------------->|  wasmd module |
  │   └───────────────┘   │                          └───────────────┘
  └───────────────────────┘
              ↓
            done
```
- **AnteHandler**
    - Since the account doesn't have **PubKey** yet, for signature verification, `SA SetPubKey` will set a temporary **PubKey** for this account using the `public_key` parameter in the message.
    - After successful signature verification, `SA decorator` will remove temporary **PubKey** so that `SA module` can initiate contract with this account later (action remove only needed in DeliveryTx).
- **SA module**
    - if the message meets all the checks, the module initiates a contract based on its parameters. The new contract will be linked to the pre-generated account (contract address will be same as account address). The module will then convert account to type `SmartAccount` and set **PubKey** for it. Finnaly, save account to `auth` module.

</br>

**Required**
- valid parameters.
- The account must have received the funds, so it can exist on-chain as **BaseAccount** type with an account number, sequence and empty public key before activated.
- The account address was not used to initiate any smart contract before.
- In some cases, we also allow reactivation of activated accounts that are not linked to any smart contracts.
- **code_id** must be in whitelist.
- `actiavte message` is the only message in tx.
 
</br>

## Recover Account
We provide a smart account recovery way in case the user forgets the account's private key or wants to change the new one. Recovery is simply changing the **PubKey** of an account of type **SmartAccount** with the new setting. This is not a required function so users can choose whether their account is recoverable or not.

</br>

### Message `recover`
The caller specifies the address to recover, the public key and provides credentials to prove they have the authority to perform the recovery.
```Go
type MsgRecover struct{
  // Sender is the actor who signs the message
  creator     string

  // smart-account address that want to recover
  address     string

  // New PubKey using for signature verification of this account, 
  // must be cosmos supported schemas
  public_key  Any

  // Credentials
  credentials string
}
```
The module makes a call to the `recover` method of contract that linked to smart account. If the message parameters meet the verification logic implemented in the contract, the smart account will be updated with the new **PubKey**.
- `recover` call
    ```Go
    type RecoverTx struct {
        Caller      string 
        PubKey      []byte
        Credentials []byte
    }
    ```

</br>

To illustrate this in a graph:

```plain
             tx
              ↓
  ┌───── Antehandler ─────┐
  │   ┌───────────────┐   │
  │   │  decorator 0  │   │
  │   └───────────────┘   │
  │   ┌───────────────┐   │
  │   │  decorator 1  │   │
  │   └───────────────┘   │ 
  │   ┌───────────────┐   │
  │   │  decorator 2  │   │
  │   └───────────────┘   │
  └───────────────────────┘
              ↓
      ┌────────────────┐
      │  msg recover   │ 
      └────────────────┘
              ↓
  ┌───────  module ───────┐
  │   ┌───────────────┐   │   `recover` sudo      ┌───────────────┐
  │   │   SA module   │---│---------------------->|  wasmd module |
  │   └───────────────┘   │                       └───────────────┘
  └───────────────────────┘
              ↓
            done
```
- **SA Module**
    - The `SA module` checks if the requested account is of type `SmartAccount`, if not, rejects it.
    - if `recover` call success, module will update new **PubKey** for account then save account to `auth` module.

</br>

**Required**
- valid parameters.
- Account with *address* must exists on-chain and has type **SmartAccount**.
- Account enables recovery function by implementing `recover` method in **sudo** entry point of linked smart contract.

</br> 

## Smart account Tx 
When build transaction with smart account, user must includes `Validate message` into tx. This message has type **MsgExecuteContract** and will use to trigger smart contract logic that applies to this account.

`Validate message` call to `after_execute` method of contract that linked with smart account. It's value is information of all other messages included in tx. Firstly, the module uses this message data to execute a contract's before the tx is passed to the mempool. The execute calls the `pre_execute` method for pre-validation tx, if it fails, the tx will be denied to enter the mempool and the user will not incur gas charges for it. Finnaly, after all messages included in tx are executed, `Validate message` will be executed to determine whether tx was successful or not.

</br>

To illustrate this in a graph:

```plain
  ┌──────────tx──────────┐
  │  ┌────────────────┐  │
  │  │      msg 0     │  │
  │  └────────────────┘  │
  │  ┌────────────────┐  │
  │  │      msg 1     │  │
  │  └────────────────┘  │
  |     .............    |
  │  ┌────────────────┐  │
  │  │  msg validate  │  │
  │  └────────────────┘  │
  └──────────────────────┘
```


```plain
             tx
              ↓
  ┌───── Antehandler ─────┐
  │   ┌───────────────┐   │
  │   │  decorator 0  │   │
  │   └───────────────┘   │
  │   ┌───────────────┐   │
  │   │ SA SetPubKey  │   │
  │   └───────────────┘   │
  │   ┌───────────────┐   │   `pre_execute` execute   ┌───────────────┐
  │   │ SA decorator  │---│-------------------------->|  wasmd module |
  │   └───────────────┘   │                           └───────────────┘
  │   ┌───────────────┐   │
  │   │  decorator 2  │   │
  │   └───────────────┘   │
  └───────────────────────┘
              ↓
      ┌────────────────┐
      │      msg 0     │
      └────────────────┘
      ┌────────────────┐
      │      msg 1     │
      └────────────────┘
          ...........
      ┌────────────────┐
      │  msg validate  │ -> MsgExecuteContract `after_execute`
      └────────────────┘
              ↓
  ┌───────  module ───────┐
  │   ┌───────────────┐   │
  │   │    module 1   │   │
  │   └───────────────┘   │
  │   ┌───────────────┐   │
  │   │    module 2   │   │
  │   └───────────────┘   │
  └───────────────────────┘
              ↓
            done
```
- **Antehandler**
    - tx will be identified as signed by smart account. If true, it will be redirected to `SA SetPubKey` and `SA decorator`
    - `smart account` tx will go through the `SA SetPubKey` decorator instead of the `auth SetPubKey` decorator. This will avoid the check for similarity of **Account Address** and **PubKey**. 
    
</br>

**Required**
- Signer is account with type **SmartAccount**
- `Validate message` must be the last message in tx.
- `Validate message` must has type **MsgExecutedContract** and call to `after_execute` method of smart contract that linked with account
- `Validate Message` data must be compatible with all other tx messages

</br> 

## Params
Parameters are updatable by the module's authority, typically set to the gov module account.
- `max_gas_execute`: limit how much gas can be consumed by the `pre_execute` method
- `whitelist_code_id`: determine which **code_id** can be instantiated as a `smart account`

</br> 

## WASM
To be considered as `smart account`, smart contract linked with account must implement execute methods, `after_execute` and `pre_execute`:
```Rust
// execute method
struct AfterExecute {
    pub msgs: Vec<MsgData>
}

// execute method
struct PreExecute { 
    pub msgs: Vec<MsgData>
}
```
- **MsgData**: is json encoded message
    ```Rust
    struct MsgData {
        pub type_url: String, // url type of message
        pub value:    String, // value of message
        // etc.
        //  MsgData {
        //      type_url: "/cosmos.bank.v1beta1.MsgSend",
        //      value: "{fromAddress:\"aura172r4c7mng5y6ccfqp5klwyulshx6dh2mmd2r0xnmsgugaa754kws8u96pq\",toAddress:\"aura1y3u4ht0p69gz757myr3l0fttchhw3fj2gpeznd\",amount:[{denom:\"uaura\",amount:\"200\"}]}"
        //  }
    }
    ```
- `pre_execute` method must not consumes exceed `max_gas_execute` gas 

Optional sudo method recover that activate the smart account recovery function
```Rust
// sudo method
struct Recover {
    pub caller: String,
    pub pub_key: Binary,
    pub credentials: Binary,
}
```

[smart account samples][2]

[1]: https://cosmwasm.com/
[2]: https://github.com/aura-nw/smart-account-sample/
[3]: https://github.com/CosmWasm/wasmd/blob/main/x/wasm/keeper/msg_server.go#L79-L110
[4]: https://aura-network.notion.site/Smart-Account-e69e51d6449b46dcb7c157a325dfb44f
[5]: https://github.com/cosmos/cosmos-sdk/blob/main/x/auth/types/account.go