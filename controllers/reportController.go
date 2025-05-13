package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/logsprogres"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetScoreReport menampilkan laporan nilai siswa berdasarkan filter tertentu
func GetScoreReport(c *gin.Context) {
	// Mendapatkan parameter filter
	idKelas := c.Query("id_kelas")
	idMapel := c.Query("id_mapel")
	sortBy := c.DefaultQuery("sort_by", "latest") // latest, highest
	aggregateBy := c.DefaultQuery("aggregate_by", "first") // first, highest

	// Validasi parameter
	if idKelas == "" && idMapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Minimal satu parameter filter (id_kelas atau id_mapel) harus diisi",
			Data:    nil,
		})
		return
	}

	// Inisialisasi database
	db := config.DB

	// Membuat query dasar
	query := db.Model(&logsprogres.Log{})

	// Menambahkan filter berdasarkan parameter
	if idKelas != "" {
		query = query.Where("id_kelas = ?", idKelas)
	}

	if idMapel != "" {
		query = query.Where("id_mapel = ?", idMapel)
	}

	// Menentukan urutan berdasarkan sortBy
	switch sortBy {
	case "highest":
		query = query.Order("skor DESC")
	case "latest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Mengambil semua log yang sesuai dengan filter
	var logs []logsprogres.Log
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Jika tidak ada data yang ditemukan
	if len(logs) == 0 {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Tidak ada data yang ditemukan",
			Data:    []logsprogres.Log{},
		})
		return
	}

	// Jika aggregateBy adalah "first" atau "highest", kita perlu mengelompokkan data
	if aggregateBy == "first" || aggregateBy == "highest" {
		// Map untuk menyimpan data yang sudah dikelompokkan berdasarkan email siswa
		aggregatedLogs := make(map[string]map[int]logsprogres.Log)

		// Iterasi semua log dan kelompokkan berdasarkan email dan id_mapel
		for _, log := range logs {
			// Jika email belum ada di map, buat map baru untuk email tersebut
			if _, ok := aggregatedLogs[log.Email]; !ok {
				aggregatedLogs[log.Email] = make(map[int]logsprogres.Log)
			}

			// Jika aggregateBy adalah "first", kita hanya menyimpan log pertama untuk setiap mapel
			if aggregateBy == "first" {
				// Jika id_mapel belum ada di map atau log ini lebih baru, simpan log ini
				if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.CreatedAt.Before(existingLog.CreatedAt) {
					aggregatedLogs[log.Email][log.IDMapel] = log
				}
			} else if aggregateBy == "highest" {
				// Jika id_mapel belum ada di map atau skor log ini lebih tinggi, simpan log ini
				if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.Skor > existingLog.Skor {
					aggregatedLogs[log.Email][log.IDMapel] = log
				}
			}
		}

		// Konversi map ke slice untuk respons
		resultLogs := []logsprogres.Log{}
		for _, mapelLogs := range aggregatedLogs {
			for _, log := range mapelLogs {
				resultLogs = append(resultLogs, log)
			}
		}

		// Mengembalikan hasil yang sudah diagregasi
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Data laporan nilai berhasil diambil",
			Data:    resultLogs,
		})
		return
	}

	// Jika tidak ada agregasi, kembalikan semua log
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data laporan nilai berhasil diambil",
		Data:    logs,
	})
}

