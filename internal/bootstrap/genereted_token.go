package bootstrap

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateLocalToken(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "service_role",
		"iss":  "supabase",
		"exp":  time.Now().Add(time.Hour * 24 * 365).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
