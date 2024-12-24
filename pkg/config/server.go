package config

import "github.com/amosehiguese/ecommerce-api/pkg/utils"

type serverConfig struct {
	Port    string
	Address string
}

func setServerConfig() *serverConfig {
	var s serverConfig
	utils.MustMapEnv(&s.Address, "SERVER_ADDR")
	utils.MustMapEnv(&s.Port, "SERVER_PORT")
	return &s
}
