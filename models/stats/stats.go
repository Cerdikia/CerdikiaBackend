package stats

// StatsResponse adalah struktur untuk respons statistik aplikasi
type StatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	TotalSiswa     int64 `json:"total_siswa"`
	TotalGuru      int64 `json:"total_guru"`
	TotalAdmin     int64 `json:"total_admin"`
	TotalMapel     int64 `json:"total_mapel"`
	TotalModule    int64 `json:"total_module"`
	TotalSoal      int64 `json:"total_soal"`
	TotalKelas     int64 `json:"total_kelas"`
	TotalLogs      int64 `json:"total_logs"`
	TotalBarang    int64 `json:"total_barang"`
	TotalDataSiswa int64 `json:"total_data_siswa"`
}
