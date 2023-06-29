# SMART-ACCOUNT

a [smart account][4] solution for [CosmWasm][1]-enabled chains.

In our context, `smart account` is a contract account associated with a public key, so it can be considered a programmable EOA. The difference is that unlike an EOA account where **Address** and **PubKey** must be the same, the **PubKey** of a `smart account` can be any key that the account owner set for it. 

Our new account will have `SmartAccount` type instead of `BaseAccount` or other existing types.

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
struct QueryGenerateAccount {
    // reference to the stored WASM code, must be in whitelist
    uint64 code_id;

    // the infor.sender field of the contract instantiate method
    string sender;

    // json encoded message to be passed to the contract on instantiation
    []byte init_msg;

    // public key of this account, must be cosmos supported schemas
    Any public_key;
}
``` 

When create a new EOA, users can generate their private key locally and claim their account without sending any transactions. In our smart account case, the mechanism enabling this is a `salt` field on `Instantiate2` method.
- **Formula**
    ```Go
    salt = sha512.hash(code_id | sender | init_msg | public_key)
    ```
Using a salt calculation formula, a smart account can be claimed to be owned by the user who has configured the parameters to generate the account address.

</br>

### Message `activate-account`
Allows users to activate their smart account using a pre-generated. This message will take `BaseAccount` type and convert it to `SmartAccount` type with preconfigured public key.
```Go
struct MessageActivateAccount {
    // AccountAddress is the actor who signs the message
    string account_address;

    // reference to the stored WASM code, must be in whitelist
    uint64 code_id;

    // the infor.sender field of the contract instantiate method
    string sender;

    // json encoded message to be passed to the contract on instantiation
    []byte init_msg;

    // public key of this account, must be cosmos supported schemas
    Any public_key;
}
```

This message is signed by the user's private key, and the signer's address is a pre-generated one. The module will recalculate the address based on the user input then check if it is equal to the signer, so other parameters must be the same as the configuration used to generate the account address and the signer must have enough funds to pay for the transaction.

**Required**
- valid parameters.
- The signer must have received the funds, so it can exist on-chain as `BaseAccount` type with an account number, sequence and empty public key.
- The signer address was not used to initiate any smart contract before.
- In some cases, we also allow reactivation of activated accounts that are not associated to any smart contracts.
- Code_id must be in whitelist.
 
</br>

## Recover Account
We provide a smart account recovery way in case the user forgets the account's private key or wants to change the new one. Recovery is simply changing the **PubKey** of an account of type `Smart Account` with the new setting. This is not a required function so users can choose whether their account is recoverable or not.

</br>

### Message `recover`
The caller specifies the address to recover, the public key and provides credentials to prove they have the authority to perform the recovery.
```Go
struct MsgRecover {
  // Sender is the actor who signs the message
  string creator;

  // smart-account address that want to recover
  string address;

  // New PubKey using for signature verification of this account, 
  // must be cosmos supported schemas
  Any public_key;

  // Credentials
  string credentials;
}
```
The module makes a call to the `recover` contract method specified by *address*. If the message parameters meet the verification logic implemented in the contract, the smart account will be updated with the new **PubKey**.
- `recover` call
    ```Rust
    struct Recover {
        // actor who signs the recover message
        pub caller: String,

        // new PubKey
        pub pub_key: Binary,

        // credentials 
        pub credentials: Binary,
    }
    ```

**Required**
- valid parameters.
- Account with *address* must exists on-chain and has type `SmartAccount`.
- Smart account enables recovery function by implementing `recover` method in **sudo** entry point.

</br> 

## Smart account Tx 
When build transaction with smart account, user must includes `Validate message` into tx. This message has type `MsgExecutedContract` and will use to trigger smart contract logic that applies to this account.

`Validate message` call to `after_execute` method of contract that associated with smart account. It's value is information of all other messages included in tx. Firstly, the module uses this message data to execute a contract's query before the tx is passed to the mempool. The query calls the `validate` method for basic validation tx, if it fails, the tx will be denied to enter the mempool and the user will not incur gas charges for it. Finnaly, after all messages included in tx are executed, `Validate message` will be executed to determine whether tx was successful or not.


[1]: https://cosmwasm.com/
[2]: https://github.com/aura-nw/smart-account-sample/
[3]: https://github.com/CosmWasm/wasmd/blob/main/x/wasm/keeper/msg_server.go#L79-L110
[4]: https://aura-network.notion.site/Smart-Account-e69e51d6449b46dcb7c157a325dfb44f