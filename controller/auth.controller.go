package controller

import (
	"errors"
	"net/http"

	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/fadhlimulyana20/go_backend/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	db := database.GetConnection()
	var user models.User
	var res models.JsonResponse

	if err := c.Bind(&user); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	if err := c.Validate(&user); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// Create Hashed Password
	hashedPassword, err := utils.HashPassword(user.Password)
	// If there's an Error, return internal Server Error
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}
	// Else, assign user password with hashed Password
	user.Password = hashedPassword

	result := db.Create(&user)

	if err := result.Error; err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// confirmationCode := uuid.New().String()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"firsName": user.FirstName,
		"lastName": user.LastName,
	})
}

func Login(c echo.Context) error {
	// initialize Variable
	db := database.GetConnection()
	var user models.User
	var login loginDTO
	var res models.JsonResponse

	// Binding data from user input, return error if there's an error
	if err := c.Bind(&login); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// Else, find a user with email = input email
	result := db.Where("email = ?", login.Email).First(&user)

	// Return error, if there's an error
	if err := result.Error; err != nil {
		// Check error type
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Status = http.StatusNotFound
		} else {
			res.Status = http.StatusInternalServerError
		}

		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// If there's no error, compare input password with hashed password stored in database.
	// return UnAuthorized, if not valid
	if valid := utils.ComparePassword(login.Password, user.Password); !valid {
		res.Status = http.StatusUnauthorized
		res.Message = "Wrong Password"
		return c.JSON(res.Status, res)
	}

	// Return a user and create a session if everything is OK.
	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"firsName":    user.FirstName,
		"lastName":    user.LastName,
		"isValidated": user.IsValidated,
	}
	return c.JSON(res.Status, res)

}

func Auth(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["foo"] = "bar"
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func Me(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if _, ok := sess.Values["foo"]; ok {
		return c.JSON(http.StatusOK, map[string]interface{}{"message": sess.Values["foo"]})
	}

	return c.NoContent(http.StatusOK)
}
