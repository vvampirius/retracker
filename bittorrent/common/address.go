package common

type Address string

func (self *Address) Valid() bool {
	if len(*self) == 0 { return false }
	return true
}
