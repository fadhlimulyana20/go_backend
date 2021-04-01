package routes

import (
	"net/http"

	"github.com/fadhlimulyana20/go_backend/controller"
	"github.com/fadhlimulyana20/go_backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	// Create HTTP server using echo
	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// Custom Validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	auth := e.Group("/auth")
	auth.GET("", controller.Auth)
	auth.GET("/me", controller.Me)
	auth.POST("/register", controller.Register)

	book := e.Group("/book")
	book.GET("", controller.GetAllBooks)
	book.GET("/:id", controller.GetBook)
	book.POST("/", controller.CreateBook)
	book.DELETE("/:id", controller.DeleteBook)
	book.PUT("/:id", controller.ModifyBook)

	return e
}
