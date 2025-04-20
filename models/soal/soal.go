package soal

type UploadImageResponse struct {
	Url string `json:"url"`
}

type UploadSoal struct {
	IDSoal   uint   `gorm:"column:id_soal;primaryKey;autoIncrement" json:"id_soal"`
	IDModule uint   `gorm:"column:id_module;not null" json:"id_module"`
	Soal     string `gorm:"type:longtext;not null" json:"soal"`
	Jenis    string `gorm:"type:enum('pilihan_ganda','essay');default:'pilihan_ganda'" json:"jenis"`
	OpsiA    string `gorm:"type:longtext;not null" json:"opsi_a"`
	OpsiB    string `gorm:"type:longtext;not null" json:"opsi_b"`
	OpsiC    string `gorm:"type:longtext;not null" json:"opsi_c"`
	OpsiD    string `gorm:"type:longtext;not null" json:"opsi_d"`
	Jawaban  string `gorm:"type:varchar(10);not null" json:"jawaban"`
}

// Optional: Jika kamu ingin GORM tahu nama tabelnya secara eksplisit
func (UploadSoal) TableName() string {
	return "soal"
}

type GenericSoalModelResponse struct {
	IDSoal   uint   `gorm:"column:id_soal;primaryKey;autoIncrement" json:"id_soal"`
	IDModule uint   `gorm:"column:id_module;not null" json:"id_module"`
	Soal     string `gorm:"type:longtext;not null" json:"soal"`
	Jenis    string `gorm:"type:enum('pilihan_ganda','essay');default:'pilihan_ganda'" json:"jenis"`
	OpsiA    string `gorm:"type:longtext;not null" json:"opsi_a"`
	OpsiB    string `gorm:"type:longtext;not null" json:"opsi_b"`
	OpsiC    string `gorm:"type:longtext;not null" json:"opsi_c"`
	OpsiD    string `gorm:"type:longtext;not null" json:"opsi_d"`
	Jawaban  string `gorm:"type:varchar(10);not null" json:"jawaban"`
}

type DeletDataSoalResponse struct {
	Message string `json:"message"`
}
