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
	Msgs      []Any     `json:"msgs"`
	CallInfor CallInfor `json:"call_info"`
	IsAuthz   bool      `json:"is_authz"`
}

type AfterExecuteTx struct {
	Msgs      []Any     `json:"msgs"`
	CallInfor CallInfor `json:"call_info"`
	IsAuthz   bool      `json:"is_authz"`
}

type Any struct {
	TypeURL string `json:"type_url"`
	Value   []byte `json:"value"`
}

type CallInfor struct {
	Fee        sdk.Coins `json:"fee"`
	Gas        uint64    `json:"gas"`
	FeePayer   string    `json:"fee_payer"`
	FeeGranter string    `json:"fee_granter"`
}

func ParseMessagesString(msgs []sdk.Msg) ([]Any, error) {
	msgsStr := make([]Any, 0, len(msgs))

	for _, msg := range msgs {
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}

		data := Any{
			TypeURL: "/" + proto.MessageName(msg),
			Value:   msgBytes,
		}

		msgsStr = append(msgsStr, data)
	}
	return msgsStr, nil
}
