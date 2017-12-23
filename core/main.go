package core

import (
	"net/http"
	"fmt"
	"github.com/vvampirius/retracker/core/common"
	Receiver "github.com/vvampirius/retracker/core/receiver"
	Storage "github.com/vvampirius/retracker/core/storage"
)

type Core struct {
	Config *common.Config
	Storage *Storage.Storage
	Receiver *Receiver.Receiver
}

func New(config *common.Config) *Core {
	storage := Storage.New(config)
	core := Core{
		Config: config,
		Storage: storage,
		Receiver: Receiver.New(config, storage),
	}
	http.HandleFunc("/announce", core.Receiver.Announce.HttpHandler)
	if err := http.ListenAndServe(config.Listen, nil); err != nil { // set listen port
		fmt.Println(err)
	}
	//TODO: do it with context
	return &core
}