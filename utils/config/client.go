package config

import (
	"utils/logger"
	"github.com/go-ini/ini"
)

type Cconfig struct {
	ServerIp    string
	DefaultPort int
	ServerPort  int
}

func ClientConfig() *Cconfig {
	client := &Cconfig{}
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.Error(err)
	}
	cfgSec := cfg.Section("clinets")
	client.ServerIp = cfgSec.Key("server_ip").String()
	client.ServerPort, _ = cfgSec.Key("server_port").Int()
	client.DefaultPort, _ = cfgSec.Key("default_port").Int()
	return client
}
