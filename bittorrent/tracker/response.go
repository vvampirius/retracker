package tracker

import (
	"../common"
	"github.com/zeebo/bencode"
)

type Response struct {
	Interval int `bencode:"interval"`
	Peers []common.Peer `bencode:"peers"`
}

func (self *Response) Bencode() (string, error) {
	return bencode.EncodeString(self)
}