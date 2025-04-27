package hadiah

import "time"

type TableBarang struct {
	IdBarang   uint      `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	NamaBarang string    `gorm:"column:nama_barang" json:"nama_barang"`
	Diamond    int       `json:"diamond"`
	Jumlah     int       `json:"jumlah"`
	Img        string    `gorm:"type:text" json:"img"` // Tambahan kolom img
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TableBarang) TableName() string {
	return "table_barang"
}
