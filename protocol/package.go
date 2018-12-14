package protocol

import (
	"encoding/json"
)

type Package interface {
	ParseFromArray([]byte) bool
	SerialToArray() []byte
	GetBusinessId() int32
	GetReqSerial() uint32
	GetPayload() []byte
}

type FixPackage struct {
	BusiId    int32  `json:"busiId"`
	ReqSerial uint32 `json:"reqId"`
	Payload   []byte `json:"payload"`
}

func (t *FixPackage) GetBusinessId() int32 {
	return t.BusiId
}

func (t *FixPackage) GetReqSerial() uint32 {
	return t.ReqSerial
}

func (t *FixPackage) GetPayload() []byte {
	return t.Payload
}

func (t *FixPackage) ParseFromArray(data []byte) error {
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	return nil
}

func (t *FixPackage) SerialToArray() []byte {
	byt, _ := json.Marshal(t)
	return byt
}
