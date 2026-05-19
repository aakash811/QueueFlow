package handlers

import (
	"net/http"

	"github.com/aakash811/queueflow/job-service/auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)

		return
	}

	if req.Username != "admin" || req.Password != "password" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid credentials",
			},
		)

		return
	}

	token, err := auth.GenerateToken(
		req.Username,
		"ADMIN",
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "token generation failed",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token": token,
		},
	)
}