package controllers

import (
	"coba1BE/config"
	"coba1BE/models/logsprogres"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RekapSemester mengumpulkan data dari tabel logs dan menyimpannya ke tabel data_siswa
// dengan pengelompokan berdasarkan siswa dan tahun ajaran
func RekapSemester(c *gin.Context) {
	db := config.DB

	// Validasi parameter input
	var input struct {
		TahunAjaran    string `json:"tahun_ajaran" form:"tahun_ajaran" binding:"required"` // Contoh: "2024/2025"
		Semester       string `json:"semester" form:"semester" binding:"required"`         // Contoh: "Ganjil" atau "Genap"
		DeleteLogsData bool   `json:"delete_logs_data" form:"delete_logs_data"`            // Apakah data logs akan dihapus setelah rekap
		FilterKelas    *int   `json:"filter_kelas" form:"filter_kelas"`                    // Opsional: filter berdasarkan kelas tertentu
		FilterMapel    *int   `json:"filter_mapel" form:"filter_mapel"`                    // Opsional: filter berdasarkan mapel tertentu
		StartDate      string `json:"start_date" form:"start_date"`                        // Opsional: filter berdasarkan tanggal mulai (format: YYYY-MM-DD)
		EndDate        string `json:"end_date" form:"end_date"`                            // Opsional: filter berdasarkan tanggal akhir (format: YYYY-MM-DD)
	}

	// Coba binding dari query parameters terlebih dahulu
	if err := c.ShouldBindQuery(&input); err != nil {
		// Jika gagal, coba binding dari JSON body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter tidak valid",
				"error":   err.Error(),
			})
			return
		}
	}

	// Validasi format tahun ajaran
	if !isValidTahunAjaran(input.TahunAjaran) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format tahun ajaran tidak valid",
			"format":  "YYYY/YYYY, contoh: 2024/2025",
		})
		return
	}

	// Validasi semester
	if input.Semester != "Ganjil" && input.Semester != "Genap" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Nilai semester tidak valid",
			"valid":   []string{"Ganjil", "Genap"},
		})
		return
	}

	// Memulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Membangun query untuk mengambil data logs
	logsQuery := tx.Table("logs").Select("email, id_kelas, id_mapel, id_module, skor, created_at")

	// Menambahkan filter jika ada
	if input.FilterKelas != nil {
		logsQuery = logsQuery.Where("id_kelas = ?", *input.FilterKelas)
	}

	if input.FilterMapel != nil {
		logsQuery = logsQuery.Where("id_mapel = ?", *input.FilterMapel)
	}

	// Filter berdasarkan tanggal jika disediakan
	if input.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal mulai tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		logsQuery = logsQuery.Where("created_at >= ?", startDate)
	}

	if input.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal akhir tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		// Tambahkan 1 hari ke endDate untuk mencakup seluruh hari
		endDate = endDate.Add(24 * time.Hour)
		logsQuery = logsQuery.Where("created_at < ?", endDate)
	}

	// Mengambil data logs sesuai filter
	var logs []logsprogres.LogEntry
	if err := logsQuery.Find(&logs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data logs",
			"error":   err.Error(),
		})
		return
	}

	// Jika tidak ada data yang ditemukan
	if len(logs) == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tidak ada data logs yang ditemukan sesuai filter",
		})
		return
	}

	// Map untuk mengelompokkan progres per siswa
	siswaMap := make(map[string]struct {
		IDKelas int
		Progres []logsprogres.ProgresItem
	})

	// Mengelompokkan data logs berdasarkan email siswa
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

	// Menyimpan data ke tabel data_siswa
	successCount := 0
	skippedCount := 0
	failedEmails := []string{}

	for email, data := range siswaMap {
		// Mengecek apakah data untuk siswa dengan tahun ajaran dan semester yang sama sudah ada
		var existingCount int64
		if err := tx.Table("data_siswa").Where("email = ? AND tahun_ajaran = ? AND semester = ?", email, input.TahunAjaran, input.Semester).Count(&existingCount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal memeriksa data yang sudah ada",
				"error":   err.Error(),
			})
			return
		}

		// Jika data sudah ada, lewati
		if existingCount > 0 {
			skippedCount++
			continue
		}

		// Mengkonversi data progres ke JSON
		progresJSON, err := json.Marshal(data.Progres)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal meng-encode progres JSON",
				"error":   err.Error(),
			})
			return
		}

		// Menyimpan data ke tabel data_siswa
		if err := tx.Table("data_siswa").Create(map[string]interface{}{
			"email":        email,
			"id_kelas":     data.IDKelas,
			"progres":      string(progresJSON),
			"tahun_ajaran": input.TahunAjaran,
			"semester":     input.Semester,
			"created_at":   time.Now(),
		}).Error; err != nil {
			failedEmails = append(failedEmails, email)
		} else {
			successCount++
		}
	}

	// Jika ada email yang gagal disimpan, rollback transaksi
	if len(failedEmails) > 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":       "Gagal menyimpan data untuk beberapa siswa",
			"failed_emails": failedEmails,
		})
		return
	}

	// Jika diminta, hapus data logs yang sudah direkap
	if input.DeleteLogsData {
		// Membuat query delete yang sesuai dengan filter yang digunakan untuk mengambil data
		deleteQuery := tx.Table("logs")

		if input.FilterKelas != nil {
			deleteQuery = deleteQuery.Where("id_kelas = ?", *input.FilterKelas)
		}

		if input.FilterMapel != nil {
			deleteQuery = deleteQuery.Where("id_mapel = ?", *input.FilterMapel)
		}

		if input.StartDate != "" {
			startDate, _ := time.Parse("2006-01-02", input.StartDate)
			deleteQuery = deleteQuery.Where("created_at >= ?", startDate)
		}

		if input.EndDate != "" {
			endDate, _ := time.Parse("2006-01-02", input.EndDate)
			endDate = endDate.Add(24 * time.Hour)
			deleteQuery = deleteQuery.Where("created_at < ?", endDate)
		}

		// Menghapus data logs
		result := deleteQuery.Delete(nil)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal menghapus data logs",
				"error":   result.Error.Error(),
			})
			return
		}
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan perubahan",
			"error":   err.Error(),
		})
		return
	}

	// Mengembalikan respons sukses
	c.JSON(http.StatusOK, gin.H{
		"message":       "Data siswa berhasil direkap",
		"tahun_ajaran":  input.TahunAjaran,
		"semester":      input.Semester,
		"success_count": successCount,
		"skipped_count": skippedCount,
		"logs_deleted":  input.DeleteLogsData,
	})
}

