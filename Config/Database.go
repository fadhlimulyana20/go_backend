package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Database Config
type DatabaseConfig struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"3306"`
	Database string `env:"DB_NAME"`
}

func GetDatabaseConfig() *DatabaseConfig {
	c := DatabaseConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &c
}

func (c *DatabaseConfig) GetMySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
}