// GetScoreReportSummary menampilkan ringkasan laporan nilai siswa
func GetScoreReportSummary(c *gin.Context) {
	// Mendapatkan parameter filter
	idKelas := c.Query("id_kelas")
	idMapel := c.Query("id_mapel")
	aggregateBy := c.DefaultQuery("aggregate_by", "highest") // first, highest

	// Validasi parameter
	if idKelas == "" && idMapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Minimal satu parameter filter (id_kelas atau id_mapel) harus diisi",
			Data:    nil,
		})
		return
	}

	// Inisialisasi database
	db := config.DB

	// Membuat query dasar
	query := db.Model(&logsprogres.Log{})

	// Menambahkan filter berdasarkan parameter
	if idKelas != "" {
		query = query.Where("id_kelas = ?", idKelas)
	}

	if idMapel != "" {
		query = query.Where("id_mapel = ?", idMapel)
	}

	// Mengambil semua log yang sesuai dengan filter
	var logs []logsprogres.Log
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Jika tidak ada data yang ditemukan
	if len(logs) == 0 {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Tidak ada data yang ditemukan",
			Data:    map[string]interface{}{
				"total_siswa":  0,
				"nilai_rata2":  0,
				"nilai_min":    0,
				"nilai_max":    0,
				"distribusi": []map[string]interface{}{},
			},
		})
		return
	}

	// Map untuk menyimpan data yang sudah dikelompokkan berdasarkan email siswa
	aggregatedLogs := make(map[string]map[int]logsprogres.Log)

	// Iterasi semua log dan kelompokkan berdasarkan email dan id_mapel
	for _, log := range logs {
		// Jika email belum ada di map, buat map baru untuk email tersebut
		if _, ok := aggregatedLogs[log.Email]; !ok {
			aggregatedLogs[log.Email] = make(map[int]logsprogres.Log)
		}

		// Jika aggregateBy adalah "first", kita hanya menyimpan log pertama untuk setiap mapel
		if aggregateBy == "first" {
			// Jika id_mapel belum ada di map atau log ini lebih baru, simpan log ini
			if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.CreatedAt.Before(existingLog.CreatedAt) {
				aggregatedLogs[log.Email][log.IDMapel] = log
			}
		} else if aggregateBy == "highest" {
			// Jika id_mapel belum ada di map atau skor log ini lebih tinggi, simpan log ini
			if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.Skor > existingLog.Skor {
				aggregatedLogs[log.Email][log.IDMapel] = log
			}
		}
	}

	// Konversi map ke slice untuk analisis
	resultLogs := []logsprogres.Log{}
	for _, mapelLogs := range aggregatedLogs {
		for _, log := range mapelLogs {
			resultLogs = append(resultLogs, log)
		}
	}

	// Hitung statistik
	totalSiswa := len(aggregatedLogs)
	totalNilai := 0
	nilaiMin := 100 // Asumsi nilai maksimal adalah 100
	nilaiMax := 0

	// Map untuk distribusi nilai
	distribusi := map[string]int{
		"A (90-100)": 0,
		"B (80-89)":  0,
		"C (70-79)":  0,
		"D (60-69)":  0,
		"E (0-59)":   0,
	}

	// Hitung total nilai dan temukan nilai min/max
	for _, log := range resultLogs {
		totalNilai += log.Skor

		if log.Skor < nilaiMin {
			nilaiMin = log.Skor
		}

		if log.Skor > nilaiMax {
			nilaiMax = log.Skor
		}

		// Hitung distribusi nilai
		switch {
		case log.Skor >= 90:
			distribusi["A (90-100)"]++
		case log.Skor >= 80:
			distribusi["B (80-89)"]++
		case log.Skor >= 70:
			distribusi["C (70-79)"]++
		case log.Skor >= 60:
			distribusi["D (60-69)"]++
		default:
			distribusi["E (0-59)"]++
		}
	}

	// Hitung rata-rata nilai
	nilaiRataRata := 0
	if len(resultLogs) > 0 {
		nilaiRataRata = totalNilai / len(resultLogs)
	}

	// Konversi distribusi ke format yang lebih baik untuk respons
	distribusiArray := []map[string]interface{}{}
	for kategori, jumlah := range distribusi {
		distribusiArray = append(distribusiArray, map[string]interface{}{
			"kategori": kategori,
			"jumlah":   jumlah,
		})
	}

	// Mengembalikan hasil ringkasan
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Ringkasan laporan nilai berhasil diambil",
		Data: map[string]interface{}{
			"total_siswa":  totalSiswa,
			"nilai_rata2":  nilaiRataRata,
			"nilai_min":    nilaiMin,
			"nilai_max":    nilaiMax,
			"distribusi":   distribusiArray,
		},
	})
}

