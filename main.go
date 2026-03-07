package main

import (
	"fmt"
	"go-book-api/config"
	"go-book-api/helper"
	"io"
	"net/http"
	_ "time/tzdata"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	helper.Logger.Infof("Activating: %s\n", config.App.Name)

	e := echo.New()

	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"success": true})
	})

	e.POST("echo", func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.Blob(http.StatusOK, "application/json", body)
	})

	if err := e.Start(fmt.Sprintf("%s:%s", config.App.Host, config.App.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
