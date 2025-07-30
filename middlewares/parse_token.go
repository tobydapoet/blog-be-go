package middlewares

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	AccountID uint   `json:"ID"`
	Type      string `json:"Type"`
	jwt.RegisteredClaims
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil {
		fmt.Println("Token parse error:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims or token not valid")
}
