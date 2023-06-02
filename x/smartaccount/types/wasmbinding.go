package types

import (
	"encoding/json"
	"fmt"

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
	Value   []byte `json:"value"`
}

func ParseMessagesString(msgs []sdk.Msg) ([]MsgData, error) {
	msgsStr := make([]MsgData, 0)

	for index, msg := range msgs {
		msgData, err := json.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("error in json marshal msg %d: %s", index, err.Error())
		}

		data := MsgData{
			TypeURL: sdk.MsgTypeURL(msg),
			Value:   msgData,
		}

		msgsStr = append(msgsStr, data)
	}
	return msgsStr, nil
}
