package controllers

import (
	"coba1BE/config"
	"coba1BE/models/kelas"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllKelas(c *gin.Context) {
	db := config.DB
	var kelas []kelas.Kelas
	if err := db.Find(&kelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kelas)
}

func GetKelasByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")
	var kelas kelas.Kelas
	if err := db.First(&kelas, "id_kelas = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas not found"})
		return
	}
	c.JSON(http.StatusOK, kelas)
}

func CreateKelas(c *gin.Context) {
	db := config.DB
	var kelas kelas.Kelas
	if err := c.ShouldBindJSON(&kelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&kelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, kelas)
}

func UpdateKelas(c *gin.Context) {
	db := config.DB
	id := c.Param("id")
	var dataKelas kelas.Kelas
	if err := db.First(&dataKelas, "id_kelas = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas not found"})
		return
	}
	if err := c.ShouldBindJSON(&dataKelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataKelas.IDKelas = 0 // Supaya tidak mengubah primary key
	if err := db.Model(&kelas.Kelas{}).Where("id_kelas = ?", id).Updates(dataKelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kelas updated successfully"})
}

func DeleteKelas(c *gin.Context) {
	db := config.DB
	id := c.Param("id")
	if err := db.Delete(&kelas.Kelas{}, "id_kelas = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kelas deleted successfully"})
}
