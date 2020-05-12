package config

import (
	"utils/logger"

	"github.com/go-ini/ini"
)

type Sconfig struct {
	DefaultPort int
}

func ServerConfig() *Sconfig {
	server := &Sconfig{}
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.Error(err)
	}
	cfgSec := cfg.Section("server")
	server.DefaultPort, _ = cfgSec.Key("default_port").Int()
	return server
}
