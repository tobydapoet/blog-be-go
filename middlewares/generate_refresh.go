package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"ID":   userID,
		"Type": "refresh",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}
