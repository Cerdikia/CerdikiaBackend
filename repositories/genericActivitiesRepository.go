package repositories

import (
	"coba1BE/config"
	genericactivities "coba1BE/models/genericActivities"
	"coba1BE/models/soal"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetGenericActivities(c *gin.Context, kelas string, isReady bool) ([]genericactivities.GenericActivitiesResponse, string) {
	var genericActivities []genericactivities.GenericActivitiesResponse
	// var result models.BaseResponseModel

	db := config.DB

	query := `SELECT 
    m.id_mapel,
    mp.mapel AS nama_mapel,
    kelas.kelas,
    COUNT(*) AS jumlah_modul
FROM 
    modules m
JOIN 
    mapel mp ON m.id_mapel = mp.id_mapel
JOIN 
  kelas ON m.id_kelas = kelas.id_kelas
WHERE 
    m.id_kelas = ? AND is_ready = ?
GROUP BY 
    m.id_mapel, mp.mapel, m.id_kelas
ORDER BY 
    m.id_mapel;`

	tmpResult := db.Raw(query, kelas, isReady).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericActivitiesAllStatusKelas(c *gin.Context, kelas string) ([]genericactivities.GenericActivitiesResponse, string) {
	var genericActivities []genericactivities.GenericActivitiesResponse
	// var result models.BaseResponseModel

	db := config.DB

	query := `SELECT 
    m.id_mapel,
    mp.mapel AS nama_mapel,
    kelas.kelas,
    COUNT(*) AS jumlah_modul
FROM 
    modules m
JOIN 
    mapel mp ON m.id_mapel = mp.id_mapel
JOIN 
  kelas ON m.id_kelas = kelas.id_kelas
WHERE 
    m.id_kelas = ?
GROUP BY 
    m.id_mapel, mp.mapel, m.id_kelas
ORDER BY 
    m.id_mapel;`

	tmpResult := db.Raw(query, kelas).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericActivitiesAllStatus(c *gin.Context) ([]genericactivities.GenericActivitiesResponse, string) {
	var genericActivities []genericactivities.GenericActivitiesResponse
	// var result models.BaseResponseModel

	db := config.DB

	query := `SELECT 
    m.id_mapel,
    mp.mapel AS nama_mapel,
    kelas.kelas,
    COUNT(*) AS jumlah_modul
FROM 
    modules m
JOIN 
    mapel mp ON m.id_mapel = mp.id_mapel
JOIN 
  kelas ON m.id_kelas = kelas.id_kelas
GROUP BY 
    m.id_mapel, mp.mapel, m.id_kelas
ORDER BY 
    m.id_mapel;`

	tmpResult := db.Raw(query).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func SpesifiedModulesKelas(c *gin.Context, email, kelas, mapel string, isFinished bool) ([]genericactivities.SpesifiedModulesKelasResponse, string) {
	var genericActivities []genericactivities.SpesifiedModulesKelasResponse
	// var result models.BaseResponseModel
	db := config.DB
	query := `
			SELECT 
    m.id_module,
    m.module,
    m.module_judul,
    m.module_deskripsi,
    m.is_ready,
    CASE 
        WHEN COUNT(l.id_logs) > 0 THEN TRUE
        ELSE FALSE
    END AS is_completed
FROM modules m
LEFT JOIN logs l 
    ON m.id_module = l.id_module 
    AND l.email = ?
WHERE m.id_kelas = ? AND m.id_mapel = ? AND is_ready = ? 
GROUP BY m.id_module`

	tmpResult := db.Raw(query, email, kelas, mapel, isFinished).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func SpesifiedModulesKelasAllStatus(c *gin.Context, email, kelas, mapel string) ([]genericactivities.SpesifiedModulesKelasResponse, string) {
	var genericActivities []genericactivities.SpesifiedModulesKelasResponse
	// var result models.BaseResponseModel
	db := config.DB
	query := `
			SELECT 
    m.id_module,
    m.module,
    m.module_judul,
    m.module_deskripsi,
    m.is_ready,
    CASE 
        WHEN COUNT(l.id_logs) > 0 THEN TRUE
        ELSE FALSE
    END AS is_completed
FROM modules m
LEFT JOIN logs l 
    ON m.id_module = l.id_module 
    AND l.email = ?
WHERE m.id_kelas = ? AND m.id_mapel = ?
GROUP BY m.id_module`

	tmpResult := db.Raw(query, email, kelas, mapel).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModulesKelas(c *gin.Context, kelas, mapel string, isFinished bool) ([]genericactivities.GenericModulesKelasResponse, string) {
	var genericActivities []genericactivities.GenericModulesKelasResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
id_module, module, module_judul, module_deskripsi
FROM modules
WHERE id_kelas = ? AND id_mapel = ? AND is_ready = ?
ORDER BY module ASC;`

	tmpResult := db.Raw(query, kelas, mapel, isFinished).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModulesKelasAllStatus(c *gin.Context, kelas, mapel string) ([]genericactivities.GenericModulesKelasResponse, string) {
	var genericActivities []genericactivities.GenericModulesKelasResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
id_module, module, module_judul, module_deskripsi
FROM modules
WHERE id_kelas = ? AND id_mapel = ?
ORDER BY module ASC;`

	tmpResult := db.Raw(query, kelas, mapel).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModules(c *gin.Context, mapel string, finished bool) ([]genericactivities.GenericModulesResponse, string) {
	var genericActivities []genericactivities.GenericModulesResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
kelas.kelas, id_module, module, module_judul, module_deskripsi, is_ready
FROM modules
JOIN
	kelas ON modules.id_kelas = kelas.id_kelas
WHERE id_mapel = ? AND is_ready = ?
ORDER BY kelas ASC;`

	tmpResult := db.Raw(query, mapel, finished).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}
	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModulesAllStatus(c *gin.Context, mapel string) ([]genericactivities.GenericModulesResponse, string) {
	var genericActivities []genericactivities.GenericModulesResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
kelas.kelas, id_module, module, module_judul, module_deskripsi, is_ready
FROM modules
JOIN
	kelas ON modules.id_kelas = kelas.id_kelas
WHERE id_mapel = ?
ORDER BY kelas ASC;`

	tmpResult := db.Raw(query, mapel).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}
	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModulesByKelas(c *gin.Context, idKelas string, finished bool) ([]genericactivities.GenericKelasResponse, string) {
	var genericActivities []genericactivities.GenericKelasResponse

	db := config.DB
	query := `SELECT modules.id_module,
  modules.module,
  modules.module_judul,
  modules.module_deskripsi,
	modules.is_ready,
  mapel.mapel
FROM modules
JOIN kelas ON modules.id_kelas = kelas.id_kelas
JOIN mapel ON modules.id_mapel = mapel.id_mapel
WHERE modules.id_kelas = ? AND modules.is_ready = ?
ORDER BY kelas.kelas ASC;`

	tmpResult := db.Raw(query, idKelas, finished).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}
	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericModulesByKelasAllStatus(c *gin.Context, idKelas string) ([]genericactivities.GenericKelasResponse, string) {
	var genericActivities []genericactivities.GenericKelasResponse

	db := config.DB
	query := `SELECT modules.id_module,
  modules.module,
  modules.module_judul,
  modules.module_deskripsi,
	modules.is_ready,
  mapel.mapel
FROM modules
JOIN kelas ON modules.id_kelas = kelas.id_kelas
JOIN mapel ON modules.id_mapel = mapel.id_mapel
WHERE modules.id_kelas = ?
ORDER BY kelas.kelas ASC;`

	tmpResult := db.Raw(query, idKelas).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}
	return genericActivities, fmt.Sprintf("Success")
}

func GetGenericSoal(c *gin.Context, mapel string) ([]soal.GenericSoalModelResponse, string) {
	var genericActivities []soal.GenericSoalModelResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
id_soal, id_module, soal, jenis, opsi_a, opsi_b, opsi_c, opsi_d, jawaban
FROM soal
WHERE id_module = ?
;`

	tmpResult := db.Raw(query, mapel).Scan(&genericActivities)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error fetching data : %e", tmpResult.Error)
	}

	if tmpResult.RowsAffected == 0 {
		return nil, fmt.Sprintf("no data found, maybe wrong in query")
	}

	return genericActivities, fmt.Sprintf("Success")
}
