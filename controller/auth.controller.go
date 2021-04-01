package controller

import (
	"net/http"

	"github.com/fadhlimulyana20/go_backend/database"
	"github.com/fadhlimulyana20/go_backend/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

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

	result := db.Create(&user)

	if err := result.Error; err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return c.JSON(res.Status, res)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"firsName": user.FirstName,
		"lastName": user.LastName,
	})
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
