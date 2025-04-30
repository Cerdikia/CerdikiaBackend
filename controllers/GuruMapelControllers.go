package controllers

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddGuruMapel(c *gin.Context) {
	db := config.DB
	var req users.GuruMapelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var inserts []users.GuruMapel

	for _, idMapel := range req.IDMapels {
		inserts = append(inserts, users.GuruMapel{
			IDGuru:  req.IDGuru,
			IDMapel: idMapel,
		})
	}

	// Gunakan INSERT IGNORE mirip ON DUPLICATE KEY UPDATE
	result := db.Clauses(
		// GORM's way to skip duplicate key
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id_guru"}, {Name: "id_mapel"}},
			DoNothing: true,
		},
	).Create(&inserts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel berhasil ditambahkan",
		"data":    inserts,
	})
}

func GetMapelByGuru(c *gin.Context) {
	db := config.DB
	idGuru := c.Param("id_guru")

	// Ambil data guru
	var guru users.GuruMapelResponse
	err := db.Raw(`
		SELECT email, nama, jabatan 
		FROM guru 
		WHERE id = ?
	`, idGuru).Scan(&guru).Error

	if err != nil || guru.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	// Ambil daftar mapel yang dimiliki guru
	var mapelList []users.Mapel
	err = db.Raw(`
		SELECT m.id_mapel, m.mapel 
		FROM mapel m
		JOIN guru_mapel gm ON gm.id_mapel = m.id_mapel
		WHERE gm.id_guru = ?
	`, idGuru).Scan(&mapelList).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil mata pelajaran"})
		return
	}

	guru.Mapel = mapelList

	c.JSON(http.StatusOK, guru)
}
