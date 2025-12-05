package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"projek/config"
	"projek/middleware"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// OPTIONAL: register user baru buat testing
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	now := time.Now()

	err := config.DB.QueryRow(
		`INSERT INTO users (username, password, created_at, created_by, modified_at, modified_by)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id`,
		req.Username, req.Password,
		now, "system",
		now, "system",
	).Scan(new(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var id int
	var password string
	err := config.DB.QueryRow(
		"SELECT id, password FROM users WHERE username = $1",
		req.Username,
	).Scan(&id, &password)

	if err == sql.ErrNoRows || password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}

	token, err := middleware.GenerateToken(id, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
