package users

import (
	"time"
)

type Siswa struct {
	Email       string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Nama        string    `json:"nama" gorm:"type:varchar(100)"`
	Kelas       string    `json:"kelas" gorm:"type:varchar(50)"`
	DateCreated time.Time `json:"date_created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName memberikan nama tabel yang eksplisit
func (Siswa) TableName() string {
	return "siswa"
}

type UserProfile struct {
	Email       string  `json:"email"`
	Nama        string  `json:"nama"`
	Role        string  `json:"role"`               // siswa/guru/admin
	Kelas       *string `json:"kelas,omitempty"`    // optional
	IdMapel     *int    `json:"id_mapel,omitempty"` // optional
	Jabatan     *string `json:"jabatan,omitempty"`
	Keterangan  *string `json:"keterangan,omitempty"`
	DateCreated string  `json:"date_created"`
}

type UserProfileReq struct {
	Email       string  `json:"email"`
	Nama        string  `json:"nama"`
	Role        string  `json:"role"`               // siswa/guru/admin
	Kelas       *string `json:"kelas,omitempty"`    // optional
	IdMapel     *int    `json:"id_mapel,omitempty"` // optional
	Jabatan     *string `json:"jabatan,omitempty"`
	Keterangan  *string `json:"keterangan,omitempty"`
	DateCreated *string `json:"date_created"`
}

type CreateAcountReq struct {
	Email string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Role  string `json:"role"`
	// Role string `json:"role" gorm:"type:enum('siswa','admin','guru');default:'siswa'"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Guru struct {
	Email       string    `gorm:"primaryKey;unique;size:100" json:"email"`
	IDMapel     int       `gorm:"not null" json:"id_mapel"`
	Nama        string    `gorm:"size:100;not null" json:"nama"`
	Jabatan     *string   `gorm:"size:100" json:"jabatan,omitempty"` // nullable
	DateCreated time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_created"`
}

// TableName memberikan nama tabel yang eksplisit
func (Guru) TableName() string {
	return "guru"
}

type Admin struct {
	Email       string    `gorm:"primaryKey;unique;size:100" json:"email"`
	Nama        string    `gorm:"size:100;not null" json:"nama"`
	Keterangan  *string   `gorm:"type:text" json:"keterangan,omitempty"` // nullable
	DateCreated time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_created"`
}

// TableName memberikan nama tabel yang eksplisit
func (Admin) TableName() string {
	return "admin"
}
