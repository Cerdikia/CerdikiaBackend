package kelas

type Kelas struct {
	IDKelas int    `gorm:"primaryKey;column:id_kelas" json:"id_kelas"`
	Kelas   string `gorm:"type:varchar(50)" json:"kelas"`
}

func (Kelas) TableName() string {
	return "kelas"
}
