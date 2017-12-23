package receiver

import (
	Announce "github.com/vvampirius/retracker/core/receiver/announce"
	Storage "github.com/vvampirius/retracker/core/storage"
	CoreCommon "github.com/vvampirius/retracker/core/common"
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