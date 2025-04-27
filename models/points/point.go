package points

type UserPoint struct {
	Email   string `json:"email" gorm:"type:varchar(255);primaryKey"`
	Diamond int    `json:"diamond,omitempty" gorm:"default:0"`
	Exp     int    `json:"exp,omitempty" gorm:"not null;default:0"`

	// Relasi dengan tabel siswa (optional, hanya jika ingin preload relasi)
	// Siswa Siswa `gorm:"foreignKey:Email;references:Email;constraint:OnDelete:CASCADE"`
}

type UserPointResponse struct {
	Email   string `json:"email" gorm:"type:varchar(255);primaryKey"`
	Diamond int    `json:"diamond" gorm:"default:0"`
	Exp     int    `json:"exp" gorm:"not null;default:0"`
}

type DiamondOrExp struct {
	Diamond *int `json:"diamond,omitempty"`
	Exp     *int `json:"exp,omitempty"`
}
