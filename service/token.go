package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// Define a super secret key. In production, load this from environment variables.
var jwtKey = []byte("your_super_secret_and_secure_key")

// Define custom claims struct, embedding jwt.RegisteredClaims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	// Set the token expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	// Create the claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Set the "exp" claim
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
		},
	}

	// Declare the token with the algorithm used for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to get the complete encoded token string
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		tokenString := strings.Split(authorization, " ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		c.Set("username", claims.Username)

		return next(c)
	}
}
