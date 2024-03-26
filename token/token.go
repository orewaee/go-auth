package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.MapClaims, key string) string {
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))

	return token
}

func VerifyToken(token, key string) bool {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if err != nil {
		return false
	}

	return parsed.Valid
}
