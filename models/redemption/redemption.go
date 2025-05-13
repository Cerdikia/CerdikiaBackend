package redemption

import "time"

// Item represents a single item in the redemption request
type RedemptionItem struct {
	IDBarang int `json:"id_barang"`
	Jumlah   int `json:"jumlah"`
}

// RedemptionRequest represents the request body for redeeming gifts
type RedemptionRequest struct {
	Email string           `json:"email"`
	Items []RedemptionItem `json:"items"`
}

// BarangInfo represents information about a gift item
type BarangInfo struct {
	IDBarang    int    `json:"id_barang" gorm:"column:id_barang"`
	NamaBarang  string `json:"nama_barang" gorm:"column:nama_barang"`
	Img         string `json:"img" gorm:"column:img"`
	Description string `json:"description" gorm:"column:description"`
	Diamond     int64  `json:"diamond" gorm:"column:diamond"`
	Jumlah      int64  `json:"jumlah" gorm:"column:jumlah"`
}

// TableName specifies the table name for BarangInfo
func (BarangInfo) TableName() string {
	return "barang"
}

// LogsPenukaranPoint represents a log entry for a gift redemption
type LogsPenukaranPoint struct {
	IDLog            int       `json:"id_log" gorm:"column:id_log;primaryKey;autoIncrement"`
	IDBarang         int       `json:"id_barang" gorm:"column:id_barang"`
	Email            string    `json:"email" gorm:"column:email"`
	Jumlah           int       `json:"jumlah" gorm:"column:jumlah"`
	TanggalPenukaran time.Time `json:"tanggal_penukaran" gorm:"column:tanggal_penukaran;default:CURRENT_TIMESTAMP"`
	KodePenukaran    string    `json:"kode_penukaran" gorm:"column:kode_penukaran"`
	StatusPenukaran  string    `json:"status_penukaran" gorm:"column:status_penukaran;default:menunggu"`
}

// TableName specifies the table name for LogsPenukaranPoint
func (LogsPenukaranPoint) TableName() string {
	return "logs_penukaran_point"
}

// UserVerified represents the verification status of a user
type UserVerified struct {
	Email          string `json:"email" gorm:"column:email;primaryKey"`
	VerifiedStatus string `json:"verified_status" gorm:"column:verified_status"`
}

// TableName specifies the table name for UserVerified
func (UserVerified) TableName() string {
	return "user_verified"
}

// UserPoints represents the points/diamonds of a user
type UserPoints struct {
	Email   string `json:"email" gorm:"column:email;primaryKey"`
	Diamond int    `json:"diamond" gorm:"column:diamond"`
	Exp     int    `json:"exp" gorm:"column:exp"`
}

// TableName specifies the table name for UserPoints
func (UserPoints) TableName() string {
	return "user_points"
}
