package main

import (
	"encoding/hex"
	"errors"
	"github.com/vvampirius/retracker/bittorrent/common"
	"github.com/zeebo/bencode"
	"net/http"
	"strings"
)

var (
	ErrNoInfohashes = errors.New(`no infohashes found`)
	ErrBadInfohash  = errors.New(`bad infohash`)
)

type ScrapeResponseHash struct {
	Complete   int `bencode:"complete"`
	Incomplete int `bencode:"incomplete"`
	Downloaded int `bencode:"downloaded"`
}

type ScrapeResponse struct {
	Files map[string]ScrapeResponseHash `bencode:"files"`
}

func (core *Core) getScrapeResponse(infoHashes []string) (ScrapeResponse, error) {
	if infoHashes == nil || len(infoHashes) == 0 {
		ErrorLog.Println(ErrNoInfohashes.Error())
		return ScrapeResponse{}, ErrNoInfohashes
	}
	scrapeResponse := ScrapeResponse{
		Files: make(map[string]ScrapeResponseHash),
	}
	core.Storage.requestsMu.Lock()
	defer core.Storage.requestsMu.Unlock()
	for _, infoHashString := range infoHashes {
		infoHash := common.InfoHash(infoHashString)
		infoHashHex := strings.ToUpper(hex.EncodeToString([]byte(infoHashString)))
		if !infoHash.Valid() {
			ErrorLog.Println(infoHashString, ErrBadInfohash.Error())
			return ScrapeResponse{}, ErrBadInfohash
		}
		requestInfoHash, found := core.Storage.Requests[infoHash]
		if !found {
			DebugLog.Printf("%s not found in storage", infoHashHex)
			continue
		}
		srh := ScrapeResponseHash{}
		for _, peerRequest := range requestInfoHash {
			if peerRequest.Event == `competed` || peerRequest.Left == 0 {
				srh.Complete++
			} else {
				srh.Incomplete++
			}
		}
		DebugLog.Printf("%s\tComplete (seed): %d\tIncomplete (leech): %d", infoHashHex, srh.Complete, srh.Incomplete)
		srh.Downloaded = srh.Complete // Unfortunately, we do not collect statistics to present the actual value.

		scrapeResponse.Files[infoHashString] = srh
	}
	return scrapeResponse, nil
}

func (core *Core) httpScrapeHandler(w http.ResponseWriter, r *http.Request) {
	xrealip := r.Header.Get(`X-Real-IP`)
	DebugLog.Printf("%s %s %s '%s' '%s'\n", r.Method, r.RemoteAddr, xrealip, r.RequestURI, r.UserAgent())
	query := r.URL.Query()
	scrapeResponse, err := core.getScrapeResponse(query["info_hash"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encoder := bencode.NewEncoder(w)
	if err := encoder.Encode(scrapeResponse); err != nil {
		ErrorLog.Println(err.Error())
	}
}
