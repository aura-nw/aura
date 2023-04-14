# IBC MIDDLEWARE

This module is the stack of the middleware and transfer modules that allows for cross-chain smart contract calls.

As with IBC [v4.3.0](https://github.com/cosmos/ibc-go/releases/tag/v4.3.0), the mechanism for enabling this is a memo field on every ICS20 transfer packet. which show whether ICS20 transfer is normal or whether ICS20 transfer requires contract call action.

# Memo Field
```Javascrip
    "memo": {
        "wasm": {
            "contract": "contract_address",
            "msg": {
                "raw_message_fields": "raw_message_data",
            }
        }
    }
``` 
The `wasm` field in `memo` indicates the contract call requirement for this package
`wasm` format must include these fields:
-   `contract`: valid contract address of receiving chain, must be the same as the *receiver* field in the ICS20 Package
-   `msg`: message call for the above contract

# Receive Flow
We override the **IBC Module**'s OnRecvPacket method to parse and validate the ICS20 Packet's memo field, and then make a contract call if the memo field is valid.

* Parse and Validate memo
    - check if `memo` is a json string
    - check if `memo` has field `wasm`
    - check if `memo['wasm']` has fileds `contract` and `msg`
    - check if `memo['wasm']['contract']` is valid Bench32 address and be the same as the *receiver* field in the ICS20 Package
    - check if `memo['wasm']` has field `msg` and it is json string

* 



