package ante

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UserOps struct {
	Messages string `json:"messages"`
}

type ValidateUserOps struct {
	Validate UserOps `json:"validate"`
}

type PreExecuteUserOps struct {
	PreExecute UserOps `json:"pre_execute"`
}

type ValidateUserOpsResponse = bool

type MsgData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func ParseMessagesString(msgs []sdk.Msg) ([]MsgData, error) {
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
