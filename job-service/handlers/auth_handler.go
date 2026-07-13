package handlers

import (
	"net/http"
	"time"

	"github.com/aakash811/queueflow/job-service/auth"
	"github.com/aakash811/queueflow/job-service/db"
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

	var user struct {
		id           string
		passwordHash string
		role         string
	}

	err = db.DB.QueryRow(
		c.Request.Context(),
		"SELECT id, password_hash, role FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.id, &user.passwordHash, &user.role)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid credentials",
			},
		)

		return
	}

	if !auth.CheckPasswordHash(req.Password, user.passwordHash) {
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
		user.role,
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
			"expires_at": time.Now().Add(24 * time.Hour).Unix(),
		},
	)
}