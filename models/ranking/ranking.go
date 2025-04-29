package ranking

type ModuleResponse struct {
	IDModule        int    `json:"id_module"`
	IDMapel         int    `json:"id_mapel"`
	IDKelas         int    `json:"id_kelas"`
	Module          int    `json:"module"`
	ModuleJudul     string `json:"module_judul"`
	ModuleDeskripsi string `json:"module_deskripsi"`
	IsReady         bool   `json:"is_ready"`
	IsCompleted     bool   `json:"is_completed"`
}

type RankingResponse struct {
	Ranking int    `json:"ranking"`
	Nama    string `json:"nama"`
	Kelas   string `json:"kelas"`
	Exp     int    `json:"exp"`
}
