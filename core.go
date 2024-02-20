package main

type Core struct {
	Config   *Config
	Storage  *Storage
	Receiver *Receiver
}

func NewCore(config *Config) *Core {
	storage := NewStorage(config)
	core := Core{
		Config:   config,
		Storage:  storage,
		Receiver: NewReceiver(config, storage),
	}
	return &core
}
