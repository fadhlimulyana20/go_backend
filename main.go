package main

import (
	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/routes"
)

func main() {
	// Create db connection
	database.Init()

	// Initialize Routes
	e := routes.Init()
	e.Logger.Fatal(e.Start(":5000"))
}
