package ante

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// just test with 2 type of message
func parseMessagesString(msgs []sdk.Msg) ([]MsgData, error) {
	var msgsArray []MsgData

	for index, msg := range msgs {
		splitStr := strings.Split(reflect.TypeOf(msg).String(), ".")
		msgType := splitStr[len(splitStr)-1]

		msgData, err := json.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("error in json marshal msg %d: %s", index, err.Error())
		}

		data := MsgData{
			Type: msgType,
			Data: string(msgData),
		}

		msgsArray = append(msgsArray, data)
	}

	return msgsArray, nil
}
