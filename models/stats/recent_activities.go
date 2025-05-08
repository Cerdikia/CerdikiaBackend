package stats

import (
	"time"
)

// RecentActivity adalah struktur untuk menampilkan aktivitas terakhir
type RecentActivity struct {
	IDLogs    uint      `json:"id_logs"`
	Email     string    `json:"email"`
	IDKelas   int       `json:"id_kelas"`
	IDMapel   int       `json:"id_mapel"`
	IDModule  int       `json:"id_module"`
	Skor      int       `json:"skor"`
	CreatedAt time.Time `json:"created_at"`
	// Tambahan informasi yang bisa ditampilkan
	NamaSiswa string `json:"nama_siswa,omitempty"`
	NamaMapel string `json:"nama_mapel,omitempty"`
	NamaKelas string `json:"nama_kelas,omitempty"`
}

// RecentActivitiesResponse adalah struktur untuk respons aktivitas terakhir dengan pagination
type RecentActivitiesResponse struct {
	Activities  []RecentActivity `json:"activities"`
	TotalCount  int64            `json:"total_count"`
	CurrentPage int              `json:"current_page"`
	Limit       int              `json:"limit"`
	TotalPages  int              `json:"total_pages"`
}
