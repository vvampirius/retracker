package common

type InfoHash string

func (self *InfoHash) Valid() bool {
	if len(*self) == 20 { return true }
	return false
}