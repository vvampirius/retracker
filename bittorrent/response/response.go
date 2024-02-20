package response

import (
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/zeebo/bencode"
	"fmt"
)

type Response struct {
	Interval int `bencode:"interval"`
	Peers []common.Peer `bencode:"peers"`
}

func (self *Response) Bencode(compacted bool) (string, error) {
	if compacted {
		response := self.Compacted()
		return response.Bencode()
	}
	return bencode.EncodeString(self)
}

func (self *Response) Compacted() ResponseCompacted {
	response := ResponseCompacted{
		Interval: self.Interval,
	}
	response.ReloadPeers(self.Peers)
	return response
}

func Load(b []byte) (*Response, error) {
	response := Response{}
	if err := bencode.DecodeBytes(b, &response); err!=nil {
		responseCompacted := ResponseCompacted{}
		if err := bencode.DecodeBytes(b, &responseCompacted); err==nil {
			response = responseCompacted.Response()
			fmt.Println()
		} else { return nil, err }
	}
	return &response, nil
}

