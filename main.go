package main

import (
	"context"
	"fmt"

	"github.com/fadhlimulyana20/go_backend/config"
	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/fadhlimulyana20/go_backend/routes"
)

var ctx = context.Background()

func main() {
	// Create db connection
	database.Init()
	db := database.GetConnection()

	db.AutoMigrate(&models.Book{}, &models.User{})

	// Connect to redis
	rc := &config.RedisConfig{}
	rc.Init()
	rdb := rc.GetConnection()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Redis connected successfully")
	}

	// Initialize Routes
	e := routes.Init()
	e.Logger.Fatal(e.Start(":5000"))
}
