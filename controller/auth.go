package controller

import (
	"fmt"
	"go-book-api/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AuthController struct{}

func (ctrl *AuthController) Token(c *echo.Context) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	if body.Username == "admin" && body.Password == "password" {
		token, err := service.GenerateJWT(body.Username)
		if err != nil {
			fmt.Printf("Error generating token: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		fmt.Printf("Generated Token: %s\n", token)

		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}

	return c.NoContent(http.StatusUnauthorized)
}
