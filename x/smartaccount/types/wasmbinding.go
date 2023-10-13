package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type AccountMsg struct {
	PreExecuteTx   *PreExecuteTx   `json:"pre_execute,omitempty"`
	AfterExecuteTx *AfterExecuteTx `json:"after_execute,omitempty"`
	RecoverTx      *RecoverTx      `json:"recover,omitempty"`
}

type RecoverTx struct {
	Caller      string `json:"caller"`
	PubKey      []byte `json:"pub_key"`
	Credentials []byte `json:"credentials"`
}

type PreExecuteTx struct {
	Msgs      []MsgData `json:"msgs"`
	CallInfor CallInfor `json:"call_info"`
	IsAuthz   bool      `json:"is_authz"`
}

type AfterExecuteTx struct {
	Msgs      []MsgData `json:"msgs"`
	CallInfor CallInfor `json:"call_info"`
	IsAuthz   bool      `json:"is_authz"`
}

type MsgData struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value"`
}

type CallInfor struct {
	Fee        sdk.Coins `json:"fee"`
	Gas        uint64    `json:"gas"`
	FeePayer   string    `json:"fee_payer"`
	FeeGranter string    `json:"fee_granter"`
}

func ParseMessagesString(msgs []sdk.Msg) ([]MsgData, error) {
	msgsStr := make([]MsgData, 0)

	for _, msg := range msgs {
		msgData, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}

		data := MsgData{
			TypeURL: proto.MessageName(msg),
			Value:   string(msgData),
		}

		msgsStr = append(msgsStr, data)
	}
	return msgsStr, nil
}
