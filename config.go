package main

import (
	"github.com/vvampirius/retracker/common"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	AnnounceResponseInterval int
	Listen                   string
	Debug                    bool
	Age                      float64
	XRealIP                  bool
	Forwards                 []common.Forward
	ForwardTimeout           int
}

func (config *Config) ReloadForwards(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	defer f.Close()
	forwards := make([]common.Forward, 0)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&forwards); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	config.Forwards = forwards
	return nil
}
