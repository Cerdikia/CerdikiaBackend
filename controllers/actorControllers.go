package controllers

import (
	"coba1BE/config"
	"strconv"

	// "coba1BE/controllers"
	"coba1BE/models"
	"coba1BE/models/points"
	"coba1BE/models/users"
	"coba1BE/repositories"
	"coba1BE/services"
	"fmt"
	"net/http"
	"net/mail"
	"os"
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
			// Set default profile image if not provided
			if siswa.ImageProfile == "" {
				basuUrl := os.Getenv("BASEURL")
				siswa.ImageProfile = fmt.Sprintf("%s/uploads/default_user.png", basuUrl)
			}

			err = db.Create(&siswa).Error
			if err == nil {
				errPoint := repositories.CreatePointFirst(siswa.Email)
				errVerivied := repositories.CreateAcountVerifiedFirst(siswa.Email)
				message := repositories.CreateUserEnergyFirstTime(siswa.Email)
				if errPoint == nil || errVerivied == nil || strings.Contains(strings.ToLower(message), "success") {
					c.JSON(201, gin.H{
						"message": "User dengan email " + siswa.Email + " berhasil dibuat",
					})
					return
				} else if errPoint != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error di errPoint": errPoint.Error()})
					return
				} else if errVerivied != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error di errVerivied": errVerivied.Error()})
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error di CreateUserEnergyFirstTime": message})
					return
				}
			}
		}
	case "guru":
		var guru users.Guru
		if err = c.ShouldBindJSON(&guru); err == nil {
			if guru.ImageProfile == "" {
				basuUrl := os.Getenv("BASEURL")
				guru.ImageProfile = fmt.Sprintf("%s/uploads/default_user.png", basuUrl)
			}
			err = db.Create(&guru).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan email " + guru.Email + " berhasil dibuat",
				})
				return
			}
		}
	case "kepalaSekolah":
		var guru users.Guru
		if err = c.ShouldBindJSON(&guru); err == nil {
			if guru.ImageProfile == "" {
				basuUrl := os.Getenv("BASEURL")
				guru.ImageProfile = fmt.Sprintf("%s/uploads/default_user.png", basuUrl)
			}

			guru.Jabatan = "kepala sekolah"

			err = db.Create(&guru).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan email " + guru.Email + " berhasil dibuat",
				})
				return
			}
		}
	case "admin":
		var admin users.Admin
		if err = c.ShouldBindJSON(&admin); err == nil {
			if admin.ImageProfile == "" {
				basuUrl := os.Getenv("BASEURL")
				admin.ImageProfile = fmt.Sprintf("%s/uploads/default_user.png", basuUrl)
			}
			err = db.Create(&admin).Error
			if err == nil {
				c.JSON(201, gin.H{
					"message": "User dengan email " + admin.Email + " berhasil dibuat",
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
	response := repositories.GetAllSiswa()

	c.JSON(http.StatusOK, response)
}

func GetDataActor(c *gin.Context) {
	var response models.BaseResponseModel
	role := c.Param("role")
	email := c.Param("email")

	// chek apakah parameter di isi
	if role == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	if email != "" {
		response = repositories.GetDataActorByRoleAndEmail(role, email)
	} else {
		response = repositories.GetDataActor(role)
	}

	c.JSON(http.StatusOK, response)
}

func GetUsers(c *gin.Context) {
	response := repositories.GetAllUsers()

	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context) {
	// var response models.BaseResponseModel
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

	response, message := repositories.GetUserByEmail(email, role)

	if strings.Contains(strings.ToLower(message), "error bad request") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
		return
	}

	if strings.Contains(strings.ToLower(message), "no data found") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "error : no data found",
			Data:    nil,
		})
		return
	}

	if strings.Contains(strings.ToLower(message), "successfully") {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Data retrieved successfully",
			Data:    response,
		})
		return
	}
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

func DeleteDataActor(c *gin.Context) {
	var user users.LoginRequest
	db := config.DB

	if err := c.ShouldBindJSON(&user); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid request",
			Data:    nil,
		})
		return
	}

	var query string
	switch user.Role {
	case "siswa":
		query = "DELETE FROM siswa WHERE email = ?"
	case "admin":
		query = "DELETE FROM admin WHERE email = ?"
	case "guru":
		query = "DELETE FROM guru WHERE email = ?"
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid role parameter",
		})
		return
	}

	// Eksekusi raw query
	if err := db.Exec(query, user.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete " + user.Email,
		})
		return
	}

	// Jika sukses
	c.JSON(http.StatusOK, gin.H{
		"message": user.Role + " data successfully deleted",
	})
}

func GetPoint(c *gin.Context) {
	var response models.BaseResponseModel

	email, errmail := services.GetUserEmailFromToken(c)
	if errmail != nil {
		response = models.BaseResponseModel{
			Message: errmail.Error(),
			Data:    nil,
		}
		fmt.Println(errmail)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userPoint, message := repositories.GetUserPoint(email)

	fmt.Println("user diamond : " + strconv.Itoa(userPoint.Diamond))
	fmt.Println("user exp : " + strconv.Itoa(userPoint.Exp))

	if strings.Contains(strings.ToLower(message), "success") {
		response = models.BaseResponseModel{
			Message: message,
			Data:    userPoint,
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		response = models.BaseResponseModel{
			Message: message,
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
}

func UpdatePoint(c *gin.Context) {
	var response models.BaseResponseModel
	var input points.DiamondOrExp

	if err := c.ShouldBindJSON(&input); err != nil {
		response = models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	email, errmail := services.GetUserEmailFromToken(c)
	if errmail != nil {
		response = models.BaseResponseModel{
			Message: errmail.Error(),
			Data:    nil,
		}
		fmt.Println(errmail)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := config.DB
	// GORM akan hanya update field yang tidak nil
	tx := db.Model(&points.UserPoint{}).
		Where("email = ?", email).Updates(input)

	// fmt.Println(err.Error())

	if tx.Error != nil {
		response = models.BaseResponseModel{
			Message: tx.Error.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if tx.RowsAffected == 0 {
		response = models.BaseResponseModel{
			Message: "Data tidak ditemukan",
			Data:    nil,
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	response = models.BaseResponseModel{
		Message: "User point updated",
		Data:    input,
	}
	c.JSON(http.StatusOK, response)
}

func UpdateSiswaImageProfile(c *gin.Context) {
	email := c.Param("email")
	role := c.Param("role")
	db := config.DB
	validRoles := map[string]bool{"siswa": true, "admin": true, "guru": true}

	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Gunakan: siswa, admin, guru"})
		return
	}

	uploadResult, message := repositories.ReciveAndStoreImage(c)
	if uploadResult == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	switch role {
	case "siswa":
		if err := db.Model(&users.Siswa{}).Where("email = ?", email).Update("image_profile", uploadResult.Url).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
			return
		}
	case "guru":
		if err := db.Model(&users.Guru{}).Where("email = ?", email).Update("image_profile", uploadResult.Url).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
			return
		}
	case "admin":
		if err := db.Model(&users.Admin{}).Where("email = ?", email).Update("image_profile", uploadResult.Url).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Image profile %s berhasil diperbarui", role),
		"image_profile": uploadResult.Url,
	})
}
