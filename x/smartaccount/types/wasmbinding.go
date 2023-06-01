package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountMsg struct {
	ValidateTx     *ValidateTx     `json:"validate,omitempty"`
	AfterExecuteTx *AfterExecuteTx `json:"after_execute,omitempty"`
}

type ValidateTx struct {
	Msgs []MsgData `json:"msgs"`
}

type AfterExecuteTx struct {
	Msgs []MsgData `json:"msgs"`
}

type MsgData struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value"`
}

func ParseMessagesString(msgs []sdk.Msg) ([]MsgData, error) {
	var msgsStr []MsgData

	for index, msg := range msgs {
		// get msg type name
		splitStr := strings.Split(reflect.TypeOf(msg).String(), ".")
		msgType := splitStr[len(splitStr)-1]

		msgData, err := json.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("error in json marshal msg %d: %s", index, err.Error())
		}

		data := MsgData{
			TypeURL: msgType,
			Value:   string(msgData),
		}

		msgsStr = append(msgsStr, data)
	}
	return msgsStr, nil
}
