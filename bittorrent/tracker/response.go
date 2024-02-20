package tracker

import (
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/zeebo/bencode"
)

type Response struct {
	Interval int `bencode:"interval"`
	Peers []common.Peer `bencode:"peers"`
	//Peers []byte `bencode:"peers"`
}

func (self *Response) Bencode() (string, error) {
	return bencode.EncodeString(self)
}

func (self *Response) New(bencoded string) (*Response, error) {
	r := Response{}
	if err := bencode.DecodeString(bencoded, &r); err!=nil {
		return nil, err
	}
	return &r, nil
}