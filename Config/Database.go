package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Database Config
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Driver   string
}

func (c *DatabaseConfig) GetDatabaseConfig() {
	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c.User = os.Getenv("DB_USER")
	c.Password = os.Getenv("DB_PASSWORD")
	c.Host = os.Getenv("DB_HOST")
	c.Port = os.Getenv("DB_PORT")
	c.Database = os.Getenv("DB_NAME")
	c.Driver = os.Getenv("DB_DRIVER")
}

func (c *DatabaseConfig) GetMySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
}

func (c *DatabaseConfig) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", c.Host, c.User, c.Password, c.Database, c.Port)
}
