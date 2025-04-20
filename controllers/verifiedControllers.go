package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Beingverifieds(c *gin.Context) {
	response, message := repositories.GetSiswaBeingVerified()

	if strings.Contains(strings.ToLower(message), "success") {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: message,
			Data:    response,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
		return
	}
}

func Beingverified(c *gin.Context) {
	email, errmail := services.GetUserEmailFromToken(c)
	if errmail != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error di adminContollrs": errmail.Error()})
		return
	}
	fmt.Println(email)

	// chek apakah parameter di isi
	if email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	// chek apakah parameter berformat email
	_, err := mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "bad parameter",
			Data:    nil,
		})
		return
	}

	response, message := repositories.GetSiswaBeingVerifiedByEmail(email)

	if strings.Contains(strings.ToLower(message), "success") {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: message,
			Data:    response,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
		return
	}
}

func UpdateUserVerifiedBatch(c *gin.Context) {
	var input []users.UserVerified
	db := config.DB

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	for _, item := range input {
		var existing users.UserVerified

		// Cek apakah datanya ada
		err := db.First(&existing, "email = ?", item.Email).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Data tidak ditemukan
			c.JSON(http.StatusNotFound, models.BaseResponseModel{
				Message: fmt.Sprintf("User not found: %s", item.Email),
				Data:    nil,
			})
			return
		} else if err != nil {
			// Error lain
			c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
				Message: fmt.Sprintf("User not found: %s", item.Email),
				Data:    nil,
			})
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}

		// Kalau ditemukan, baru cek perlu update atau tidak
		if existing.Verified != item.Verified {
			err := db.Model(&users.UserVerified{}).
				Where("email = ?", item.Email).
				Update("verified", item.Verified).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
					Message: fmt.Sprintf("Failed to update: %s", item.Email),
					Data:    nil,
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Update successful",
		Data:    input,
	})
}
