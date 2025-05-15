package controllers

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// BatchGuruMapelRequest defines the request structure for batch operations on guru-mapel relationships
type BatchGuruMapelRequest struct {
	Teachers []users.GuruMapelRequest `json:"teachers" binding:"required"`
}

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

// UpdateGuruMapel updates the subject relationships for a teacher by replacing all existing relationships
func UpdateGuruMapel(c *gin.Context) {
	db := config.DB
	var req users.GuruMapelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the teacher exists
	var count int64
	if err := tx.Table("guru").Where("id = ?", req.IDGuru).Count(&count).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa data guru"})
		return
	}

	if count == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	// Delete all existing relationships for this teacher
	if err := tx.Where("id_guru = ?", req.IDGuru).Delete(&users.GuruMapel{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus relasi lama"})
		return
	}

	// Create new relationships
	var inserts []users.GuruMapel
	for _, idMapel := range req.IDMapels {
		inserts = append(inserts, users.GuruMapel{
			IDGuru:  req.IDGuru,
			IDMapel: idMapel,
		})
	}

	if len(inserts) > 0 {
		if err := tx.Create(&inserts).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat relasi baru"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel berhasil diperbarui",
		"data":    inserts,
	})
}

// DeleteGuruMapel deletes a specific teacher-subject relationship
func DeleteGuruMapel(c *gin.Context) {
	db := config.DB
	idGuru := c.Param("id_guru")
	idMapel := c.Param("id_mapel")

	// Convert string parameters to uint
	guruID, err := strconv.ParseUint(idGuru, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID guru tidak valid"})
		return
	}

	mapelID, err := strconv.ParseUint(idMapel, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mapel tidak valid"})
		return
	}

	// Delete the specific relationship
	result := db.Where("id_guru = ? AND id_mapel = ?", guruID, mapelID).Delete(&users.GuruMapel{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus relasi guru-mapel"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relasi guru-mapel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel berhasil dihapus",
	})
}

// DeleteAllGuruMapel deletes all subject relationships for a specific teacher
func DeleteAllGuruMapel(c *gin.Context) {
	db := config.DB
	idGuru := c.Param("id_guru")

	// Convert string parameter to uint
	guruID, err := strconv.ParseUint(idGuru, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID guru tidak valid"})
		return
	}

	// Check if the teacher exists
	var count int64
	if err := db.Table("guru").Where("id = ?", guruID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa data guru"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	// Delete all relationships for this teacher
	result := db.Where("id_guru = ?", guruID).Delete(&users.GuruMapel{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus relasi guru-mapel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Semua relasi guru-mapel berhasil dihapus",
		"count":   result.RowsAffected,
	})
}

// BatchAddGuruMapel adds subject relationships for multiple teachers at once
func BatchAddGuruMapel(c *gin.Context) {
	db := config.DB
	var batchReq BatchGuruMapelRequest
	if err := c.ShouldBindJSON(&batchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	responseData := make(map[uint][]users.GuruMapel)
	for _, teacherReq := range batchReq.Teachers {
		// Check if the teacher exists
		var count int64
		if err := tx.Table("guru").Where("id = ?", teacherReq.IDGuru).Count(&count).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal memeriksa data guru",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		if count == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "Guru tidak ditemukan",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		// Create relationships for this teacher
		var inserts []users.GuruMapel
		for _, idMapel := range teacherReq.IDMapels {
			inserts = append(inserts, users.GuruMapel{
				IDGuru:  teacherReq.IDGuru,
				IDMapel: idMapel,
			})
		}

		// Use INSERT IGNORE to skip duplicates
		result := tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "id_guru"}, {Name: "id_mapel"}},
				DoNothing: true,
			},
		).Create(&inserts)

		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal menambahkan relasi guru-mapel",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		responseData[teacherReq.IDGuru] = inserts
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel batch berhasil ditambahkan",
		"data":    responseData,
	})
}

// BatchUpdateGuruMapel updates subject relationships for multiple teachers at once
func BatchUpdateGuruMapel(c *gin.Context) {
	db := config.DB
	var batchReq BatchGuruMapelRequest
	if err := c.ShouldBindJSON(&batchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	responseData := make(map[uint][]users.GuruMapel)
	for _, teacherReq := range batchReq.Teachers {
		// Check if the teacher exists
		var count int64
		if err := tx.Table("guru").Where("id = ?", teacherReq.IDGuru).Count(&count).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal memeriksa data guru",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		if count == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "Guru tidak ditemukan",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		// Delete all existing relationships for this teacher
		if err := tx.Where("id_guru = ?", teacherReq.IDGuru).Delete(&users.GuruMapel{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal menghapus relasi lama",
				"id_guru": teacherReq.IDGuru,
			})
			return
		}

		// Create new relationships for this teacher
		var inserts []users.GuruMapel
		for _, idMapel := range teacherReq.IDMapels {
			inserts = append(inserts, users.GuruMapel{
				IDGuru:  teacherReq.IDGuru,
				IDMapel: idMapel,
			})
		}

		if len(inserts) > 0 {
			if err := tx.Create(&inserts).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Gagal membuat relasi baru",
					"id_guru": teacherReq.IDGuru,
				})
				return
			}
		}

		responseData[teacherReq.IDGuru] = inserts
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel batch berhasil diperbarui",
		"data":    responseData,
	})
}

// BatchDeleteGuruMapel deletes all subject relationships for multiple teachers
func BatchDeleteGuruMapel(c *gin.Context) {
	db := config.DB
	var request struct {
		TeacherIDs []uint `json:"teacher_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	results := make(map[uint]int64)
	for _, guruID := range request.TeacherIDs {
		// Check if the teacher exists
		var count int64
		if err := tx.Table("guru").Where("id = ?", guruID).Count(&count).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal memeriksa data guru",
				"id_guru": guruID,
			})
			return
		}

		if count == 0 {
			// Skip non-existent teachers but continue with others
			results[guruID] = 0
			continue
		}

		// Delete all relationships for this teacher
		result := tx.Where("id_guru = ?", guruID).Delete(&users.GuruMapel{})
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Gagal menghapus relasi guru-mapel",
				"id_guru": guruID,
			})
			return
		}

		results[guruID] = result.RowsAffected
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relasi guru-mapel batch berhasil dihapus",
		"results":  results,
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
