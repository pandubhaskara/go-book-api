package router

import (
	"go-book-api/controller"
	"go-book-api/service"

	"github.com/labstack/echo/v5"
)

var (
	authController controller.AuthController = controller.AuthController{}
	bookController controller.BookController = controller.BookController{}
)

func UseSubroute(e *echo.Echo) {
	e.POST("/auth/token", authController.Token)

	books := e.Group("/books", service.ValidateToken)

	books.GET("", bookController.Index)
	books.GET("/:id", bookController.Detail)
	books.POST("", bookController.Create)
	books.PUT("/:id", bookController.Update)
	books.DELETE("/:id", bookController.Delete)
}
