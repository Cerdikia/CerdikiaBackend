package repositories

import (
	"coba1BE/config"
	"coba1BE/models/genericActivities"
	"coba1BE/models/logsprogres"
	"coba1BE/models/soal"
	"coba1BE/models/stats"
	"math"
)

// GetAllStats mengambil semua statistik aplikasi
func GetAllStats() stats.StatsResponse {
	db := config.DB
	var statsResponse stats.StatsResponse

	// Hitung total users (gabungan siswa, guru, admin)
	db.Raw(`
		SELECT 
			COUNT(*) as total_users 
		FROM (
			SELECT email FROM siswa
			UNION
			SELECT email FROM guru
			UNION
			SELECT email FROM admin
		) as all_users
	`).Scan(&statsResponse.TotalUsers)

	// Hitung total siswa
	db.Table("siswa").Count(&statsResponse.TotalSiswa)

	// Hitung total guru
	db.Table("guru").Count(&statsResponse.TotalGuru)

	// Hitung total admin
	db.Table("admin").Count(&statsResponse.TotalAdmin)

	// Hitung total mapel
	db.Model(&genericActivities.Mapel{}).Count(&statsResponse.TotalMapel)

	// Hitung total module
	db.Model(&genericActivities.Module{}).Count(&statsResponse.TotalModule)

	// Hitung total soal
	db.Model(&soal.UploadSoal{}).Count(&statsResponse.TotalSoal)

	// Hitung total kelas
	db.Table("kelas").Count(&statsResponse.TotalKelas)

	// Hitung total logs
	db.Table("logs").Count(&statsResponse.TotalLogs)

	// Hitung total barang
	db.Table("barang").Count(&statsResponse.TotalBarang)

	// Hitung total data siswa (rekap semester)
	// db.Table("data_siswa").Count(&statsResponse.TotalDataSiswa)

	return statsResponse
}

// GetRecentActivities mengambil aktivitas terakhir dari tabel logs dengan pagination
func GetRecentActivities(page, limit int) stats.RecentActivitiesResponse {
	db := config.DB
	var response stats.RecentActivitiesResponse
	var totalCount int64

	// Set nilai default untuk pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 7 // Default limit 7 data
	}

	// Hitung total logs
	db.Table("logs").Count(&totalCount)

	// Hitung offset
	offset := (page - 1) * limit

	// Ambil data logs dengan pagination dan join dengan tabel lain untuk mendapatkan informasi tambahan
	var activities []stats.RecentActivity

	// Query dasar untuk mendapatkan logs dengan informasi tambahan
	query := `
		SELECT 
			l.id_logs, 
			l.email, 
			l.id_kelas, 
			l.id_mapel, 
			l.id_module, 
			l.skor, 
			l.created_at,
			s.nama as nama_siswa,
			m.mapel as nama_mapel,
			k.kelas as nama_kelas
		FROM logs l
		LEFT JOIN siswa s ON l.email = s.email
		LEFT JOIN mapel m ON l.id_mapel = m.id_mapel
		LEFT JOIN kelas k ON l.id_kelas = k.id_kelas
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?
	`

	db.Raw(query, limit, offset).Scan(&activities)

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Buat response
	response = stats.RecentActivitiesResponse{
		Activities:  activities,
		TotalCount:  totalCount,
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  totalPages,
	}

	return response
}

// GetRecentActivitiesWithLogs mengambil aktivitas terakhir dari tabel logs dengan pagination
// dan juga mengembalikan objek logs asli
func GetRecentActivitiesWithLogs(page, limit int) (stats.RecentActivitiesResponse, []logsprogres.Log) {
	db := config.DB
	var response stats.RecentActivitiesResponse
	var logs []logsprogres.Log
	var totalCount int64

	// Set nilai default untuk pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 7 // Default limit 7 data
	}

	// Hitung total logs
	db.Table("logs").Count(&totalCount)

	// Hitung offset
	offset := (page - 1) * limit

	// Ambil data logs original
	db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs)

	// Ambil data logs dengan informasi tambahan
	var activities []stats.RecentActivity

	// Query untuk mendapatkan logs dengan informasi tambahan
	query := `
		SELECT 
			l.id_logs, 
			l.email, 
			l.id_kelas, 
			l.id_mapel, 
			l.id_module, 
			l.skor, 
			l.created_at,
			s.nama as nama_siswa,
			m.mapel as nama_mapel,
			k.nama_kelas as nama_kelas
		FROM logs l
		LEFT JOIN siswa s ON l.email = s.email
		LEFT JOIN mapel m ON l.id_mapel = m.id_mapel
		LEFT JOIN kelas k ON l.id_kelas = k.id_kelas
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?
	`

	db.Raw(query, limit, offset).Scan(&activities)

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Buat response
	response = stats.RecentActivitiesResponse{
		Activities:  activities,
		TotalCount:  totalCount,
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  totalPages,
	}

	return response, logs
}
