package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID uint, role string, email string, avatar string, name string) (string, error) {
	claims := jwt.MapClaims{
		"ID":     userID,
		"Role":   role,
		"Email":  email,
		"Name":   name,
		"Avatar": avatar,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}