// isValidTahunAjaran memvalidasi format tahun ajaran (YYYY/YYYY)
func isValidTahunAjaran(tahunAjaran string) bool {
	// Format yang diharapkan: YYYY/YYYY, contoh: 2024/2025
	var tahunAwal, tahunAkhir int
	if _, err := fmt.Sscanf(tahunAjaran, "%d/%d", &tahunAwal, &tahunAkhir); err != nil {
		return false
	}

	// Validasi tahun (antara 2000 dan 2100)
	if tahunAwal < 2000 || tahunAwal > 2100 || tahunAkhir < 2000 || tahunAkhir > 2100 {
		return false
	}

	// Tahun akhir harus lebih besar dari tahun awal
	if tahunAkhir != tahunAwal+1 {
		return false
	}

	return true
}

// EditTahunAjaran memperbarui tahun ajaran pada data yang sudah ada
func EditTahunAjaran(c *gin.Context) {
	db := config.DB
	var input struct {
		TahunAjaranLama string `json:"tahun_ajaran_lama" binding:"required"`
		TahunAjaranBaru string `json:"tahun_ajaran_baru" binding:"required"`
		Semester        string `json:"semester" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validasi format tahun ajaran
	if !isValidTahunAjaran(input.TahunAjaranLama) || !isValidTahunAjaran(input.TahunAjaranBaru) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format tahun ajaran tidak valid",
			"format":  "YYYY/YYYY, contoh: 2024/2025",
		})
		return
	}

	// Validasi semester
	if input.Semester != "Ganjil" && input.Semester != "Genap" && input.Semester != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Nilai semester tidak valid",
			"valid":   []string{"Ganjil", "Genap", ""},
		})
		return
	}

	// Memulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Membuat query dasar
	query := tx.Model(&logsprogres.DataSiswa{}).Where("tahun_ajaran = ?", input.TahunAjaranLama)

	// Jika semester disediakan, tambahkan ke filter
	if input.Semester != "" {
		query = query.Where("semester = ?", input.Semester)
	}

	// Melakukan update
	result := query.Update("tahun_ajaran", input.TahunAjaranBaru)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memperbarui data",
			"error":   result.Error.Error(),
		})
		return
	}

	// Jika tidak ada data yang diperbarui
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tidak ada data yang ditemukan dengan tahun ajaran dan semester tersebut",
		})
		return
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan perubahan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Tahun ajaran berhasil diperbarui",
		"rows_affected":     result.RowsAffected,
		"tahun_ajaran_lama": input.TahunAjaranLama,
		"tahun_ajaran_baru": input.TahunAjaranBaru,
		"semester":          input.Semester,
	})
}

// GetDataSiswa mendapatkan data siswa berdasarkan ID data
func GetDataSiswa(c *gin.Context) {
	db := config.DB
	idData := c.Param("id_data") // ID data siswa yang ingin dibaca

	if idData == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID data diperlukan",
		})
		return
	}

	// Struktur untuk menyimpan data siswa
	var dataSiswa struct {
		IDData      uint      `json:"id_data" gorm:"column:id_data"`
		Email       string    `json:"email" gorm:"column:email"`
		IDKelas     int       `json:"id_kelas" gorm:"column:id_kelas"`
		Progres     string    `json:"progres" gorm:"column:progres"`
		TahunAjaran string    `json:"tahun_ajaran" gorm:"column:tahun_ajaran"`
		Semester    string    `json:"semester" gorm:"column:semester"`
		CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	}

	// Menggunakan gorm.ErrRecordNotFound untuk menangani kasus data tidak ditemukan
	if err := db.Table("data_siswa").Where("id_data = ?", idData).First(&dataSiswa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data siswa tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mengambil data siswa",
				"error":   err.Error(),
			})
		}
		return
	}

	// Mengurai data progres dari JSON string ke array objek
	var progresItems []logsprogres.ProgresItem
	if err := json.Unmarshal([]byte(dataSiswa.Progres), &progresItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengurai data progres",
			"error":   err.Error(),
		})
		return
	}

	// Mengembalikan data dalam format yang lebih terstruktur
	c.JSON(http.StatusOK, gin.H{
		"message": "Data siswa berhasil diambil",
		"data": gin.H{
			"id_data":      dataSiswa.IDData,
			"email":        dataSiswa.Email,
			"id_kelas":     dataSiswa.IDKelas,
			"progres":      progresItems,
			"tahun_ajaran": dataSiswa.TahunAjaran,
			"semester":     dataSiswa.Semester,
			"created_at":   dataSiswa.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// GetAllDataSiswa mendapatkan semua data rekap semester dengan filter opsional
func GetAllDataSiswa(c *gin.Context) {
	db := config.DB

	// Parameter filter opsional
	tahunAjaran := c.Query("tahun_ajaran")
	idKelas := c.Query("id_kelas")
	email := c.Query("email")

	// Debug: Log the request parameters
	fmt.Println("GetAllDataSiswa request params - tahunAjaran:", tahunAjaran, "idKelas:", idKelas, "email:", email)

	// Struktur untuk menyimpan data siswa
	type DataSiswaResponse struct {
		IDData      uint      `json:"id_data" gorm:"column:id_data"`
		Email       string    `json:"email" gorm:"column:email"`
		IDKelas     int       `json:"id_kelas" gorm:"column:id_kelas"`
		Kelas       string    `json:"kelas" gorm:"column:kelas"`
		NamaSiswa   string    `json:"nama_siswa" gorm:"column:nama_siswa"`
		Progres     string    `json:"progres" gorm:"column:progres"`
		TahunAjaran string    `json:"tahun_ajaran" gorm:"column:tahun_ajaran"`
		CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
		Semester    string    `json:"semester" form:"semester"`
		// CreatedAt field removed as it doesn't exist in the data_siswa table
	}

	// Membuat query dasar dengan join ke tabel siswa dan kelas untuk mendapatkan informasi tambahan
	query := db.Table("data_siswa ds").
		Select("ds.id_data, ds.email, ds.id_kelas, k.kelas, s.nama as nama_siswa, ds.progres, ds.tahun_ajaran, ds.created_at, ds.semester").
		Joins("LEFT JOIN siswa s ON ds.email = s.email").
		Joins("LEFT JOIN kelas k ON ds.id_kelas = k.id_kelas")

	// Debug: Count total records before applying filters
	var totalCount int64
	if err := db.Table("data_siswa").Count(&totalCount).Error; err != nil {
		fmt.Println("Error counting total records:", err)
	} else {
		fmt.Println("Total records in data_siswa table:", totalCount)
	}

	// Menambahkan filter jika ada
	if tahunAjaran != "" {
		query = query.Where("ds.tahun_ajaran = ?", tahunAjaran)
	}

	// Semester filter is removed as the column doesn't exist in the database

	if idKelas != "" {
		query = query.Where("ds.id_kelas = ?", idKelas)
	}

	if email != "" {
		query = query.Where("ds.email = ?", email)
	}

	// Mengurutkan data berdasarkan ID (terbaru lebih dulu)
	query = query.Order("ds.id_data DESC")

	// Mengambil data
	var dataSiswa []DataSiswaResponse
	if err := query.Find(&dataSiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data siswa",
			"error":   err.Error(),
		})
		return
	}

	// Debug: Log the number of records found
	fmt.Println("Records found after applying filters:", len(dataSiswa))

	// Jika tidak ada data yang ditemukan, return empty array but don't exit
	if len(dataSiswa) == 0 {
		// Instead of returning immediately, we'll continue with an empty array
		// This ensures the API response format is consistent
		dataSiswa = []DataSiswaResponse{}
	}

	// Memproses data untuk respons
	responseData := make([]map[string]interface{}, 0, len(dataSiswa))
	for _, data := range dataSiswa {
		// Mengurai data progres dari JSON string ke array objek
		var progresItems []logsprogres.ProgresItem
		if err := json.Unmarshal([]byte(data.Progres), &progresItems); err != nil {
			// Jika gagal mengurai, gunakan string progres asli
			responseData = append(responseData, map[string]interface{}{
				"id_data":      data.IDData,
				"email":        data.Email,
				"id_kelas":     data.IDKelas,
				"kelas":        data.Kelas,
				"nama_siswa":   data.NamaSiswa,
				"progres":      data.Progres, // String JSON asli
				"tahun_ajaran": data.TahunAjaran,
				"created_at":   data.CreatedAt,
				"semester":     data.Semester,
			})
		} else {
			// Jika berhasil mengurai, gunakan array objek progres
			responseData = append(responseData, map[string]interface{}{
				"id_data":      data.IDData,
				"email":        data.Email,
				"id_kelas":     data.IDKelas,
				"kelas":        data.Kelas,
				"nama_siswa":   data.NamaSiswa,
				"progres":      progresItems,
				"tahun_ajaran": data.TahunAjaran,
				"created_at":   data.CreatedAt,
				"semester":     data.Semester,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data rekap semester berhasil diambil",
		"count":   len(responseData),
		"data":    responseData,
	})
}

// RekapSemesterSiswa mengumpulkan data dari tabel logs untuk satu siswa dan menyimpannya ke tabel data_siswa
// kemudian menghapus data logs untuk siswa tersebut
func RekapSemesterSiswa(c *gin.Context) {
	db := config.DB

	// Validasi parameter input
	var input struct {
		Email       string `json:"email" form:"email" binding:"required"`
		TahunAjaran string `json:"tahun_ajaran" form:"tahun_ajaran" binding:"required"` // Contoh: "2024/2025"
		DeleteLogs  bool   `json:"delete_logs" form:"delete_logs" binding:""`           // Apakah data logs akan dihapus setelah rekap
	}

	// Coba binding dari query parameters terlebih dahulu
	if err := c.ShouldBindQuery(&input); err != nil {
		// Jika gagal, coba binding dari JSON body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter tidak valid",
				"error":   err.Error(),
			})
			return
		}
	}

	// Validasi format tahun ajaran
	if !isValidTahunAjaran(input.TahunAjaran) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format tahun ajaran tidak valid",
			"format":  "YYYY/YYYY, contoh: 2024/2025",
		})
		return
	}

	// Memulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Memeriksa apakah siswa ada
	var siswaCount int64
	if err := tx.Table("siswa").Where("email = ?", input.Email).Count(&siswaCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memeriksa data siswa",
			"error":   err.Error(),
		})
		return
	}

	if siswaCount == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Siswa dengan email tersebut tidak ditemukan",
		})
		return
	}

	// Mendapatkan data id_kelas siswa
	var siswa struct {
		Email   string `gorm:"column:email"`
		IDKelas int    `gorm:"column:id_kelas"`
		Nama    string `gorm:"column:nama"`
	}

	if err := tx.Table("siswa").Select("email, id_kelas, nama").Where("email = ?", input.Email).First(&siswa).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mendapatkan data siswa",
			"error":   err.Error(),
		})
		return
	}

	// Mengambil data logs untuk siswa tersebut
	var logs []logsprogres.Log
	if err := tx.Table("logs").Where("email = ?", input.Email).Find(&logs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data logs",
			"error":   err.Error(),
		})
		return
	}

	// Jika tidak ada data logs untuk siswa tersebut
	if len(logs) == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tidak ada data logs untuk siswa tersebut",
		})
		return
	}

	// Mengkonversi data logs ke format ProgresItem
	progresItems := make([]logsprogres.ProgresItem, 0, len(logs))
	for _, log := range logs {
		progresItems = append(progresItems, logsprogres.ProgresItem{
			IDMapel:   log.IDMapel,
			IDModule:  log.IDModule,
			Skor:      log.Skor,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Mengkonversi ke JSON
	progresJSON, err := json.Marshal(progresItems)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengkonversi data ke JSON",
			"error":   err.Error(),
		})
		return
	}

	// Memeriksa apakah sudah ada data rekap untuk siswa dan tahun ajaran yang sama
	var existingData logsprogres.DataSiswa
	var isUpdate bool
	result := tx.Table("data_siswa").Where("email = ? AND tahun_ajaran = ?", input.Email, input.TahunAjaran).First(&existingData)
	if result.Error == nil {
		// Data sudah ada, update saja
		isUpdate = true
	} else if result.Error != gorm.ErrRecordNotFound {
		// Error selain record not found
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memeriksa data rekap yang sudah ada",
			"error":   result.Error.Error(),
		})
		return
	}

	// Menyimpan data ke tabel data_siswa
	var idData uint
	if isUpdate {
		// Update data yang sudah ada
		if err := tx.Table("data_siswa").Where("id_data = ?", existingData.IDData).Updates(map[string]interface{}{
			"progres": string(progresJSON),
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal memperbarui data rekap",
				"error":   err.Error(),
			})
			return
		}
		idData = existingData.IDData
	} else {
		// Buat data baru
		newData := map[string]interface{}{
			"email":        input.Email,
			"id_kelas":     siswa.IDKelas,
			"progres":      string(progresJSON),
			"tahun_ajaran": input.TahunAjaran,
		}

		result := tx.Table("data_siswa").Create(newData)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal menyimpan data rekap",
				"error":   result.Error.Error(),
			})
			return
		}

		// Mendapatkan ID data yang baru dibuat
		var lastInsertID uint
		if err := tx.Table("data_siswa").Select("id_data").Where("email = ? AND tahun_ajaran = ?", input.Email, input.TahunAjaran).Order("id_data DESC").Limit(1).Scan(&lastInsertID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mendapatkan ID data yang baru dibuat",
				"error":   err.Error(),
			})
			return
		}
		idData = lastInsertID
	}

	// Jika diminta untuk menghapus data logs
	if input.DeleteLogs {
		if err := tx.Table("logs").Where("email = ?", input.Email).Delete(nil).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal menghapus data logs",
				"error":   err.Error(),
			})
			return
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan perubahan",
			"error":   err.Error(),
		})
		return
	}

	// Mengembalikan respons sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Rekap semester untuk siswa berhasil dibuat",
		"data": map[string]interface{}{
			"id_data":      idData,
			"email":        input.Email,
			"nama_siswa":   siswa.Nama,
			"id_kelas":     siswa.IDKelas,
			"tahun_ajaran": input.TahunAjaran,
			"jumlah_logs":  len(logs),
			"logs_deleted": input.DeleteLogs,
		},
	})
}

// RekapSemesterAllSiswa mengumpulkan data dari tabel logs untuk semua siswa dan menyimpannya ke tabel data_siswa
// kemudian menghapus data logs jika diminta
func RekapSemesterAllSiswa(c *gin.Context) {
	db := config.DB

	// Validasi parameter input
	var input struct {
		TahunAjaran    string `json:"tahun_ajaran" form:"tahun_ajaran" binding:"required"` // Contoh: "2024/2025"
		DeleteLogsData bool   `json:"delete_logs_data" form:"delete_logs_data" binding:""` // Apakah data logs akan dihapus setelah rekap
		FilterKelas    *int   `json:"filter_kelas" form:"filter_kelas"`                    // Opsional: filter berdasarkan kelas tertentu
		StartDate      string `json:"start_date" form:"start_date"`                        // Opsional: filter berdasarkan tanggal mulai (format: YYYY-MM-DD)
		EndDate        string `json:"end_date" form:"end_date"`                            // Opsional: filter berdasarkan tanggal akhir (format: YYYY-MM-DD)
		// CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
		// CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
		Semester string `json:"semester" form:"semester"`
	}

	// Coba binding dari query parameters terlebih dahulu
	if err := c.ShouldBindQuery(&input); err != nil {
		// Jika gagal, coba binding dari JSON body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter tidak valid",
				"error":   err.Error(),
			})
			return
		}
	}

	// Validasi format tahun ajaran
	if !isValidTahunAjaran(input.TahunAjaran) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format tahun ajaran tidak valid",
			"format":  "YYYY/YYYY, contoh: 2024/2025",
		})
		return
	}

	// Memulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mendapatkan daftar semua siswa (atau filter berdasarkan kelas jika ada)
	var siswaQuery *gorm.DB = tx.Table("siswa").Select("email, id_kelas, nama")
	if input.FilterKelas != nil {
		fmt.Println("masuk kondisi filtered kelas")
		siswaQuery = siswaQuery.Where("id_kelas = ?", *input.FilterKelas)
	}

	var siswaList []struct {
		Email   string `gorm:"column:email"`
		IDKelas int    `gorm:"column:id_kelas"`
		Nama    string `gorm:"column:nama"`
	}

	if err := siswaQuery.Find(&siswaList).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mendapatkan daftar siswa",
			"error":   err.Error(),
		})
		return
	}

	if len(siswaList) == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Tidak ada siswa yang ditemukan",
		})
		return
	}

	// Membangun query untuk filter logs berdasarkan tanggal jika disediakan
	logsQuery := tx.Table("logs")

	// Debug: Print jumlah total logs di database
	var totalLogsCount int64
	if err := tx.Table("logs").Count(&totalLogsCount).Error; err != nil {
		fmt.Println("Error counting logs:", err)
	} else {
		fmt.Println("Total logs in database:", totalLogsCount)
	}

	// Filter berdasarkan tanggal jika disediakan
	if input.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal mulai tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		var totalLogsCountfiler1 int64
		logsQuery = logsQuery.Where("created_at >= ?", startDate)
		logsQuery.Count(&totalLogsCountfiler1)
		fmt.Println("Total logs after start date filter:", totalLogsCountfiler1)
	}

	if input.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal akhir tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		// Tambahkan 1 hari ke endDate untuk mencakup seluruh hari
		endDate = endDate.Add(24 * time.Hour)
		logsQuery = logsQuery.Where("created_at < ?", endDate)
		var totalLogsCountfiler2 int64
		logsQuery.Count(&totalLogsCountfiler2)
		fmt.Println("Total logs after end date filter:", totalLogsCountfiler2)
	}

	// Variabel untuk menyimpan hasil rekap
	var hasilRekap []map[string]interface{}
	// Variabel untuk menghitung jumlah total logs yang diproses
	totalLogsProcessed := 0

	// Memproses setiap siswa
	for _, siswa := range siswaList {
		fmt.Println("Processing student:", siswa.Email)

		// Create a new query for each student to avoid issues with the shared logsQuery
		studentLogsQuery := tx.Table("logs")

		// Apply the same date filters as the main query
		if input.StartDate != "" {
			startDate, _ := time.Parse("2006-01-02", input.StartDate)
			studentLogsQuery = studentLogsQuery.Where("created_at >= ?", startDate)
		}

		if input.EndDate != "" {
			endDate, _ := time.Parse("2006-01-02", input.EndDate)
			endDate = endDate.Add(24 * time.Hour)
			studentLogsQuery = studentLogsQuery.Where("created_at < ?", endDate)
		}

		// Add the student email filter
		studentLogsQuery = studentLogsQuery.Where("email = ?", siswa.Email)

		// Debug: Count logs for this student before retrieving them
		var studentLogsCount int64
		studentLogsQuery.Count(&studentLogsCount)
		fmt.Println("Logs count for student", siswa.Email, ":", studentLogsCount)

		// Retrieve the logs
		var logs []logsprogres.Log
		if logsErr := studentLogsQuery.Find(&logs).Error; logsErr != nil {
			fmt.Println("Error retrieving logs for student:", siswa.Email, "error:", logsErr)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mengambil data logs untuk siswa: " + siswa.Email,
				"error":   logsErr.Error(),
			})
			return
		}

		// Process student data whether they have logs or not
		if len(logs) == 0 {
			fmt.Println("No logs found for student:", siswa.Email)
			// Create empty progress data for students with no logs
			emptyProgresJSON, _ := json.Marshal([]logsprogres.ProgresItem{})

			// Check if there's already an entry for this student and academic year
			var existingData logsprogres.DataSiswa
			result := tx.Table("data_siswa").Where("email = ? AND tahun_ajaran = ?", siswa.Email, input.TahunAjaran).First(&existingData)

			var idData uint
			if result.Error == nil {
				// Update existing entry
				idData = existingData.IDData
				if err := tx.Table("data_siswa").Where("id_data = ?", existingData.IDData).Updates(map[string]interface{}{
					"progres": string(emptyProgresJSON),
				}).Error; err != nil {
					fmt.Println("Error updating data for student:", siswa.Email, "error:", err)
					continue
				}
			} else if result.Error == gorm.ErrRecordNotFound {
				// Create new entry
				newData := map[string]interface{}{
					"email":        siswa.Email,
					"id_kelas":     siswa.IDKelas,
					"progres":      string(emptyProgresJSON),
					"tahun_ajaran": input.TahunAjaran,
				}

				result := tx.Table("data_siswa").Create(newData)
				if result.Error != nil {
					fmt.Println("Error creating data for student:", siswa.Email, "error:", result.Error)
					continue
				}

				// Get the ID of the newly created entry
				var lastInsertID uint
				if err := tx.Table("data_siswa").Select("id_data").Where("email = ? AND tahun_ajaran = ?", siswa.Email, input.TahunAjaran).Order("id_data DESC").Limit(1).Scan(&lastInsertID).Error; err != nil {
					fmt.Println("Error getting ID for student:", siswa.Email, "error:", err)
					continue
				}
				idData = lastInsertID
			} else {
				// Other error
				fmt.Println("Error checking existing data for student:", siswa.Email, "error:", result.Error)
				continue
			}

			// Add to results
			hasilRekap = append(hasilRekap, map[string]interface{}{
				"id_data":      idData,
				"email":        siswa.Email,
				"nama_siswa":   siswa.Nama,
				"id_kelas":     siswa.IDKelas,
				"tahun_ajaran": input.TahunAjaran,
				"jumlah_logs":  0,
			})
			continue
		}

		// If we get here, the student has logs
		fmt.Println("Found", len(logs), "logs for student:", siswa.Email)

		totalLogsProcessed += len(logs)

		fmt.Println("ada logs yang terbaca")

		// Mengkonversi data logs ke format ProgresItem
		progresItems := make([]logsprogres.ProgresItem, 0, len(logs))
		for _, log := range logs {
			progresItems = append(progresItems, logsprogres.ProgresItem{
				IDMapel:   log.IDMapel,
				IDModule:  log.IDModule,
				Skor:      log.Skor,
				CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		// Mengkonversi ke JSON
		progresJSON, err := json.Marshal(progresItems)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mengkonversi data ke JSON untuk siswa: " + siswa.Email,
				"error":   err.Error(),
			})
			return
		}

		// Memeriksa apakah sudah ada data rekap untuk siswa dan tahun ajaran yang sama
		var existingData logsprogres.DataSiswa
		var idData uint

		result := tx.Table("data_siswa").Where("email = ? AND tahun_ajaran = ?", siswa.Email, input.TahunAjaran).First(&existingData)
		if result.Error == nil {
			// Data sudah ada, update saja
			idData = existingData.IDData

			if err := tx.Table("data_siswa").Where("id_data = ?", existingData.IDData).Updates(map[string]interface{}{
				"progres": string(progresJSON),
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal memperbarui data rekap untuk siswa: " + siswa.Email,
					"error":   err.Error(),
				})
				return
			}
		} else if result.Error != gorm.ErrRecordNotFound {
			// Error selain record not found
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal memeriksa data rekap yang sudah ada untuk siswa: " + siswa.Email,
				"error":   result.Error.Error(),
			})
			return
		} else {
			// Buat data baru
			newData := map[string]interface{}{
				"email":        siswa.Email,
				"id_kelas":     siswa.IDKelas,
				"progres":      string(progresJSON),
				"tahun_ajaran": input.TahunAjaran,
				"created_at":   input.CreatedAt,
				"semester":     input.Semester,
			}

			result := tx.Table("data_siswa").Create(newData)
			if result.Error != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal menyimpan data rekap untuk siswa: " + siswa.Email,
					"error":   result.Error.Error(),
				})
				return
			}

			// Mendapatkan ID data yang baru dibuat
			var lastInsertID uint
			if err := tx.Table("data_siswa").Select("id_data").Where("email = ? AND tahun_ajaran = ?", siswa.Email, input.TahunAjaran).Order("id_data DESC").Limit(1).Scan(&lastInsertID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal mendapatkan ID data yang baru dibuat untuk siswa: " + siswa.Email,
					"error":   err.Error(),
				})
				return
			}
			idData = lastInsertID
		}

		// Menambahkan hasil rekap ke array
		hasilRekap = append(hasilRekap, map[string]interface{}{
			"id_data":      idData,
			"email":        siswa.Email,
			"nama_siswa":   siswa.Nama,
			"id_kelas":     siswa.IDKelas,
			"tahun_ajaran": input.TahunAjaran,
			"jumlah_logs":  len(logs),
		})
	}

	// Create entries for students even if they don't have logs
	if len(hasilRekap) == 0 {
		// Instead of returning an error, we'll create empty entries for all students
		for _, siswa := range siswaList {
			// Create empty progres JSON array
			emptyProgresJSON, _ := json.Marshal([]logsprogres.ProgresItem{})

			// Create new data entry with empty progres
			newData := map[string]interface{}{
				"email":        siswa.Email,
				"id_kelas":     siswa.IDKelas,
				"progres":      string(emptyProgresJSON),
				"tahun_ajaran": input.TahunAjaran,
			}

			result := tx.Table("data_siswa").Create(newData)
			if result.Error != nil {
				continue // Skip this student if there's an error
			}

			// Get the ID of the newly created entry
			var lastInsertID uint
			if err := tx.Table("data_siswa").Select("id_data").Where("email = ? AND tahun_ajaran = ?", siswa.Email, input.TahunAjaran).Order("id_data DESC").Limit(1).Scan(&lastInsertID).Error; err != nil {
				continue // Skip this student if there's an error
			}

			// Add to results
			hasilRekap = append(hasilRekap, map[string]interface{}{
				"id_data":      lastInsertID,
				"email":        siswa.Email,
				"nama_siswa":   siswa.Nama,
				"id_kelas":     siswa.IDKelas,
				"tahun_ajaran": input.TahunAjaran,
				"jumlah_logs":  0,
			})
		}

		// If still no entries were created, return an error
		if len(hasilRekap) == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Tidak ada data siswa yang dapat direkap",
			})
			return
		}
	}

	// Jika diminta untuk menghapus data logs
	if input.DeleteLogsData {
		// Membuat query untuk menghapus logs berdasarkan filter yang sama
		deleteQuery := tx.Table("logs")

		// Menambahkan filter kelas jika ada
		if input.FilterKelas != nil {
			deleteQuery = deleteQuery.Where("id_kelas = ?", *input.FilterKelas)
		}

		// Filter berdasarkan tanggal jika disediakan
		if input.StartDate != "" {
			startDate, _ := time.Parse("2006-01-02", input.StartDate)
			deleteQuery = deleteQuery.Where("created_at >= ?", startDate)
		}

		if input.EndDate != "" {
			endDate, _ := time.Parse("2006-01-02", input.EndDate)
			// Tambahkan 1 hari ke endDate untuk mencakup seluruh hari
			endDate = endDate.Add(24 * time.Hour)
			deleteQuery = deleteQuery.Where("created_at < ?", endDate)
		}

		// Membuat array email siswa yang diproses untuk filter penghapusan
		emails := make([]string, 0, len(hasilRekap))
		for _, rekap := range hasilRekap {
			emails = append(emails, rekap["email"].(string))
		}

		// Menghapus logs hanya untuk siswa yang berhasil direkap
		if len(emails) > 0 {
			if err := deleteQuery.Where("email IN ?", emails).Delete(nil).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal menghapus data logs",
					"error":   err.Error(),
				})
				return
			}
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan perubahan",
			"error":   err.Error(),
		})
		return
	}

	// Mengembalikan respons sukses
	c.JSON(http.StatusOK, gin.H{
		"message":          "Rekap semester untuk semua siswa berhasil dibuat",
		"jumlah_siswa":     len(hasilRekap),
		"total_logs":       totalLogsProcessed,
		"logs_dihapus":     input.DeleteLogsData,
		"tahun_ajaran":     input.TahunAjaran,
		"detail_per_siswa": hasilRekap,
	})
}

// DeleteDataSiswa menghapus data rekap semester berdasarkan ID
func DeleteDataSiswa(c *gin.Context) {
	db := config.DB
	idData := c.Param("id_data") // ID data siswa yang ingin dihapus

	if idData == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID data diperlukan",
		})
		return
	}

	// Memulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Memeriksa apakah data ada sebelum dihapus
	var count int64
	if err := tx.Table("data_siswa").Where("id_data = ?", idData).Count(&count).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memeriksa data siswa",
			"error":   err.Error(),
		})
		return
	}

	// Jika data tidak ditemukan
	if count == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data siswa tidak ditemukan",
		})
		return
	}

	// Hapus data siswa berdasarkan id_data
	result := tx.Table("data_siswa").Where("id_data = ?", idData).Delete(nil)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus data siswa",
			"error":   result.Error.Error(),
		})
		return
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan perubahan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data siswa berhasil dihapus",
		"id_data": idData,
	})
}
