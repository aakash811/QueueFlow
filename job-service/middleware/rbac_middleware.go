package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(
	allowedRoles ...string,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exists := ctx.Get("role")

		if !exists {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "role missing",
				},
			)
			ctx.Abort()
			return
		}

		role := roleValue.(string)

		for _, allowed := range allowedRoles {
			if role == allowed {
				ctx.Next()
				return
			}
		}

		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "access denied",
			},
		)
		ctx.Abort()
	}
}