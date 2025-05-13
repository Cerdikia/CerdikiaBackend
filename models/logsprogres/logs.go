package logsprogres

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	IDLogs    uint      `gorm:"primaryKey" json:"id_logs"`
	Email     string    `gorm:"not null" json:"email"`
	IDKelas   int       `gorm:"not null" json:"id_kelas"`
	IDMapel   int       `gorm:"not null" json:"id_mapel"`
	IDModule  int       `gorm:"not null" json:"id_module"`
	Skor      int       `gorm:"not null" json:"skor"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Optional: Jika kamu ingin GORM tahu nama tabelnya secara eksplisit
func (Log) TableName() string {
	return "logs"
}

type ProgresItem struct {
	IDMapel   int    `json:"id_mapel"`
	IDModule  int    `json:"id_module"`
	Skor      int    `json:"skor"`
	CreatedAt string `json:"created_at"`
}

type LogEntry struct {
	Email     string
	IDKelas   int
	IDMapel   int
	IDModule  int
	Skor      int
	CreatedAt string
}

type DataSiswa struct {
	IDData      uint `gorm:"column:id_data;primaryKey;" json:"id_data"`
	Email       string
	IDKelas     int
	Progres     datatypes.JSON `gorm:"type:json"`
	TahunAjaran string
}

// Optional: Jika kamu ingin GORM tahu nama tabelnya secara eksplisit
func (DataSiswa) TableName() string {
	return "data_siswa"
}