// GetStudentScoreComparison membandingkan nilai siswa dengan rata-rata kelas
func GetStudentScoreComparison(c *gin.Context) {
	// Mendapatkan parameter
	email := c.Query("email")
	idKelas := c.Query("id_kelas")
	idMapel := c.Query("id_mapel")
	aggregateBy := c.DefaultQuery("aggregate_by", "highest") // first, highest

	// Validasi parameter
	if email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter email harus diisi",
			Data:    nil,
		})
		return
	}

	if idKelas == "" && idMapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Minimal satu parameter filter (id_kelas atau id_mapel) harus diisi",
			Data:    nil,
		})
		return
	}

	// Inisialisasi database
	db := config.DB

	// Membuat query dasar untuk siswa tertentu
	queryStudent := db.Model(&logsprogres.Log{}).Where("email = ?", email)

	// Membuat query dasar untuk semua siswa (untuk perbandingan)
	queryAll := db.Model(&logsprogres.Log{})

	// Menambahkan filter berdasarkan parameter
	if idKelas != "" {
		queryStudent = queryStudent.Where("id_kelas = ?", idKelas)
		queryAll = queryAll.Where("id_kelas = ?", idKelas)
	}

	if idMapel != "" {
		queryStudent = queryStudent.Where("id_mapel = ?", idMapel)
		queryAll = queryAll.Where("id_mapel = ?", idMapel)
	}

	// Mengambil log untuk siswa tertentu
	var studentLogs []logsprogres.Log
	if err := queryStudent.Find(&studentLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data siswa: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Mengambil log untuk semua siswa
	var allLogs []logsprogres.Log
	if err := queryAll.Find(&allLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data semua siswa: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Jika tidak ada data yang ditemukan untuk siswa
	if len(studentLogs) == 0 {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Tidak ada data yang ditemukan untuk siswa ini",
			Data:    map[string]interface{}{
				"student_scores": []interface{}{},
				"class_average":  0,
				"comparison":    "Tidak ada data",
			},
		})
		return
	}

	// Agregasi data siswa berdasarkan mata pelajaran
	studentScores := make(map[int]int)
	for _, log := range studentLogs {
		if aggregateBy == "highest" {
			if existingScore, ok := studentScores[log.IDMapel]; !ok || log.Skor > existingScore {
				studentScores[log.IDMapel] = log.Skor
			}
		} else { // first
			if _, ok := studentScores[log.IDMapel]; !ok {
				studentScores[log.IDMapel] = log.Skor
			}
		}
	}

	// Agregasi data semua siswa berdasarkan mata pelajaran
	allScores := make(map[int][]int)
	for _, log := range allLogs {
		if _, ok := allScores[log.IDMapel]; !ok {
			allScores[log.IDMapel] = []int{}
		}
		allScores[log.IDMapel] = append(allScores[log.IDMapel], log.Skor)
	}

	// Hitung rata-rata kelas untuk setiap mata pelajaran
	classAverages := make(map[int]float64)
	for mapelID, scores := range allScores {
		total := 0
		for _, score := range scores {
			total += score
		}
		classAverages[mapelID] = float64(total) / float64(len(scores))
	}

	// Buat respons perbandingan
	comparisons := []map[string]interface{}{}
	for mapelID, studentScore := range studentScores {
		classAvg := classAverages[mapelID]
		difference := float64(studentScore) - classAvg
		status := ""

		if difference > 0 {
			status = "Di atas rata-rata"
		} else if difference < 0 {
			status = "Di bawah rata-rata"
		} else {
			status = "Sama dengan rata-rata"
		}

		comparisons = append(comparisons, map[string]interface{}{
			"id_mapel":      mapelID,
			"student_score": studentScore,
			"class_average": classAvg,
			"difference":   difference,
			"status":       status,
		})
	}

	// Mengembalikan hasil perbandingan
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Perbandingan nilai siswa berhasil diambil",
		Data:    comparisons,
	})
}

// GetScoreProgress mendapatkan progres nilai siswa dari waktu ke waktu
func GetScoreProgress(c *gin.Context) {
	// Mendapatkan parameter
	email := c.Query("email")
	idMapel := c.Query("id_mapel")

	// Validasi parameter
	if email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter email harus diisi",
			Data:    nil,
		})
		return
	}

	if idMapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter id_mapel harus diisi",
			Data:    nil,
		})
		return
	}

	// Konversi id_mapel ke integer
	mapelID, err := strconv.Atoi(idMapel)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter id_mapel harus berupa angka",
			Data:    nil,
		})
		return
	}

	// Inisialisasi database
	db := config.DB

	// Mengambil log untuk siswa dan mapel tertentu, diurutkan berdasarkan waktu
	var logs []logsprogres.Log
	if err := db.Where("email = ? AND id_mapel = ?", email, mapelID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Jika tidak ada data yang ditemukan
	if len(logs) == 0 {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Tidak ada data progres yang ditemukan",
			Data:    []interface{}{},
		})
		return
	}

	// Format data untuk respons
	progress := []map[string]interface{}{}
	for _, log := range logs {
		progress = append(progress, map[string]interface{}{
			"id_logs":    log.IDLogs,
			"email":      log.Email,
			"id_mapel":   log.IDMapel,
			"id_module":  log.IDModule,
			"skor":       log.Skor,
			"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Mengembalikan hasil progres
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data progres nilai berhasil diambil",
		Data:    progress,
	})
}

