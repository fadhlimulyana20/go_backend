package database

import (
	"fmt"

	"github.com/fadhlimulyana20/go_backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error
var db *gorm.DB

func Init() {
	c := config.GetDatabaseConfig()

	dsn := c.GetMySqlDSN()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected succesfully")
}

func GetConnection() *gorm.DB {
	return db
}
