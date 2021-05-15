package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/fadhlimulyana20/go_backend/utils/constant"
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
		res.Status = constant.StatusFail
		fmt.Println(result.Error.Error())
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = constant.StatusSuccess
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
			res.Status = constant.StatusFail
			res.Message = fmt.Sprintf("Book with id %s is not found", id)
			return c.JSON(http.StatusNotFound, res)
		}
		res.Status = constant.StatusError
		res.Message = result.Error.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = constant.StatusSuccess
	res.Data = book

	return c.JSON(http.StatusOK, res)
}

func CreateBook(c echo.Context) error {
	db := database.GetConnection()
	var res models.JsonResponse
	b := &models.CreateBookDTO{}

	// Validate data binding,
	if err := c.Bind(b); err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Validate Field, if field not valid then return error
	if err := c.Validate(b); err != nil {
		res.Status = constant.StatusFail
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// else if data and field valid, create object to datavase
	book := models.Book{
		Title:       b.Title,
		Description: b.Description,
	}
	result := db.Create(&book)

	// If there is any error when creating record in database, return error
	if result.Error != nil {
		res.Status = constant.StatusError
		res.Message = result.Error.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// if everithing are ok, return data that has been created.
	res.Status = constant.StatusSuccess
	res.Data = book

	return c.JSON(http.StatusOK, res)
}

func ModifyBook(c echo.Context) error {
	db := database.GetConnection()
	id := c.Param("id")
	var res models.JsonResponse
	b := &models.CreateBookDTO{}

	// Validate data binding,
	if err := c.Bind(b); err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Validate Field, if field not valid then return error
	if err := c.Validate(b); err != nil {
		res.Status = constant.StatusFail
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// else if data and field valid, create object to datavase
	book := models.Book{}
	result := db.First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			res.Status = constant.StatusFail
			res.Message = fmt.Sprintf("Book with id %s is not found", id)
			return c.JSON(http.StatusNotFound, res)
		}
		res.Status = constant.StatusError
		res.Message = result.Error.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	book.Title = b.Title
	book.Description = b.Description
	db.Save(&book)

	// if everithing are ok, return data that has been created.
	res.Status = constant.StatusSuccess
	res.Data = book

	return c.JSON(http.StatusOK, res)
}

func DeleteBook(c echo.Context) error {
	db := database.GetConnection()
	id := c.Param("id")
	var book models.Book
	var res models.JsonResponse

	result := db.Delete(&book, id)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Status = constant.StatusFail
			res.Message = err.Error()
			return c.JSON(http.StatusNotFound, res)
		}

		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = constant.StatusSuccess
	res.Message = fmt.Sprintf("Book with id %s is successfully deleted", id)
	return c.JSON(http.StatusOK, res)
}