// GetAllStudentsReport mendapatkan laporan semua siswa termasuk yang belum mengerjakan mata pelajaran
func GetAllStudentsReport(c *gin.Context) {
	// Mendapatkan parameter filter
	idKelas := c.Query("id_kelas")
	idMapel := c.Query("id_mapel")
	aggregateBy := c.DefaultQuery("aggregate_by", "highest") // first, highest

	// Validasi parameter
	if idKelas == "" && idMapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Minimal satu parameter filter (id_kelas atau id_mapel) harus diisi",
			Data:    nil,
		})
		return
	}

	// Inisialisasi database
	db := config.DB

	// Membuat query untuk mendapatkan semua siswa
	var siswaQuery string
	var queryParams []interface{}

	if idKelas != "" {
		siswaQuery = "SELECT s.email, s.nama, s.id_kelas, k.kelas FROM siswa s " +
			"LEFT JOIN kelas k ON s.id_kelas = k.id_kelas " +
			"WHERE s.id_kelas = ?"
		queryParams = append(queryParams, idKelas)
	} else {
		siswaQuery = "SELECT s.email, s.nama, s.id_kelas, k.kelas FROM siswa s " +
			"LEFT JOIN kelas k ON s.id_kelas = k.id_kelas"
	}

	type SiswaData struct {
		Email   string `gorm:"column:email" json:"email"`
		Nama    string `gorm:"column:nama" json:"nama"`
		IDKelas int    `gorm:"column:id_kelas" json:"id_kelas"`
		Kelas   string `gorm:"column:kelas" json:"kelas"`
	}

	var siswaList []SiswaData
	if err := db.Raw(siswaQuery, queryParams...).Scan(&siswaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data siswa: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Jika tidak ada siswa yang ditemukan
	if len(siswaList) == 0 {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: "Tidak ada data siswa yang ditemukan",
			Data:    []interface{}{},
		})
		return
	}

	// Membuat query untuk mendapatkan log nilai
	var logsQuery string
	queryParams = []interface{}{} // Reset params

	if idMapel != "" {
		logsQuery = "SELECT * FROM logs WHERE id_mapel = ?"
		queryParams = append(queryParams, idMapel)

		if idKelas != "" {
			logsQuery += " AND id_kelas = ?"
			queryParams = append(queryParams, idKelas)
		}
	} else if idKelas != "" {
		logsQuery = "SELECT * FROM logs WHERE id_kelas = ?"
		queryParams = append(queryParams, idKelas)
	} else {
		logsQuery = "SELECT * FROM logs"
	}

	var logs []logsprogres.Log
	if err := db.Raw(logsQuery, queryParams...).Scan(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal mengambil data logs: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Map untuk menyimpan data yang sudah dikelompokkan berdasarkan email siswa
	aggregatedLogs := make(map[string]map[int]logsprogres.Log)

	// Iterasi semua log dan kelompokkan berdasarkan email dan id_mapel
	for _, log := range logs {
		// Jika email belum ada di map, buat map baru untuk email tersebut
		if _, ok := aggregatedLogs[log.Email]; !ok {
			aggregatedLogs[log.Email] = make(map[int]logsprogres.Log)
		}

		// Jika aggregateBy adalah "first", kita hanya menyimpan log pertama untuk setiap mapel
		if aggregateBy == "first" {
			// Jika id_mapel belum ada di map atau log ini lebih baru, simpan log ini
			if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.CreatedAt.Before(existingLog.CreatedAt) {
				aggregatedLogs[log.Email][log.IDMapel] = log
			}
		} else if aggregateBy == "highest" {
			// Jika id_mapel belum ada di map atau skor log ini lebih tinggi, simpan log ini
			if existingLog, ok := aggregatedLogs[log.Email][log.IDMapel]; !ok || log.Skor > existingLog.Skor {
				aggregatedLogs[log.Email][log.IDMapel] = log
			}
		}
	}

	// Membuat respons untuk semua siswa
	type StudentReportItem struct {
		Email       string                 `json:"email"`
		Nama        string                 `json:"nama"`
		IDKelas     int                    `json:"id_kelas"`
		Kelas       string                 `json:"kelas"`
		HasActivity bool                   `json:"has_activity"`
		Scores      []map[string]interface{} `json:"scores,omitempty"`
	}

	result := []StudentReportItem{}

	// Iterasi semua siswa
	for _, siswa := range siswaList {
		// Membuat item laporan untuk siswa ini
		reportItem := StudentReportItem{
			Email:       siswa.Email,
			Nama:        siswa.Nama,
			IDKelas:     siswa.IDKelas,
			Kelas:       siswa.Kelas,
			HasActivity: false,
			Scores:      []map[string]interface{}{},
		}

		// Jika siswa memiliki log
		if logMap, ok := aggregatedLogs[siswa.Email]; ok && len(logMap) > 0 {
			reportItem.HasActivity = true

			// Tambahkan skor untuk setiap mapel
			for _, log := range logMap {
				reportItem.Scores = append(reportItem.Scores, map[string]interface{}{
					"id_logs":    log.IDLogs,
					"id_mapel":   log.IDMapel,
					"id_module":  log.IDModule,
					"skor":       log.Skor,
					"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		} else if idMapel != "" {
			// Jika ada filter mapel, tambahkan skor kosong untuk mapel tersebut
			mapelID, _ := strconv.Atoi(idMapel)
			reportItem.Scores = append(reportItem.Scores, map[string]interface{}{
				"id_mapel":  mapelID,
				"skor":      0,
				"status":    "Belum mengerjakan",
			})
		}

		// Tambahkan ke hasil
		result = append(result, reportItem)
	}

	// Mengembalikan hasil
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data laporan semua siswa berhasil diambil",
		Data:    result,
	})
}
