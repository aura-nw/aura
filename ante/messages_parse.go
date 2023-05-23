package ante

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

/*
{"sender":"aura14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9swserkw","contract":"aura14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9swserkw","msg":{"pre_validate":{"messages":"{\"bank\":{\"send\":{\"to_address\":\"aura19ecqv8ga40jrpsltetnafj7lazll0mwtk2q5h3\",\"amount\":[{\"denom\":\"uaura\",\"amount\":\"500\"}]}}}"}},"funds":[]}
{"from_address":"aura14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9swserkw","to_address":"aura19ecqv8ga40jrpsltetnafj7lazll0mwtk2q5h3","amount":[{"denom":"uaura","amount":"500"}]}
*/

type MsgType string

const (
	Bank MsgType = "bank"
	Wasm MsgType = "wasm"
)

type WasmTypeMsg string

const (
	WasmExecute                 WasmTypeMsg = "execute"
	WasmClearAdmin              WasmTypeMsg = "clear_admin"
	WasmUpdateAdmin             WasmTypeMsg = "update_admin"
	WasmIBCCloseChannel         WasmTypeMsg = "ibc_close_channel"
	WasmIBCSend                 WasmTypeMsg = "ibc_send"
	WasmInstantiateContract     WasmTypeMsg = "instantiate_contract"
	WasmInstantiateContract2    WasmTypeMsg = "instantiate_contract_2"
	WasmMigrateContract         WasmTypeMsg = "migrate_contract"
	WasmStoreCode               WasmTypeMsg = "store_code"
	WasmUpdateInstantiateConfig WasmTypeMsg = "update_instantiate_config"
)

type BankTypeMsg string

const (
	BankSend      BankTypeMsg = "send"
	BankMultiSend BankTypeMsg = "multi_send"
)

type MsgData struct {
	Type    MsgType `json:"type"`
	SubType string  `json:"sub_type"`
	Data    string  `json:"data"`
}

// just test with 2 type of message
func parseMessagesString(msgs []sdk.Msg) []MsgData {

	var msgsTypeArray []MsgData

	for _, msg := range msgs {
		switch msg := msg.(type) {
		// wasm types
		case *wasmtypes.MsgExecuteContract:
			msgStr, _ := json.Marshal(msg)

			msgWasm := MsgData{
				Type:    Wasm,
				SubType: string(WasmExecute),
				Data:    string(msgStr),
			}

			msgsTypeArray = append(msgsTypeArray, msgWasm)
		case *wasmtypes.MsgClearAdmin:
			break
		case *wasmtypes.MsgUpdateAdmin:
			break
		case *wasmtypes.MsgIBCCloseChannel:
			break
		case *wasmtypes.MsgIBCSend:
			break
		case *wasmtypes.MsgInstantiateContract:
			break
		case *wasmtypes.MsgInstantiateContract2:
			break
		case *wasmtypes.MsgMigrateContract:
			break
		case *wasmtypes.MsgStoreCode:
			break
		case *wasmtypes.MsgUpdateInstantiateConfig:
			break

		// bank types
		case *banktypes.MsgSend:
			msgStr, _ := json.Marshal(msg)

			msgBank := MsgData{
				Type:    Bank,
				SubType: string(BankSend),
				Data:    string(msgStr),
			}

			msgsTypeArray = append(msgsTypeArray, msgBank)
		case *banktypes.MsgMultiSend:
			break
		default:
			break
		}
	}

	return msgsTypeArray
}
