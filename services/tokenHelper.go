package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ExtractEmailFromToken extracts the email (username) from a JWT token string
func ExtractEmailFromToken(tokenString string) (string, error) {
	// Get the secret key from environment
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key" // Fallback secret key
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if username exists in claims
		if username, exists := claims["username"]; exists {
			// Convert to string
			if usernameStr, ok := username.(string); ok {
				return usernameStr, nil
			}
			return "", fmt.Errorf("username is not a string")
		}
		return "", fmt.Errorf("username not found in token claims")
	}

	return "", fmt.Errorf("invalid token")
}

// ExtractEmailFromAuthHeader extracts the email from an Authorization header
func ExtractEmailFromAuthHeader(authHeader string) (string, error) {
	// Check if the header is in the correct format
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// Extract the token
	tokenString := authHeader[7:] // Remove "Bearer " prefix
	return ExtractEmailFromToken(tokenString)
}
