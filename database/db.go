package database

import (
	"fmt"
	"log"

	"github.com/fadhlimulyana20/go_backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error
var db *gorm.DB

func Init() {
	c := &config.DatabaseConfig{}

	c.GetDatabaseConfig()

	if c.Driver == "mysql" {
		dsn := c.GetMySqlDSN()
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			fmt.Printf("Cannot connect to %s database", c.Driver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("Database connected succesfully")
		}
	}

	if c.Driver == "postgres" {
		dsn := c.GetPostgresDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			fmt.Printf("Cannot connect to %s database", c.Driver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("Database connected succesfully")
		}
	}
}

func GetConnection() *gorm.DB {
	return db
}
