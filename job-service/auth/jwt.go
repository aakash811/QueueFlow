package auth

import (
	"time"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	username string,
	role string,
) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role": role,
		"exp": time.Now().Add(
			24 * time.Hour,
		).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(config.AppConfig.JWTSecret),
	)
}