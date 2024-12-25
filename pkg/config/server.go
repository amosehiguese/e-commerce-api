package config

import "github.com/amosehiguese/ecommerce-api/pkg/utils"

type serverConfig struct {
	Port string
}

func setServerConfig() *serverConfig {
	var s serverConfig
	utils.MustMapEnv(&s.Port, "SERVER_PORT")
	return &s
}
