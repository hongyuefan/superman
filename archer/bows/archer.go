package bows

import (
	"fmt"
)

type ArcherCmd struct {
	Exchange  string
	BusiId    int32
	ReqSerial uint32
	PayLoad   []byte
}

func (c *ArcherCmd) ToString() string {
	return fmt.Sprintf("exchange:%s,business:%v,reqSerial:%v,payload:%s", c.Exchange, c.BusiId, c.ReqSerial, string(c.PayLoad))
}

type Archer interface {
	Init() error
	Run(chan *ArcherCmd)
}
