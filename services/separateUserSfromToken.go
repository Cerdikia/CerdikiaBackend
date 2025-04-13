package services

import (
	"fmt"

	// "github.com/dgrijalva/jwt-go"

	// "net/http"
	"github.com/gin-gonic/gin"
)

// func GetUserIDFromToken(c *gin.Context) (string, error) {
// 	// Mendapatkan token dari header Authorization
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		return "", fmt.Errorf("no token provided")
// 	}

// 	// Validasi dan parsing token
// 	token, err := ValidateToken(tokenString)
// 	if err != nil || !token.Valid {
// 		return "", fmt.Errorf("invalid or expired token")
// 	}

// 	// Ambil klaim dari token (misalnya user_id)
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return "", fmt.Errorf("invalid token claims")
// 	}

// 	// Ambil user_id dari klaim token
// 	userID, ok := claims["username"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("user ID not found in token")
// 	}

// 	return userID, nil
// }

func GetUserEmailFromToken(c *gin.Context) (string, error) {
	// claims, ok := c.MustGet("claims").(map[string]interface{})
	claims, ok := c.MustGet("username").(map[string]interface{})

	if !ok {
		return "", fmt.Errorf("invalid username")
	}

	// username := claims["username"].(string)
	username := claims["username"].(string)
	fmt.Println("token username from separateUser : " + username)
	// role := claims["role"].(string)

	return username, nil
}
