package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	// Create HTTP server using echo
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	return e
}
