package controllers

import (
	"coba1BE/config"
	"coba1BE/models/logsprogres"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RekapSemester(c *gin.Context) {
	db := config.DB
	tahunAjaran := c.Query("tahun_ajaran") // Contoh: "2024/2025"
	if tahunAjaran == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tahun_ajaran diperlukan"})
		return
	}

	var logs []logsprogres.LogEntry
	if err := db.Table("logs").Select("email, id_kelas, id_mapel, id_module, skor, created_at").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data logs"})
		return
	}

	// Map untuk mengelompokkan progres per siswa
	siswaMap := make(map[string]struct {
		IDKelas int
		Progres []logsprogres.ProgresItem
	})

	for _, log := range logs {
		entry := siswaMap[log.Email]
		entry.IDKelas = log.IDKelas
		entry.Progres = append(entry.Progres, logsprogres.ProgresItem{
			IDMapel:   log.IDMapel,
			IDModule:  log.IDModule,
			Skor:      log.Skor,
			CreatedAt: log.CreatedAt,
		})
		siswaMap[log.Email] = entry
	}

	// Simpan ke data_siswa
	for email, data := range siswaMap {
		progresJSON, err := json.Marshal(data.Progres)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-encode progres JSON"})
			return
		}

		if err := db.Table("data_siswa").Create(map[string]interface{}{
			"email":        email,
			"id_kelas":     data.IDKelas,
			"progres":      string(progresJSON),
			"tahun_ajaran": tahunAjaran,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal menyimpan data untuk %s", email)})
			return
		}
	}

	// Optional: hapus data logs setelah dipindah
	// db.Exec("DELETE FROM logs")
	// Hapus data logs setelah rekap berhasil
	if err := db.Exec("DELETE FROM logs").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data siswa berhasil direkap untuk tahun ajaran " + tahunAjaran})
}

func EditTahunAjaran(c *gin.Context) {
	db := config.DB
	var input struct {
		TahunAjaranLama string `json:"tahun_ajaran_lama"`
		TahunAjaranBaru string `json:"tahun_ajaran_baru"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if input.TahunAjaranLama == "" || input.TahunAjaranBaru == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun ajaran lama dan baru harus diisi"})
		return
	}

	result := db.Model(&logsprogres.DataSiswa{}).Where("tahun_ajaran = ?", input.TahunAjaranLama).Update("tahun_ajaran", input.TahunAjaranBaru)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Tahun ajaran berhasil diperbarui",
		"rows_affected":     result.RowsAffected,
		"tahun_ajaran_baru": input.TahunAjaranBaru,
	})
}

func GetDataSiswa(c *gin.Context) {
	db := config.DB
	idData := c.Param("id_data") // ID data siswa yang ingin dibaca

	if idData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_data diperlukan"})
		return
	}

	var dataSiswa struct {
		Email       string `json:"email"`
		IDKelas     int    `json:"id_kelas"`
		Progres     string `json:"progres"`
		TahunAjaran string `json:"tahun_ajaran"`
	}

	if err := db.Table("data_siswa").Where("id_data = ?", idData).First(&dataSiswa).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data siswa tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":        dataSiswa.Email,
		"id_kelas":     dataSiswa.IDKelas,
		"progres":      dataSiswa.Progres,
		"tahun_ajaran": dataSiswa.TahunAjaran,
	})
}

func GetAllDataSiswa(c *gin.Context) {
	db := config.DB

	var dataSiswa []struct {
		IDData      int    `json:"id_data"`
		Email       string `json:"email"`
		IDKelas     int    `json:"id_kelas"`
		Progres     string `json:"progres"`
		TahunAjaran string `json:"tahun_ajaran"`
	}

	if err := db.Table("data_siswa").Find(&dataSiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	c.JSON(http.StatusOK, dataSiswa)
}

func DeleteDataSiswa(c *gin.Context) {
	db := config.DB
	idData := c.Param("id_data") // ID data siswa yang ingin dihapus

	if idData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_data diperlukan"})
		return
	}

	// Hapus data siswa berdasarkan id_data
	if err := db.Table("data_siswa").Where("id_data = ?", idData).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data siswa berhasil dihapus",
	})
}
