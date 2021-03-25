package Config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env"
	"github.com/go-pg/pg/v10"
)

// PostgresConfig persists the config for our PostgreSQL database connection
type PostgresConfig struct {
	URL      string `env:"DATABASE_URL"` // DATABASE_URL will be used in preference if it exists
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

// GetConnection returns our pg database connection
// usage:
// db := config.GetConnection()
// defer db.Close()
func GetConnection() *pg.DB {
	c := GetPostgresConfig()

	// if DATABASE_URL is valid, we will use its constituent values to preferences
	validConfig, err := validPostgresURL(c.URL)
	if err == nil {
		c = validConfig
	}

	db := pg.Connect(&pg.Options{
		Addr:     ":" + c.Port,
		User:     c.User,
		Password: c.Password,
		Database: c.Database,
	})

	return db
}

// GetPostgresConfig returns a PostgresConfig pointer with the correct Postgres Config values
func GetPostgresConfig() *PostgresConfig {
	c := PostgresConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &c
}

func validPostgresURL(URL string) (*PostgresConfig, error) {
	if URL == "" || strings.TrimSpace(URL) == "" {
		return nil, errors.New("database url is blank")
	}

	validURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	c := &PostgresConfig{}
	c.URL = URL
	c.Host = validURL.Host
	c.Database = strings.Replace(validURL.Path, "/", "", 1)
	c.Port = validURL.Port()
	c.User = validURL.User.Username()
	c.Password, _ = validURL.User.Password()
	return c, nil
}
