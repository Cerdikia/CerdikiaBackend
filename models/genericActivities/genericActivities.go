package genericactivities

// // Model untuk tabel mapel
// type Mapel struct {
// 		IDMapel int    `gorm:"column:id_mapel;primaryKey"`
// 		Mapel   string `gorm:"column:mapel"`
// }

// // Model untuk tabel modules
// type Module struct {
// 		IDModule        int    `gorm:"column:id_module;primaryKey;autoIncrement"`
// 		IDMapel         int    `gorm:"column:id_mapel"`
// 		Kelas           string `gorm:"column:kelas;type:varchar(50)"`
// 		Module          int    `gorm:"column:module"`
// 		ModuleJudul     string `gorm:"column:module_judul;type:text"`
// 		ModuleDeskripsi string `gorm:"column:module_deskripsi;type:text"`

// 		// Relasi ke tabel Mapel
// 		MataPelajaran Mapel `gorm:"foreignKey:IDMapel"`
// }

// Model untuk hasil query jumlah modul per mapel
type GenericActivitiesResponse struct {
	IDMapel     int    `json:"id_mapel" gorm:"column:id_mapel"`
	NamaMapel   string `json:"nama_mapel" gorm:"column:nama_mapel"`
	Kelas       string `json:"kelas" gorm:"column:kelas"`
	JumlahModul int    `json:"jumlah_modul" gorm:"column:jumlah_modul"`
}

type GenericModulesKelasResponse struct {
	// id_module, module, module_judul, module_deskripsi
	IDModule        int    `json:"id_module" gorm:"column:id_module"`
	Module          int    `json:"module" gorm:"column:module"`
	ModuleJudul     string `json:"module_judul" gorm:"column:module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi" gorm:"column:module_deskripsi"`
}

type GenericModulesResponse struct {
	// id_module, module, module_judul, module_deskripsi
	Kelas           string `json:"kelas" gorm:"column:kelas"`
	IDModule        int    `json:"id_module" gorm:"column:id_module"`
	Module          int    `json:"module" gorm:"column:module"`
	ModuleJudul     string `json:"module_judul" gorm:"column:module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi" gorm:"column:module_deskripsi"`
}
