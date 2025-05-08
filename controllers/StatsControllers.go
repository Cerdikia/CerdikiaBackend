package controllers

import (
	"coba1BE/config"
	"coba1BE/models/genericActivities"
	"coba1BE/models/soal"
	"coba1BE/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetStats adalah endpoint lama yang hanya menampilkan jumlah module dan soal
// Dipertahankan untuk kompatibilitas dengan aplikasi yang sudah ada
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

// GetAllStats adalah endpoint baru yang menampilkan semua statistik aplikasi
// Menggabungkan semua statistik yang sebelumnya diambil dari beberapa endpoint berbeda
func GetAllStats(c *gin.Context) {
	// Ambil semua statistik dari repository
	statsResponse := repositories.GetAllStats()

	// Kirim response
	c.JSON(http.StatusOK, gin.H{
		"message": "Statistik aplikasi berhasil diambil",
		"data": statsResponse,
	})
}

// GetRecentActivities adalah endpoint untuk menampilkan aktivitas terakhir dari tabel logs
// dengan pagination (default: page 1, limit 7)
func GetRecentActivities(c *gin.Context) {
	// Ambil query params untuk pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "7")

	// Konversi ke int
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 7
	}

	// Ambil aktivitas terakhir dari repository
	activitiesResponse := repositories.GetRecentActivities(pageInt, limitInt)

	// Kirim response
	c.JSON(http.StatusOK, gin.H{
		"message": "Aktivitas terakhir berhasil diambil",
		"data": activitiesResponse,
	})
}
