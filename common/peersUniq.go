package common

import "github.com/vvampirius/retracker/bittorrent/common"

func PeerInPeers(peers []common.Peer, peer common.Peer) bool {
	for _, p := range peers {
		if p.IP == peer.IP && p.Port == peer.Port {
			return true
		}
	}
	return false
}

func PeersUniq(peers []common.Peer) []common.Peer {
	prs := make([]common.Peer, 0)
	for _, peer := range peers {
		if !PeerInPeers(prs, peer) {
			prs = append(prs, peer)
		}
	}
	return prs
}
