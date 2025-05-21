package soal

// ImportedModuleData represents a module with its questions
type ImportedModuleData struct {
	JudulModule     string `json:"judul_module"`
	DeskripsiModule string `json:"deskripsi_module"`
	Soal            []ImportedSoalData `json:"soal"`
}

// ImportedSoalData represents a question with its options and answer
type ImportedSoalData struct {
	Soal    string `json:"soal"`
	Jenis   string `json:"jenis"`
	OpsiA   string `json:"opsi_a"`
	OpsiB   string `json:"opsi_b"`
	OpsiC   string `json:"opsi_c"`
	OpsiD   string `json:"opsi_d"`
	Jawaban string `json:"jawaban"`
}

// ImportedData represents the structure of imported data for a class and subject
type ImportedData struct {
	Kelas  string             `json:"kelas"`
	Mapel  string             `json:"mapel"`
	Module []ImportedModuleData `json:"module"`
}

// SaveImportedDataResponse represents the response after saving imported data
type SaveImportedDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ModulesCreated int `json:"modules_created"`
		SoalCreated    int `json:"soal_created"`
	} `json:"data"`
}
