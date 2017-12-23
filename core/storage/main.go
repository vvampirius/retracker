package storage

import (
	CoreCommon "github.com/vvampirius/retracker/core/common"
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/vvampirius/retracker/bittorrent/tracker"
	"time"
	"log"
	"os"
)

type Storage struct {
	Config *CoreCommon.Config
	Requests map[common.InfoHash]map[common.PeerID]tracker.Request
	Logger *log.Logger
}

func (self *Storage) Update(request tracker.Request) {
	if _, ok := self.Requests[request.InfoHash]; !ok {
		self.Requests[request.InfoHash] = make(map[common.PeerID]tracker.Request)
	}
	self.Requests[request.InfoHash][request.PeerID] = request

}

func (self *Storage) Delete(request tracker.Request) {
	delete(self.Requests[request.InfoHash], request.PeerID) //TODO: test this
}

func (self *Storage) GetPeers(infoHash common.InfoHash) []common.Peer {
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
		logger := *self.Logger
		logger.SetPrefix(`purgeRoutine() `)
		if self.Config.Debug { logger.Printf("In memory %d hashes\n", len(self.Requests)) }
		for hash, requests := range self.Requests {
			if self.Config.Debug { logger.Printf("%d peer in hash %x\n",len(requests), hash) }
			for peerId, request := range requests {
				timestampDelta := request.TimeStampDelta()
				if self.Config.Debug { logger.Printf(" %x %s:%d %v\n", peerId, request.Peer().IP, request.Peer().Port, timestampDelta) }
				if timestampDelta > self.Config.Age {
					logger.Printf("delete peer %x in hash %x\n", peerId, hash)
					delete(self.Requests[hash], peerId)
				}
			}
			if len(requests) == 0 {
				logger.Printf("delete hash %x\n", hash)
				delete(self.Requests, hash)
			}
		}
	}
}


func New(config *CoreCommon.Config) *Storage {
	logger := log.New(os.Stdout, `storage# `, log.Flags())
	storage := Storage{
		Config: config,
		Requests: make(map[common.InfoHash]map[common.PeerID]tracker.Request),
		Logger: logger,
	}
	go storage.purgeRoutine()
	return &storage
}