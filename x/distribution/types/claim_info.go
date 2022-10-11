package types

import (
	"time"
)

type ClaimInfo struct {
	ClaimTime     time.Time `json:"claim_time"`
	ClaimBlockNum int64     `json:"claim_block_num"`
	Address       string    `json:"address"`
}

func (c ClaimInfo) Reset() {}

func (c ClaimInfo) String() string {
	//TODO implement me
	panic("implement me")
}

func (c ClaimInfo) ProtoMessage() {}

func (c ClaimInfo) IsEmpty() bool {
	return c == ClaimInfo{}
}
