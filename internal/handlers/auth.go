package handlers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad request format"})
		return
	}

	// simple auth check using login credentials from .env, I could implement a full login & register system if needed.
	if req.Username != os.Getenv("ADMIN_USER") || req.Password != os.Getenv("ADMIN_PASSWORD") {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	// generate JWT
	secret := os.Getenv("JWT_SECRET")
	var expMinutes int = 60 // default expiration time in minutes
	if val, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION")); err == nil {
		expMinutes = int(val.Minutes())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp": time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}