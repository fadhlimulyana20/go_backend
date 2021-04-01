package main

import (
	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/fadhlimulyana20/go_backend/routes"
)

func main() {
	// Create db connection
	database.Init()
	db := database.GetConnection()

	db.AutoMigrate(&models.Book{}, &models.User{})

	// Initialize Routes
	e := routes.Init()
	e.Logger.Fatal(e.Start(":5000"))
}
