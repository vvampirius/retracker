package common

import "net/url"

type Forward struct {
	Name string
	Uri  string
	//TODO:
	//Proxy string
	Ip   string
	Host string
	//Dns []string
}

func (forward *Forward) GetName() string {
	if forward.Name != `` {
		return forward.Name
	}
	u, err := url.Parse(forward.Uri)
	if err == nil && u.Host != `` {
		return u.Host
	}
	return forward.Uri
}
