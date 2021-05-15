package controller

import (
	"context"
	"errors"
	"fmt"
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
var db = database.GetConnection()

func Register(c echo.Context) error {
	// initialize variable
	user := &models.User{}
	var res models.JsonResponse

	// Binding input data
	if err := c.Bind(&user); err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Validate Input data
	if err := c.Validate(&user); err != nil {
		res.Status = constant.StatusFail
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// Create Hashed Password
	err := user.BeforeSave(db)
	// If there's an Error, return internal Server Error
	if err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Create new user record
	newUser, err := user.SaveUser(db)
	if err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Create Confirmation url
	cu := utils.ConfirmationUrl{}
	code, err := cu.Create(newUser.ID)
	if err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// sendEmail
	if err := config.SendEmail(newUser.Email, "Password Confirmation", code); err != nil {
		log.Fatal(err)
	}

	res.Status = "success"
	res.Data = map[string]interface{}{
		"id":          newUser.ID,
		"email":       newUser.Email,
		"firstName":   newUser.FirstName,
		"lastName":    newUser.LastName,
		"isValidated": newUser.IsValidated,
		"createdAt":   newUser.CreatedAt,
		"updatedAt":   newUser.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
}

func Login(c echo.Context) error {
	// initialize Variable
	var user models.User
	var login loginDTO
	var res models.JsonResponse

	// Binding data from user input, return error if there's an error
	if err := c.Bind(&login); err != nil {
		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Else, find a user with email = input email
	result := db.Where("email = ?", login.Email).First(&user)

	// Return error, if there's an error
	if err := result.Error; err != nil {
		// Check error type
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Status = constant.StatusFail
			res.Message = err.Error()
			return c.JSON(http.StatusNotFound, res)
		}

		res.Status = constant.StatusError
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	// If there's no error, compare input password with hashed password stored in database.
	// return UnAuthorized, if not valid
	if valid := utils.ComparePassword(login.Password, user.Password); !valid {
		res.Status = constant.StatusFail
		res.Message = "Wrong Password"
		return c.JSON(http.StatusUnauthorized, res)
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

	res.Status = constant.StatusSuccess
	res.Message = fmt.Sprintf("Welcome %d", user.ID)
	res.Data = map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"firsName":    user.FirstName,
		"lastName":    user.LastName,
		"isValidated": user.IsValidated,
		"createdAt":   user.CreatedAt,
		"updatedAt":   user.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)

}

func ConfirmUser(c echo.Context) error {
	// Initialize Varaible
	var res models.JsonResponse
	var user models.User
	token := c.Param("token")

	// If there's no token provided, then return bad request
	if token == "" {
		res.Status = constant.StatusFail
		res.Message = "Token Required!"
		return c.JSON(http.StatusBadRequest, res)
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
		res.Status = constant.StatusFail
		res.Message = "Token Expired"
		return c.JSON(http.StatusBadRequest, res)
	}

	// Update user record to validated
	db.Model(&user).Where("id", userId).Update("is_validated", true)

	// return a message
	res.Status = constant.StatusSuccess
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
