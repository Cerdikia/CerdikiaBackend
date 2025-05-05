package genericActivities

// // Model untuk tabel mapel
type Mapel struct {
	IDMapel int    `gorm:"primaryKey;column:id_mapel" json:"id_mapel"`
	Mapel   string `gorm:"type:varchar(255);not null" json:"mapel"`
}

func (Mapel) TableName() string {
	return "mapel"
}

// // Model untuk tabel modules
type Module struct {
	IDModule        int    `gorm:"primaryKey;column:id_module" json:"id_module"`
	IDMapel         int    `gorm:"not null" json:"id_mapel"`
	IDKelas         int    `gorm:"not null" json:"id_kelas"`
	Module          *int   `json:"module"` // nullable
	ModuleJudul     string `json:"module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi"`
	IsReady         bool   `gorm:"column:is_ready;default:false" json:"is_ready"`

	// Optional preload relasi
	// Mapel Mapel `gorm:"foreignKey:IDMapel;references:IDMapel" json:"-"`
	// Kelas Kelas `gorm:"foreignKey:IDKelas;references:IDKelas" json:"-"`
}

func (Module) TableName() string {
	return "modules"
}

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
	IsReady         bool   `json:"is_ready" gorm:"column:is_ready"`
}

type SpesifiedModulesKelasResponse struct {
	// id_module, module, module_judul, module_deskripsi
	IDModule        int    `json:"id_module" gorm:"column:id_module"`
	Module          int    `json:"module" gorm:"column:module"`
	ModuleJudul     string `json:"module_judul" gorm:"column:module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi" gorm:"column:module_deskripsi"`
	IsReady         bool   `json:"is_ready" gorm:"column:is_ready"`
	IsCompleted     bool   `json:"is_completed" gorm:"column:is_completed"`
}

type GenericModulesResponse struct {
	// id_module, module, module_judul, module_deskripsi
	Kelas           string `json:"kelas" gorm:"column:kelas"`
	IDModule        int    `json:"id_module" gorm:"column:id_module"`
	Module          int    `json:"module" gorm:"column:module"`
	ModuleJudul     string `json:"module_judul" gorm:"column:module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi" gorm:"column:module_deskripsi"`
	IsReady         bool   `json:"is_ready" gorm:"column:is_ready"`
}

type GenericKelasResponse struct {
	// id_module, module, module_judul, module_deskripsi
	IDModule        int    `json:"id_module" gorm:"column:id_module"`
	Module          int    `json:"module" gorm:"column:module"`
	ModuleJudul     string `json:"module_judul" gorm:"column:module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi" gorm:"column:module_deskripsi"`
	IsReady         bool   `json:"is_ready" gorm:"column:is_ready"`
	Mapel           string `gorm:"column:mapel;type:varchar(255);not null" json:"mapel"`
}
