package response

import (
	"bytes"
	"encoding/binary"
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/zeebo/bencode"
	"net"
)

type ResponseCompacted struct {
	Interval int    `bencode:"interval"`
	Peers4   []byte `bencode:"peers"`
	Peers6   []byte `bencode:"peers6"`
}

func (self *ResponseCompacted) Bencode() (string, error) {
	return bencode.EncodeString(self)
}

func (self *ResponseCompacted) Response() Response {
	response := Response{
		Interval: self.Interval,
		Peers:    self.Peers(),
	}
	return response
}

func (self *ResponseCompacted) Peers() []common.Peer {
	peers := make([]common.Peer, 0)
	peers = append(peers, self.peers4Parse()...)
	peers = append(peers, self.peers6Parse()...)
	return peers
}

func (self *ResponseCompacted) peers4Parse() []common.Peer {
	peers := make([]common.Peer, 0)
	buf := bytes.NewBuffer(self.Peers4)
	for true {
		ipBytes := make([]byte, 4)
		if n, err := buf.Read(ipBytes); err != nil || n != 4 {
			break
		}
		ipInt := binary.BigEndian.Uint32(ipBytes)
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, ipInt)
		portBytes := make([]byte, 2)
		if n, err := buf.Read(portBytes); err != nil || n != 2 {
			break
		}
		port := binary.BigEndian.Uint16(portBytes)
		peer := common.Peer{
			IP:   common.Address(ip.String()),
			Port: int(port),
		}
		peers = append(peers, peer)
	}
	return peers
}

func (self *ResponseCompacted) peers6Parse() []common.Peer {
	peers := make([]common.Peer, 0)
	buf := bytes.NewBuffer(self.Peers6)
	for true {
		ipBytes := make([]byte, 16)
		if n, err := buf.Read(ipBytes); err != nil || n != 4 {
			break
		}
		ipInt := binary.BigEndian.Uint32(ipBytes) //64?
		ip := make(net.IP, 16)
		binary.BigEndian.PutUint32(ip, ipInt)
		portBytes := make([]byte, 2)
		if n, err := buf.Read(portBytes); err != nil || n != 2 {
			break
		}
		port := binary.BigEndian.Uint16(portBytes)
		peer := common.Peer{
			IP:   common.Address(ip.String()),
			Port: int(port),
		}
		peers = append(peers, peer)
	}
	return peers
}

func (self *ResponseCompacted) ReloadPeers(peers []common.Peer) {
	peers4 := bytes.Buffer{}
	peers6 := bytes.Buffer{}
	for _, peer := range peers {
		if ip, err := peer.IP.IPv4(); err == nil { // if IPv4
			peers4.Write([]byte(ip)) // write IP to buf
			portBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(portBytes, uint16(peer.Port))
			peers4.Write(portBytes) // write port to buf
		} else if ip, err := peer.IP.IPv6(); err == nil { // if not IPv4 -> check for IPv6
			DebugLog.Printf("IPv6 peer in compacted mode: %s:%d\n", ip, peer.Port)
			peers6.Write([]byte(ip)) // write IP to buf
			portBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(portBytes, uint16(peer.Port))
			peers6.Write(portBytes) // write port to buf
		}
	}
	self.Peers4 = peers4.Bytes()
	self.Peers6 = peers6.Bytes()
}
