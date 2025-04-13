package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func GenerateToken(username string) (string, string, error) {
	// Access Token (Berlaku 15 menit)
	accessTokenClaims := jwt.MapClaims{
		"username": username,
		// "exp":      time.Now().Add(15 * time.Minute).Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (Berlaku 7 hari)
	refreshTokenClaims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

// ================= VERIS 1 ==================
// ValidateToken memvalidasi token JWT
// func ValidateToken(tokenString string) (*jwt.Token, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }

// ================= VERIS 2 ==================
// ValidateToken memvalidasi token JWT dan memeriksa apakah token valid
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse token dan verifikasi dengan secretKey
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi metode signing yang digunakan pada token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	// Jika ada error saat parsing token
	if err != nil {
		fmt.Println("error saat parsing token : " + err.Error())
		return nil, err
	}

	// Pastikan token valid (cek expiration dan signature)
	if !token.Valid {
		fmt.Println("cek expiration dan signature")
		return nil, errors.New("invalid token")
	}

	return token, nil
}
