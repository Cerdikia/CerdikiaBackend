package controllers

import (
	// "coba1BE/controllers"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login handler
func Login(c *gin.Context) {
	var user users.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validasi user dari hardcoded data
	// if storedPassword, ok := users[user.Username]; !ok || storedPassword != user.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// 	return
	// }

	fmt.Println("email : " + user.Email)
	userStatus := repositories.GetUserByEmail(user.Email).Message
	if userStatus != "Data retrieved successfully" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := services.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, users.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
