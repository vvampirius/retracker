package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvampirius/retracker/bittorrent/common"
	Response "github.com/vvampirius/retracker/bittorrent/response"
	"github.com/vvampirius/retracker/bittorrent/tracker"
	CoreCommon "github.com/vvampirius/retracker/common"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type ReceiverAnnounce struct {
	Config      *Config
	Storage     *Storage
	Prometheus  *Prometheus
	TempStorage *TempStorage
}

func (ra *ReceiverAnnounce) httpHandler(w http.ResponseWriter, r *http.Request) {
	if ra.Prometheus != nil {
		ra.Prometheus.Requests.Inc()
	}
	xrealip := r.Header.Get(`X-Real-IP`)
	DebugLog.Printf("%s %s %s '%s' '%s'\n", r.Method, r.RemoteAddr, xrealip, r.RequestURI, r.UserAgent())
	remoteAddr := ra.getRemoteAddr(r, xrealip)
	remotePort := r.URL.Query().Get(`port`)
	infoHash := r.URL.Query().Get(`info_hash`)
	if ra.Config.Debug {
		DebugLog.Printf("hash: '%x', remote addr: %s:%s", infoHash, remoteAddr, remotePort)
	}
	response := ra.ProcessAnnounce(
		remoteAddr,
		infoHash,
		r.URL.Query().Get(`peer_id`),
		remotePort,
		r.URL.Query().Get(`uploaded`),
		r.URL.Query().Get(`downloaded`),
		r.URL.Query().Get(`left`),
		r.URL.Query().Get(`ip`),
		r.URL.Query().Get(`numwant`),
		r.URL.Query().Get(`event`),
	)
	compacted := false
	if r.URL.Query().Get(`compact`) == `1` {
		compacted = true
	}
	d, err := response.Bencode(compacted)
	if err != nil {
		ErrorLog.Println(err.Error())
		return
	}
	fmt.Fprint(w, d)
	/*
		if ra.Config.Debug {
			DebugLog.Printf("Bencode: %s\n", d)
		}
	*/
}

func (ra *ReceiverAnnounce) getRemoteAddr(r *http.Request, xrealip string) string {
	if ra.Config.XRealIP && xrealip != `` {
		return xrealip
	}
	return ra.parseRemoteAddr(r.RemoteAddr, `127.0.0.1`)
}

func (ra *ReceiverAnnounce) parseRemoteAddr(in, def string) string {
	address := def
	r := regexp.MustCompile(`(.*):\d+$`)
	if match := r.FindStringSubmatch(in); len(match) == 2 {
		address = match[1]
	}
	return address
}

func (ra *ReceiverAnnounce) ProcessAnnounce(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
	event string) *Response.Response {
	if request, err := tracker.MakeRequest(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
		event, DebugLog); err == nil {

		response := Response.Response{
			Interval: 30,
		}

		if request.Event != `stopped` {
			ra.Storage.Update(*request)
			response.Peers = ra.Storage.GetPeers(request.InfoHash)
			response.Peers = append(response.Peers, ra.makeForwards(*request)...)
		} else {
			ra.Storage.Delete(*request)
		}

		return &response
	}

	return nil
}

func (ra *ReceiverAnnounce) makeForwards(request tracker.Request) []common.Peer {
	peers := make([]common.Peer, 0)
	forwardsCount := len(ra.Config.Forwards)
	if forwardsCount > 0 {
		ch := make(chan []common.Peer, forwardsCount)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(ra.Config.ForwardTimeout))
		for _, v := range ra.Config.Forwards {
			go ra.makeForward(v, request, ch, ctx)
		}
		for i := 0; i < forwardsCount; i++ {
			peers = append(peers, <-ch...)
		}
		peers = CoreCommon.PeersUniq(peers)
		if ra.Config.Debug {
			DebugLog.Printf("%x has uniq peers: %d\n", request.InfoHash, len(peers))
		}
	}
	return peers
}

func (ra *ReceiverAnnounce) makeForward(forward CoreCommon.Forward, request tracker.Request, ch chan<- []common.Peer, ctx context.Context) {
	peers := make([]common.Peer, 0)
	uri := fmt.Sprintf("%s?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d", forward.Uri, url.QueryEscape(string(request.InfoHash)),
		url.QueryEscape(string(request.PeerID)), request.Port, request.Uploaded, request.Downloaded, request.Left)
	if forward.Ip != `` {
		uri = fmt.Sprintf("%s&ip=%s&ipv4=%s", uri, forward.Ip, forward.Ip) //TODO: check for IPv4
	}
	hash := fmt.Sprintf("%x", request.InfoHash)
	forwardName := forward.GetName()
	if ra.Config.Debug {
		if forward.Ip != `` {
			DebugLog.Printf("Announce %x to %s with IP %s", hash, forwardName, forward.Ip)
		} else {
			DebugLog.Printf("Announce %x to %s", hash, forwardName)
		}
		//DebugLog.Println(uri)
	}

	rqst, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		ErrorLog.Println(err)
		ch <- peers
		return
	}
	client := http.Client{}
	response, err := client.Do(rqst)
	if err != nil {
		ErrorLog.Printf("Announce %x to %s got error: %s", hash, forwardName, err.Error())
		if ra.Prometheus != nil {
			ra.Prometheus.ForwarderStatus.With(prometheus.Labels{`name`: forwardName, `status`: `error`}).Inc()
		}
		ch <- peers
		return
	}
	defer response.Body.Close()
	if ra.Prometheus != nil {
		ra.Prometheus.ForwarderStatus.With(prometheus.Labels{`name`: forwardName, `status`: fmt.Sprintf("%d", response.StatusCode)}).Inc()
	}
	if response.StatusCode != http.StatusOK {
		ErrorLog.Printf("Announce %x to %s got status: %s", request.InfoHash, forward.GetName(), response.Status)
		ch <- peers
		return
	}
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLog.Printf("Announce %x from %s read error: %s", request.InfoHash, forward.GetName(), err.Error())
		if ra.Prometheus != nil {
			ra.Prometheus.ForwarderStatus.With(prometheus.Labels{`name`: forwardName, `status`: fmt.Sprintf("%d", response.StatusCode)}).Inc()
		}
		ch <- peers
		return
	}
	tempFilename := ``
	if ra.Config.Debug {
		tempFilename = ra.TempStorage.SaveBencodeFromForwarder(payload, fmt.Sprintf("%x", request.InfoHash), uri)
	}
	bitResponse, err := Response.Load(payload)
	if err != nil {
		ErrorLog.Printf("Announce %x from %s parse error: %s", request.InfoHash, forward.GetName(), err.Error())
		if tempFilename == `` {
			tempFilename = ra.TempStorage.SaveBencodeFromForwarder(payload, fmt.Sprintf("%x", request.InfoHash), uri)
		}
		ErrorLog.Fatalln(`Check file`, tempFilename)
		if ra.Prometheus != nil {
			ra.Prometheus.ForwarderStatus.With(prometheus.Labels{`name`: forwardName, `status`: fmt.Sprintf("%d", response.StatusCode)}).Inc()
		}
		ch <- peers
		return
	}
	if ra.Config.Debug {
		DebugLog.Printf("Announce %x to %s got %d peers", request.InfoHash, forward.GetName(), len(bitResponse.Peers))
	}
	peers = append(peers, bitResponse.Peers...)
	ch <- peers
}

func NewReceiverAnnounce(config *Config, storage *Storage) *ReceiverAnnounce {
	announce := ReceiverAnnounce{
		Config:  config,
		Storage: storage,
	}
	return &announce
}
