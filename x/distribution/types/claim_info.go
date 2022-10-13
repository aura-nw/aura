package types

import (
	"fmt"
	"time"
)

type ClaimInfo struct {
	ClaimTime     time.Time `json:"claim_time"`
	ClaimBlockNum int64     `json:"claim_block_num"`
}

func (c ClaimInfo) Reset() {}

func (c ClaimInfo) String() string {
	return fmt.Sprintf("claim_time: %s - claim_block_num: %d", c.ClaimTime.String(), c.ClaimBlockNum)
}

func (c ClaimInfo) ProtoMessage() {}

func (c ClaimInfo) IsEmpty() bool {
	return c == ClaimInfo{}
}
