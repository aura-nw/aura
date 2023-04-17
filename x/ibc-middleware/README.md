# IBC MIDDLEWARE 

This module is the stack of the middleware and transfer modules that allows for cross-chain smart contract calls. 

As with IBC [v4.3.0](https://github.com/cosmos/ibc-go/releases/tag/v4.3.0), the  mechanism for enabling this is a memo field on every ICS20 transfer packet. which show whether ICS20 transfer is normal or whether ICS20 transfer requires contract  call action. 

## Memo Field 
```Javascript 
    "memo": { 
        "wasm": { 
            "contract": "contract_address", 
            "msg": { 
                "raw_message_fields": "raw_message_data", 
            } 
        } 
    } 
```  
The `wasm` field in `memo` indicates the contract call requirement for this package <br /> 
`wasm` format must include these fields: 
-   `contract`: valid contract address of receiving chain, must be the same as the *receiver* field in the ICS20 Package 
-   `msg`: message call for the above contract 

## Receive Ics20 Packet Flow 
We override the **IBC Module**'s OnRecvPacket method to parse and validate the ICS20 Packet's memo field, and then make a contract call if the memo field is valid. 

* Parse and Validate memo 
    - check if income packet is Ics20 Packet
    - check if `memo` is a json string 
    - check if `memo` has field `wasm` 
    - check if `memo['wasm']` has fields `contract` and `msg` 
    - check if `memo['wasm']['contract']` is valid Bench32 address and be the same  as the *receiver* field in the ICS20 Package 
    - check if `memo['wasm']` has field `msg` and it is json string 

* Create `intermediary account` and Send fund to this address 
    - `intermediary account` = `Bech32(Hash("ibc-wasm-hook-intermediary" || channelID || sender))` 
    - replace ICS20 Packet's receiver by `intermediary account` then execute **IBC Module**'s OnRecvPacket base method with this packet

* Execute Contract Call
    - change ICS20 pakcet's denom to IBC local denom
    - execute Contract Call with:
        - `sender`: `intermediary account`
        - `contract`: `memo['wasm']['contract']`
        - `msg`: `memo['wasm']['msg']`
        - `funds`: ICS20 Packet's funds after change of denom
    - return Ack

Reasons of creating `intermediary account` are
- we cannot trust the sender of an IBC packet
- IBC Middleware only validates the memo string, not what it is and what it can do
i.e the counterparty chain can lie about sender and make the ICS20 Packet to require a contract call in our chain that needs owner privileges

## Send Ics20 Packet Flow
To read and determine whether the ICS20 Packet contains a particular memo field that necessitates the module to store the packet information in order to execute the future contract's callback. we override the SendPacket method of the **ICS4 Wrapper**.

* Process Send ICS20 Packet
    - check if income packet is Ics20 Packet
    - check if `memo` is a json string and has field `ibc_callback`
    - the callback metadata be removed from the pakcet because it has already been processed. Remove the memo from the data so that the packet is transmitted without it if the callback is the only key in the memo that is usable.
    - execute **ICS4 Wrapper**'s SendPacket base method with this packet

* Store packet information
    - check if `memo['ibc_callback']` is valid Bench32 address, because will use it as a callback contract address
    - store information needed for callback with
        - `key`: `("%s::%d", Ics20Packet.GetSourceChannel(), Ics20Packet.GetSequence())`
        - `value`: `memo['ibc_callback']`

## Process Acknowledgement Flow
To handle the Acknowledgment returned when using ICS20 packet sending request contract callback. We override the **IBC Module**'s OnAcknowledgementPacket method.
- execute the **IBC Module**'s OnAcknowledgementPacket base method on this packet
- read module's store with `key` = `("%s::%d", Ics20Packet.GetSourceChannel(), Ics20Packet.GetSequence())` to get callback contract address
- check if callback contract address is valid Bench32 address
- make a call to `sudo` entry point of callback contract with message:
    `{"ibc_lifecycle_complete": {"ibc_ack": {"channel": "%s", "sequence": %d, "ack": %s, "success": %s}}}, SourceChannel, Sequence, ackAsJson, success))`
    - `SourceChannel`: ICS20 Packet's SourceChannel
    - `Sequence`: ICS20 Packet's Sequence
    - `ackAsJson`: json form of return acknowledgement 
    - `success`: **true** or **false**
- delete callback data from module'store using `key` above if call success

`sudo` is for operations that the chain itself initiate, i.e. native modules.It allows privilleged calls into the contract that no external user can do. We employ it for "ibc_lifecycle_complete", a method that only the IBC Middleware module is capable of calling.

## Process Timeout Flow
To handle the Timeout Event Returned when using ICS20 packet sending request contract callback. We override the **IBC Module**'s OnTimeoutPacket method.
- execute the **IBC Module**'s OnTimeoutPacket base method on this packet
- read module's store with `key` = `("%s::%d", Ics20Packet.GetSourceChannel(), Ics20Packet.GetSequence())` to get callback contract address
- check if callback contract address is valid Bench32 address
- make a call to `sudo` entry point of callback contract with message:
    `{"ibc_lifecycle_complete": {"ibc_timeout": {"channel": "%s", "sequence": %d}}},SourceChannel, Sequence))`
    - `SourceChannel`: ICS20 Packet's SourceChannel
    - `Sequence`: ICS20 Packet's Sequence
- delete callback data from module'store using `key` above

