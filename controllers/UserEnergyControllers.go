package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserEnergy(c *gin.Context) {
	db := config.DB
	email := c.Param("email")

	var energy users.UserEnergy
	if err := db.Where("email = ?", email).First(&energy).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Request Failed, no data found or query wrong",
			Data:    nil,
		})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"email":  energy.Email,
	// 	"energy": energy.Energy,
	// })

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Request successful",
		Data:    &energy,
	})
}

func CreateUserEnergy(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	message := repositories.CreateUserEnergyFirstTime(input.Email)

	if strings.Contains(strings.ToLower(message), "success") {
		c.JSON(http.StatusOK, message)
	} else {
		c.JSON(http.StatusInternalServerError, message)
	}
}

func UseEnergyForAll(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad parameter",
		})
		return
	}

	// chek apakah parameter berformat email
	_, err := mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad parameter",
		})
		return
	}

	if err := services.UseEnergy(email); err != nil {
		if err.Error() == "Energy 0, Please wait 10 minuts to recharge" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something wrong with server or email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Energy successfully used on students",
	})
}

func AddEnergyForAll(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad parameter",
		})
		return
	}

	// chek apakah parameter berformat email
	_, err := mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad parameter",
		})
		return
	}

	if err := services.AddEnergy(email); err != nil {
		if err.Error() == "energy is at maximum" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something wrong with server or email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Energy successfully used on students",
	})
}
