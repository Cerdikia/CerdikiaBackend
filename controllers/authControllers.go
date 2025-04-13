package controllers

import (
	// "coba1BE/controllers"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login handler
func LoginControler(c *gin.Context) {
	var user users.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid request",
			Data:    nil,
		})
		return
	}

	fmt.Println("email from auth controlers : " + user.Email)
	userStatus := repositories.GetUserByEmail(user.Email, user.Role).Message
	if userStatus != "Data retrieved successfully" {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		c.JSON(http.StatusUnauthorized, models.BaseResponseModel{
			Message: "Invalid credentials",
			Data:    nil,
		})
		return
	}

	accessToken, refreshToken, err := services.GenerateToken(user.Email)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Token generation failed",
			Data:    nil,
		})
		return
	}

	// c.JSON(http.StatusOK, users.LoginResponse{
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// })

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "E-mail " + user.Email + " Successfuly Login",
		Data: users.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})

}

// Refresh Token handler
func RefreshToken(c *gin.Context) {
	var req users.LoginResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := services.ValidateToken(req.RefreshToken)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token data"})
		return
	}

	// username, err := getDataToken(c)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// Generate new access token
	accessToken, refreshToken, err := services.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "E-mail " + username + " Successfuly Login",
		Data: users.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})
}

// func getDataToken(c *gin.Context) (string, error) {
// 	var req users.LoginResponse
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return "", fmt.Errorf("Invalid request")
// 	}

// 	token, err := services.ValidateToken(req.RefreshToken)
// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
// 		return "", fmt.Errorf("Invalid or expired refresh token")
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims"})
// 		return "", fmt.Errorf("Invalid refresh token claims")
// 	}

// 	username, ok := claims["username"].(string)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token data"})
// 		return "", fmt.Errorf("Invalid refresh token data")
// 	}
// 	return username, nil
// }
