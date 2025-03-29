package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"coba1BE/repositories"

	// "database/sql"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var user users.User

	// Bind JSON dari request ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi role: jika kosong, set ke "siswa"
	validRoles := map[string]bool{"siswa": true, "admin": true, "guru": true}

	// Jika role kosong, set default "siswa"
	if user.Role == "" {
		user.Role = "siswa"
	} else {
		// Pastikan role yang dikirim valid
		user.Role = strings.ToLower(user.Role) // Buat case insensitive
		if !validRoles[user.Role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Gunakan: siswa, admin, guru"})
			return
		}
	}

	// Simpan ke database
	// if err := database.DB.Create(&user).Error; err != nil {
	db := config.DB
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat", "user": user})
}

func GetUsers(c *gin.Context) {
	var response models.BaseResponseModel

	response = repositories.GetAllUsers()

	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context) {
	var response models.BaseResponseModel

	// Tangkap parameter email dari query
	email := c.Param("email")

	// fmt.Println("email controller : " + email)
	response = repositories.GetUserByEmail(email)
	c.JSON(http.StatusOK, response)
}
