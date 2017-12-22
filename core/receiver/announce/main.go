package announce

import (
	"log"
	"os"
	Storage "../../storage"
	CoreCommon "../../common"
)

type Announce struct {
	Config *CoreCommon.Config
	Logger *log.Logger
	Storage *Storage.Storage
}

func New(config *CoreCommon.Config, storage *Storage.Storage) *Announce {
	announce := Announce{
		Config: config,
		Logger: log.New(os.Stdout, `announce `, log.Flags()),
		Storage: storage,
	}
	return &announce
}