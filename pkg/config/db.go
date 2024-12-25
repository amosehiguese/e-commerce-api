package config

import (
	"fmt"
	"os"

	"github.com/amosehiguese/ecommerce-api/pkg/utils"
)

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SslMode  string
}

func setDatabaseConfig() *databaseConfig {
	var d databaseConfig
	utils.MustMapEnv(&d.Host, "DB_HOST")
	utils.MustMapEnv(&d.User, "DB_USER")
	utils.MustMapEnv(&d.Password, "DB_PASSWORD")
	utils.MustMapEnv(&d.SslMode, "DB_SSLMODE")
	d.Port = utils.GetEnvAsInt("DB_PORT")
	d.Name = os.Getenv("DB_NAME")

	return &d
}

func (d *databaseConfig) ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.Name, d.SslMode)
}

func (d *databaseConfig) ConnStringDefaultDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.SslMode)
}
