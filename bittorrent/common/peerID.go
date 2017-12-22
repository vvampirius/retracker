package common

type PeerID string

func (self *PeerID) Valid() bool {
	if len(*self) == 20 { return true }
	return false
}