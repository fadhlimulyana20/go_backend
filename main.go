package main

import (
	"context"

	"github.com/fadhlimulyana20/go_backend/Config"
	"github.com/gin-gonic/gin"
)

var err error

func main() {
	r := gin.Default()

	db := Config.GetConnection()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	route(r)

	r.Run()
}
