package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"
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
		fmt.Println(message)
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
		return
	}
}

func Beingverified(c *gin.Context) {
	email := c.Query("email")
	// email, errmail := services.GetUserEmailFromToken(c)
	// if errmail != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error di adminContollrs": errmail.Error()})
	// 	return
	// }
	fmt.Println(email)

	// chek apakah parameter di isi
	if email == "" || email == " " {
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
	type VerificationRequest struct {
		Students []struct {
			Email              string `json:"email"`
			VerificationStatus string `json:"verification_status"`
		} `json:"students"`
	}

	var request VerificationRequest
	db := config.DB

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	fmt.Println("student.Email : " + request.Students[0].Email)
	fmt.Println("student.VerificationStatus : " + request.Students[0].VerificationStatus)

	var updatedStudents []map[string]any

	for _, student := range request.Students {
		var existing users.UserVerified

		// Validate verification status
		if student.VerificationStatus != "accept" && student.VerificationStatus != "rejected" && student.VerificationStatus != "waiting" {
			fmt.Println("Validate verification status error")
			fmt.Println("Validate status : " + student.VerificationStatus)
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: fmt.Sprintf("Invalid verification status for student %s: %s. Must be 'accept', 'rejected', or 'waiting'", student.Email, student.VerificationStatus),
				Data:    nil,
			})
			return
		}

		// Check if the student exists in the verified table
		err := db.First(&existing, "email = ?", student.Email).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If not found, check if the student exists in the siswa table
			var siswa struct {
				Email string
			}
			err = db.Table("siswa").Where("email = ?", student.Email).First(&siswa).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, models.BaseResponseModel{
					Message: fmt.Sprintf("Student not found: %s", student.Email),
					Data:    nil,
				})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
					Message: fmt.Sprintf("Database error: %s", err.Error()),
					Data:    nil,
				})
				return
			}

			// Create new verification record
			newVerification := users.UserVerified{
				Email:          student.Email,
				VerifiedStatus: student.VerificationStatus,
			}
			err = db.Create(&newVerification).Error
			if err != nil {
				fmt.Println("Create new verification record error")
				c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
					Message: fmt.Sprintf("Failed to create verification record: %s", err.Error()),
					Data:    nil,
				})
				return
			}
		} else if err != nil {
			// Other database error
			c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
				Message: fmt.Sprintf("Database error: %s", err.Error()),
				Data:    nil,
			})
			return
		} else {
			// Update existing record
			fmt.Println("student.Email : " + student.Email)
			fmt.Println("student.VerificationStatus : " + student.VerificationStatus)
			err := db.Model(&users.UserVerified{}).
				Where("email = ?", student.Email).
				Update("verified_status", student.VerificationStatus).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
					Message: fmt.Sprintf("Failed to update: %s", err.Error()),
					Data:    nil,
				})
				return
			}
		}

		// Add to updated students list
		updatedStudents = append(updatedStudents, map[string]any{
			"email":           student.Email,
			"verified_status": student.VerificationStatus,
		})
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Verification status updated successfully",
		Data:    updatedStudents,
	})
}
