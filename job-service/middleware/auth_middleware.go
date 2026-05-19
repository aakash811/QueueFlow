package middleware

import (
	"net/http"
	"strings"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader(
			"Authorization",
		)

		if authHeader == "" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Missing Token",
				},
			)

			c.Abort()

			return
		}

		parts := strings.Split(
			authHeader,
			" ",
		)

		if len(parts) != 2 ||
			parts[0] != "Bearer" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Invalid Authorization Header",
				},
			)

			c.Abort()

			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,

			func(token *jwt.Token) (interface{}, error) {

				return []byte(
					config.AppConfig.JWTSecret,
				), nil
			},
		)

		if err != nil || !token.Valid {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Invalid Token",
				},
			)

			c.Abort()

			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set(
			"role",
			claims["role"],
		)

		c.Set(
			"username",
			claims["username"],
		)

		c.Next()
	}
}
