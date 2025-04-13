package services

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("error form parseToken : " + err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("error Claims token")
		return nil, err
	}

	// Convert jwt.MapClaims to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}

	return result, nil
}
