package repositories

import (
	"coba1BE/config"
	"coba1BE/models/soal"
	"fmt"
)

// SaveImportedData saves the imported data to the database
func SaveImportedData(data []soal.ImportedData) (*soal.SaveImportedDataResponse, string) {
	db := config.DB
	response := &soal.SaveImportedDataResponse{
		Success: true,
		Message: "Data imported successfully",
	}

	totalModulesCreated := 0
	totalSoalCreated := 0

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, importedData := range data {
		// Get kelas ID
		var kelasID uint
		kelasIDQuery := "SELECT id_kelas FROM kelas WHERE kelas = ?"
		if err := tx.Raw(kelasIDQuery, importedData.Kelas).Scan(&kelasID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Sprintf("Error finding kelas: %s", err.Error())
		}
		if kelasID == 0 {
			tx.Rollback()
			return nil, fmt.Sprintf("Kelas '%s' not found", importedData.Kelas)
		}

		// Check if mapel exists (don't create if not exists)
		var mapelID uint
		mapelIDQuery := "SELECT id_mapel FROM mapel WHERE LOWER(mapel) = LOWER(?);"
		if err := tx.Raw(mapelIDQuery, importedData.Mapel).Scan(&mapelID).Error; err != nil {
			fmt.Println("Error finding mapel:", err)
			tx.Rollback()
			return nil, fmt.Sprintf("Error finding mapel: %s", err.Error())
		}

		if mapelID == 0 {
			// Mapel doesn't exist, return error
			tx.Rollback()
			return nil, fmt.Sprintf("Mapel '%s' tidak ditemukan. Silahkan buat mapel terlebih dahulu sebelum mengimpor data.", importedData.Mapel)
		}

		// Process modules
		for _, moduleData := range importedData.Module {
			// Create module
			createModuleQuery := "INSERT INTO modules (id_mapel, id_kelas, module_judul, module_deskripsi ) VALUES (?, ?, ?, ?)"
			moduleResult := tx.Exec(createModuleQuery, mapelID, kelasID, moduleData.JudulModule, moduleData.DeskripsiModule)
			if moduleResult.Error != nil {
				tx.Rollback()
				return nil, fmt.Sprintf("Error creating module: %s", moduleResult.Error.Error())
			}
			totalModulesCreated++

			// Get the newly created module ID
			var moduleID uint
			if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&moduleID).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Sprintf("Error getting new module ID: %s", err.Error())
			}

			// Process soal
			for _, soalData := range moduleData.Soal {
				// Create soal
				newSoal := soal.UploadSoal{
					IDModule: moduleID,
					Soal:     soalData.Soal,
					Jenis:    soalData.Jenis,
					OpsiA:    soalData.OpsiA,
					OpsiB:    soalData.OpsiB,
					OpsiC:    soalData.OpsiC,
					OpsiD:    soalData.OpsiD,
					Jawaban:  soalData.Jawaban,
				}

				if err := tx.Create(&newSoal).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Sprintf("Error creating soal: %s", err.Error())
				}
				totalSoalCreated++
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Sprintf("Error committing transaction: %s", err.Error())
	}

	response.Data.ModulesCreated = totalModulesCreated
	response.Data.SoalCreated = totalSoalCreated
	return response, "Data imported successfully"
}
