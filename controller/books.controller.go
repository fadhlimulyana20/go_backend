package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Method to get collection of book
func GetAllBooks(c echo.Context) error {
	db := database.GetConnection()
	var books []models.Book
	var res models.JsonResponse

	result := db.Find(&books)

	if result.Error != nil {
		res.Status = http.StatusInternalServerError
		fmt.Println(result.Error.Error())
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = books

	return c.JSON(http.StatusOK, res)
}

// Method to get single object of Book
func GetBook(c echo.Context) error {
	db := database.GetConnection()
	id := c.Param("id")
	var book models.Book
	var res models.JsonResponse

	result := db.First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			res.Status = http.StatusNotFound
			res.Message = fmt.Sprintf("Book with id %s is not found", id)
			return c.JSON(http.StatusNotFound, res)
		}
		res.Status = http.StatusInternalServerError
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = book

	return c.JSON(http.StatusOK, res)
}
