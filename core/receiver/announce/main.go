package announce

import (
	"log"
	"os"
	Storage "github.com/vvampirius/retracker/core/storage"
	CoreCommon "github.com/vvampirius/retracker/core/common"
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