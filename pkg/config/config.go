package config

import "github.com/amosehiguese/ecommerce-api/pkg/utils"

type Config struct {
	Env      string
	Server   *serverConfig
	Database *databaseConfig
	JWT      *jwtConfig
}

var c Config

func initConfig() *Config {
	c.Server = setServerConfig()
	c.Database = setDatabaseConfig()
	c.JWT = setJwtConfig()
	utils.MustMapEnv(&c.Env, "ECOMM_ENV")

	return &c
}

func Get() *Config {
	if c.Server == nil || c.Database == nil {
		c = *initConfig()
	}
	return &c
}
