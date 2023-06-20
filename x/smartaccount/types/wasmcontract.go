package types

import (
	"crypto/sha512"
	"encoding/json"
)

type InstantiateSalt struct {
	Owner   string `json:"owner"`
	CodeID  uint64 `json:"code_id"`
	InitMsg []byte `json:"init_msg"`
	PubKey  []byte `json:"pub_key"`
}

func GenerateSalt(owner string, codeId uint64, initMsg []byte, pubKey []byte) ([]byte, error) {
	salt := InstantiateSalt{
		Owner:   owner,
		CodeID:  codeId,
		InitMsg: initMsg,
		PubKey:  pubKey,
	}

	salt_bytes, err := json.Marshal(salt)
	if err != nil {
		return nil, err
	}

	// instantiate2 salt max length is 64 bytes, so need hash here
	salt_hashed := sha512.Sum512(salt_bytes)

	return salt_hashed[:], nil
}
