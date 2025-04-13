package controllers

import (
	"coba1BE/config"
	// "coba1BE/controllers"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	role := c.Param("role")
	var err error
	db := config.DB

	switch role {
	case "siswa":
		var siswa users.Siswa
		if err = c.ShouldBindJSON(&siswa); err == nil {
			err = db.Create(&siswa).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan nama " + siswa.Nama + " berhasil dibuat",
				})
				return
			}
		}
	case "guru":
		var guru users.Guru
		if err = c.ShouldBindJSON(&guru); err == nil {
			err = db.Create(&guru).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan nama " + guru.Nama + " berhasil dibuat",
				})
				return
			}
		}
	case "admin":
		var admin users.Admin
		if err = c.ShouldBindJSON(&admin); err == nil {
			err = db.Create(&admin).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan nama " + admin.Nama + " berhasil dibuat",
				})
				return
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	// check error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetSiswa(c *gin.Context) {
	var response models.BaseResponseModel

	response = repositories.GetAllSiswa()

	c.JSON(http.StatusOK, response)
}

func GetDataActor(c *gin.Context) {
	var response models.BaseResponseModel
	role := c.Param("role")

	// chek apakah parameter di isi
	if role == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}
	response = repositories.GetDataActor(role)

	c.JSON(http.StatusOK, response)
}

func GetUsers(c *gin.Context) {
	var response models.BaseResponseModel

	response = repositories.GetAllUsers()

	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context) {
	var response models.BaseResponseModel
	validRoles := map[string]bool{"siswa": true, "admin": true, "guru": true}
	// Ambil query parameter "role"
	role := c.Query("role")
	fmt.Println("role : " + role)

	// Jika role kosong, set default "siswa"
	if role == "" {
		role = "siswa"
	} else {
		// Pastikan role yang dikirim valid
		role = strings.ToLower(role) // Buat case insensitive
		if !validRoles[role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Gunakan: siswa, admin, guru"})
			return
		}
	}

	email, errmail := services.GetUserEmailFromToken(c)
	if errmail != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error di adminContollrs": errmail.Error()})
		return
	}
	fmt.Println(errmail)

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

	response = repositories.GetUserByEmail(email, role)

	c.JSON(http.StatusOK, response)
}

func UpdateDataActor(c *gin.Context) {
	// var result users.UserProfileReq
	var response models.BaseResponseModel

	role := c.Param("role")

	// chek apakah parameter di isi
	if role == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	switch role {
	// ### SISWA ###
	case "siswa":
		var siswa users.Siswa
		if err := c.ShouldBindJSON(&siswa); err != nil {
			response = models.BaseResponseModel{
				Message: err.Error(),
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		result, message := repositories.UpdateDataSiswa(siswa)
		if message != "Success" {
			response = models.BaseResponseModel{
				Message: message,
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response = models.BaseResponseModel{
			Message: message,
			Data:    result,
		}

		// ### GURU ###
	case "guru":
		var guru users.Guru
		if err := c.ShouldBindJSON(&guru); err != nil {
			response = models.BaseResponseModel{
				Message: err.Error(),
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		result, message := repositories.UpdateDataGuru(guru)
		if message != "Success" {
			response = models.BaseResponseModel{
				Message: message,
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response = models.BaseResponseModel{
			Message: message,
			Data:    result,
		}

		// ### ADMIN ###

	case "admin":
		var admin users.Admin
		if err := c.ShouldBindJSON(&admin); err != nil {
			response = models.BaseResponseModel{
				Message: err.Error(),
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		result, message := repositories.UpdateDataAdmin(admin)
		if message != "Success" {
			response = models.BaseResponseModel{
				Message: message,
				Data:    nil,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response = models.BaseResponseModel{
			Message: message,
			Data:    result,
		}

		// ### DEVAULT ###
	default:
		response = models.BaseResponseModel{
			Message: "undifine role",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
