package common

import (
	"net"
	"errors"
)

type Address string

func (self *Address) Valid() bool {
	if len(*self) == 0 { return false }
	return true
}

func (self *Address) IPv4() (net.IP, error) {
	if ip := net.ParseIP(self.String()); ip!=nil {
		if ip = ip.To4(); ip!=nil {
			return ip, nil
		}
	}
	return nil, errors.New(`Can't convert to IPv4'`)
}

func (self *Address) IPv6() (net.IP, error) {
	if ip := net.ParseIP(self.String()); ip!=nil {
		if ip = ip.To16(); ip!=nil {
			return ip, nil
		}
	}
	return nil, errors.New(`Can't convert to IPv6'`)
}

func (self *Address) String() string {
	return string(*self)
}