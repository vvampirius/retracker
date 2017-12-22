package receiver

import (
	Announce "./announce"
	Storage "../storage"
	CoreCommon "../common"
)

type Receiver struct {
	Announce *Announce.Announce
}

func New(config *CoreCommon.Config, storage *Storage.Storage) *Receiver {
	receiver := Receiver{
		Announce: Announce.New(config, storage),
	}
	return &receiver
}