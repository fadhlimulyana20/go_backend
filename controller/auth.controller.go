package controller

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/fadhlimulyana20/go_backend/config"
	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/fadhlimulyana20/go_backend/utils"
	"github.com/fadhlimulyana20/go_backend/utils/constant"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ctx = context.Background()

func Register(c echo.Context) error {
	// initialize variable
	db := database.GetConnection()
	var user models.User
	var res models.JsonResponse

	// Binding input data
	if err := c.Bind(&user); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// Validate Input data
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

	// Create new user record
	result := db.Create(&user)
	if err := result.Error; err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// Create Confirmation url
	cu := utils.ConfirmationUrl{}
	code, err := cu.Create(user.ID)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	// sendEmail
	if err := config.SendEmail(user.Email, "Password Confirmation", code); err != nil {
		log.Fatal(err)
	}

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
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["userId"] = user.ID
	sess.Save(c.Request(), c.Response())

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

func ConfirmUser(c echo.Context) error {
	// Initialize Varaible
	db := database.GetConnection()
	var res models.JsonResponse
	var user models.User
	token := c.Param("token")

	// If there's no token provided, then return bad request
	if token == "" {
		res.Status = http.StatusBadRequest
		res.Message = "Token Required!"
		return c.JSON(res.Status, res)
	}

	// Create Redis Connection
	rc := &config.RedisConfig{}
	rc.Init()
	rdb := rc.GetConnection()

	// Get userId from redis with token
	// Return token expired, if token key not found
	key := constant.ConfirmUserPrefix + token
	userId, err := rdb.Get(ctx, key).Result()
	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Token Expired"
		return c.JSON(res.Status, res)
	}

	// Update user record to validated
	db.Model(&user).Where("id", userId).Update("is_validated", true)

	// return a message
	res.Status = http.StatusOK
	res.Message = "You are now validated"
	return c.JSON(http.StatusOK, res)
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
