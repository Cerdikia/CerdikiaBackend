package controllers

import (
	"coba1BE/config"
	"coba1BE/models/genericActivities"
	"coba1BE/models/soal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	var module genericActivities.Module
	var soal soal.UploadSoal
	var moduleCount int64
	var soalCount int64
	db := config.DB

	// Hitung jumlah module
	if err := db.Model(&module).Count(&moduleCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghitung jumlah module",
			"error":   err.Error(),
		})
		return
	}

	// Hitung jumlah soal
	if err := db.Model(&soal).Count(&soalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghitung jumlah soal",
			"error":   err.Error(),
		})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Statistik berhasil diambil",
		"total_modules": moduleCount,
		"total_soal":    soalCount,
	})
}
