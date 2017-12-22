package common

type Peer struct {
	PeerID PeerID `bencode:"peer_id"`
	IP Address `bencode:"ip"`
	Port int `bencode:"port"`
}
