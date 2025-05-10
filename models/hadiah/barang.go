package hadiah

import "time"

type TableBarang struct {
	IdBarang    uint      `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	NamaBarang  string    `gorm:"column:nama_barang" json:"nama_barang"`
	Diamond     int       `json:"diamond" gorm:"column:diamond"`
	Jumlah      int       `json:"jumlah" gorm:"column:jumlah"`
	Img         string    `gorm:"type:text;column:img" json:"img"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

// For form binding
type BarangForm struct {
	Name         string `form:"name" binding:"required"`
	Quantity     int    `form:"quantity" binding:"required"`
	DiamondValue int    `form:"diamond_value" binding:"required"`
	Description  string `form:"description" binding:"required"`
}

func (TableBarang) TableName() string {
	return "barang"
}
