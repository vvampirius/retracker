package main

type Core struct {
	Config   *Config
	Storage  *Storage
	Receiver *Receiver
}

func NewCore(config *Config, tempStorage *TempStorage) *Core {
	storage := NewStorage(config)
	core := Core{
		Config:   config,
		Storage:  storage,
		Receiver: NewReceiver(config, storage),
	}
	core.Receiver.Announce.TempStorage = tempStorage
	return &core
}
