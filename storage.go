package main

import (
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/vvampirius/retracker/bittorrent/tracker"
	"sync"
	"time"
)

type Storage struct {
	Config     *Config
	Requests   map[common.InfoHash]map[common.PeerID]tracker.Request
	requestsMu sync.Mutex
}

func (self *Storage) Update(request tracker.Request) {
	self.requestsMu.Lock()
	defer self.requestsMu.Unlock()
	if _, ok := self.Requests[request.InfoHash]; !ok {
		self.Requests[request.InfoHash] = make(map[common.PeerID]tracker.Request)
	}
	self.Requests[request.InfoHash][request.PeerID] = request

}

func (self *Storage) Delete(request tracker.Request) {
	self.requestsMu.Lock()
	defer self.requestsMu.Unlock()
	delete(self.Requests[request.InfoHash], request.PeerID) //TODO: test this
}

func (self *Storage) GetPeers(infoHash common.InfoHash) []common.Peer {
	self.requestsMu.Lock()
	defer self.requestsMu.Unlock()
	peers := make([]common.Peer, 0)
	if requests, ok := self.Requests[infoHash]; ok {
		for _, request := range requests {
			peers = append(peers, request.Peer())
		}
	}
	return peers
}

func (self *Storage) purgeRoutine() {
	for true {
		time.Sleep(1 * time.Minute)
		if self.Config.Debug {
			DebugLog.Printf("In memory %d hashes\n", len(self.Requests))
			DebugLog.Println(`Locking...`)
		}
		self.requestsMu.Lock()
		for hash, requests := range self.Requests {
			if self.Config.Debug {
				DebugLog.Printf("%d peer in hash %x\n", len(requests), hash)
			}
			for peerId, request := range requests {
				timestampDelta := request.TimeStampDelta()
				if self.Config.Debug {
					DebugLog.Printf(" %x %s:%d %v\n", peerId, request.Peer().IP, request.Peer().Port, timestampDelta)
				}
				if timestampDelta > self.Config.Age {
					DebugLog.Printf("delete peer %x in hash %x\n", peerId, hash)
					delete(self.Requests[hash], peerId)
				}
			}
			if len(requests) == 0 {
				DebugLog.Printf("delete hash %x\n", hash)
				delete(self.Requests, hash)
			}
		}
		self.requestsMu.Unlock()
		if self.Config.Debug {
			DebugLog.Println(`Unlocked`)
		}
	}
}

func NewStorage(config *Config) *Storage {
	storage := Storage{
		Config:   config,
		Requests: make(map[common.InfoHash]map[common.PeerID]tracker.Request),
	}
	go storage.purgeRoutine()
	return &storage
}
