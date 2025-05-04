package users

import (
	"time"
)

type Siswa struct {
	Email        string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Nama         string    `json:"nama,omitempty" gorm:"type:varchar(100)"`
	IdKelas      *int      `json:"id_kelas,omitempty" gorm:"type:int"`
	DateCreated  time.Time `json:"date_created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	ImageProfile string    `json:"image_profile,omitempty" gorm:"type:text"`
}

// TableName memberikan nama tabel yang eksplisit
func (Siswa) TableName() string {
	return "siswa"
}

type UserProfile struct {
	ID           *int    `json:"id,omitempty"` // optional
	Email        string  `json:"email"`
	Nama         string  `json:"nama"`
	Role         string  `json:"role"`               // siswa/guru/admin
	IdKelas      *int    `json:"id_kelas,omitempty"` // optional
	IdMapel      *int    `json:"id_mapel,omitempty"` // optional
	Jabatan      *string `json:"jabatan,omitempty"`
	Keterangan   *string `json:"keterangan,omitempty"`
	DateCreated  string  `json:"date_created"`
	ImageProfile string  `json:"image_profile"`
}

type UserProfileReq struct {
	Email       string  `json:"email"`
	Nama        string  `json:"nama"`
	Role        string  `json:"role"`               // siswa/guru/admin
	IdKelas     *int    `json:"id_kelas,omitempty"` // optional
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
	Email        string  `json:"email"`
	Nama         string  `json:"nama"`
	Role         string  `json:"role"`               // siswa/guru/admin
	IdKelas      *int    `json:"id_kelas,omitempty"` // optional
	IdMapel      *int    `json:"id_mapel,omitempty"` // optional
	Jabatan      *string `json:"jabatan,omitempty"`
	Keterangan   *string `json:"keterangan,omitempty"`
	DateCreated  string  `json:"date_created"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Guru struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Nama         string    `gorm:"size:100;not null" json:"nama"`
	Jabatan      *string   `gorm:"size:100" json:"jabatan,omitempty"` // nullable
	DateCreated  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_created"`
	ImageProfile string    `json:"image_profile,omitempty" gorm:"type:text"`
}

// TableName memberikan nama tabel yang eksplisit
func (Guru) TableName() string {
	return "guru"
}

// ==========================================

type Mapel struct {
	IDMapel int    `json:"id_mapel"`
	Mapel   string `json:"mapel"`
}

type GuruMapelResponse struct {
	Email   string  `json:"email"`
	Nama    string  `json:"nama"`
	Jabatan string  `json:"jabatan"`
	Mapel   []Mapel `json:"mapel" gorm:"-"` // <== HARUS ada ini
}

// ============================================
type GuruMapel struct {
	IDGuru  uint `gorm:"column:id_guru;primaryKey"`
	IDMapel uint `gorm:"column:id_mapel;primaryKey"`
}

func (GuruMapel) TableName() string {
	return "guru_mapel"
}

type GuruMapelRequest struct {
	IDGuru   uint   `json:"id_guru" binding:"required"`
	IDMapels []uint `json:"id_mapels" binding:"required"`
}

type Admin struct {
	Email        string    `gorm:"primaryKey;unique;size:100" json:"email"`
	Nama         string    `gorm:"size:100;not null" json:"nama"`
	Keterangan   *string   `gorm:"type:text" json:"keterangan,omitempty"` // nullable
	DateCreated  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_created"`
	ImageProfile string    `json:"image_profile,omitempty" gorm:"type:text"`
}

// TableName memberikan nama tabel yang eksplisit
func (Admin) TableName() string {
	return "admin"
}

type UserVerified struct {
	Email    string `gorm:"primaryKey;size:100;not null"`
	Verified bool   `gorm:"default:false"`
}

// TableName memberikan nama tabel yang eksplisit
func (UserVerified) TableName() string {
	return "user_verified"
}
