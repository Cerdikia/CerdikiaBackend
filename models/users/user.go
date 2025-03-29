package users

import "time"

type User struct {
	Email       string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Nama        string    `json:"nama" gorm:"type:varchar(100);not null"`
	Kelas       string    `json:"kelas" gorm:"type:varchar(50)"`
	DateCreated time.Time `json:"date_created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	// DateCreated string `json:"date_created" from:"date_created"`
	Role string `json:"role" gorm:"type:enum('siswa','admin','guru');default:'siswa'"`
}

type LoginRequest struct {
	Email string `json:"email" gorm:"type:varchar(100);unique;not null"`
	// Role string `json:"role" gorm:"type:enum('siswa','admin','guru');default:'siswa'"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
