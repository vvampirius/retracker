package announce

import (
	"net/http"
	"regexp"
	"fmt"
)

func (self *Announce) HttpHandler(w http.ResponseWriter, r *http.Request) {
	xrealip := r.Header.Get(`X-Real-IP`)
	if self.Logger != nil {
		self.Logger.Printf("%s %s %s %s\n", r.RemoteAddr, xrealip, r.RequestURI, r.UserAgent())
	}
	rr := self.ProcessAnnounce(
		self.getRemoteAddr(r, xrealip),
		r.URL.Query().Get(`info_hash`),
		r.URL.Query().Get(`peer_id`),
		r.URL.Query().Get(`port`),
		r.URL.Query().Get(`uploaded`),
		r.URL.Query().Get(`downloaded`),
		r.URL.Query().Get(`left`),
		r.URL.Query().Get(`ip`),
		r.URL.Query().Get(`numwant`),
		r.URL.Query().Get(`event`),
		)
	if d, err := rr.Bencode(); err==nil {
		fmt.Fprint(w, d)
		if self.Logger != nil && self.Config.Debug {
			self.Logger.Printf("Bencode: %s\n", d)
		}
	} else { self.Logger.Println(err.Error()) }
}

func (self *Announce) getRemoteAddr(r *http.Request, xrealip string) string {
	if self.Config.XRealIP && xrealip!=`` {
		return xrealip
	}
	return self.parseRemoteAddr(r.RemoteAddr, `127.0.0.1`)
}

func (self *Announce) parseRemoteAddr(in, def string) string {
	address := def
	r := regexp.MustCompile(`(.*):\d+$`)
	if match := r.FindStringSubmatch(in); len(match)==2 { address = match[1] }
	return address
}