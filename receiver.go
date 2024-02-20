package main

type Receiver struct {
	Announce *ReceiverAnnounce
}

func NewReceiver(config *Config, storage *Storage) *Receiver {
	receiver := Receiver{
		Announce: NewReceiverAnnounce(config, storage),
	}
	return &receiver
}
