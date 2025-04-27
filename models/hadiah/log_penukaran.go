package hadiah

import "time"

type LogsPenukaranPoint struct {
	IdLog            uint       `gorm:"primaryKey;column:id_log" json:"id_log"`
	IdBarang         uint       `gorm:"column:id_barang" json:"id_barang"`
	Email            string     `gorm:"type:varchar(100)" json:"email"`
	Jumlah           int        `json:"jumlah"`
	TanggalPenukaran *time.Time `gorm:"column:tanggal_penukaran;autoCreateTime" json:"tanggal_penukaran"`
	// TanggalPenukaran time.Time `gorm:"column:tanggal_penukaran;autoCreateTime" json:"tanggal_penukaran"`

}

func (LogsPenukaranPoint) TableName() string {
	return "logs_penukaran_point"
}

type PenukaranInput struct {
	IdBarang uint   `json:"id_barang"`
	Email    string `json:"email"`
	Jumlah   int    `json:"jumlah"`
}
