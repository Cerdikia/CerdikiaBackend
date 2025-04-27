package controllers

import (
	"coba1BE/config"
	"coba1BE/models/hadiah"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TukarBarang(c *gin.Context) {
	db := config.DB

	var input hadiah.PenukaranInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Cari barang
	var barang hadiah.TableBarang
	if err := db.First(&barang, "id_barang = ?", input.IdBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	// Cek stok cukup atau tidak
	if barang.Jumlah < input.Jumlah {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak cukup"})
		return
	}

	// Kurangi jumlah barang
	barang.Jumlah -= input.Jumlah
	if err := db.Save(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengurangi stok"})
		return
	}

	// Catat ke logs_penukaran_point
	log := hadiah.LogsPenukaranPoint{
		IdBarang: input.IdBarang,
		Email:    input.Email,
		Jumlah:   input.Jumlah,
	}
	if err := db.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat penukaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penukaran berhasil",
		"data":    log,
	})
}
