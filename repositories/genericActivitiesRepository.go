package repositories

import (
	"coba1BE/config"
	genericactivities "coba1BE/models/genericActivities"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetGenericActivities(c *gin.Context, kelas string) ([]genericactivities.GenericActivitiesResponse, string) {
	var genericActivities []genericactivities.GenericActivitiesResponse
	// var result models.BaseResponseModel

	db := config.DB
	// 	query := `SELECT
	//     m.id_mapel,
	//     mp.mapel AS nama_mapel,
	//     COUNT(*) AS jumlah_modul
	// FROM
	//     modules m
	// JOIN
	//     mapel mp ON m.id_mapel = mp.id_mapel
	// GROUP BY
	//     m.id_mapel, mp.mapel
	// ORDER BY
	//     m.id_mapel;`

	query := `SELECT 
    m.id_mapel,
    mp.mapel AS nama_mapel,
    m.kelas,
    COUNT(*) AS jumlah_modul
FROM 
    modules m
JOIN 
    mapel mp ON m.id_mapel = mp.id_mapel
WHERE 
    m.kelas = ?
GROUP BY 
    m.id_mapel, mp.mapel, m.kelas
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

func GetGenericModules(c *gin.Context, kelas, mapel string) ([]genericactivities.GenericModulesResponse, string) {
	var genericActivities []genericactivities.GenericModulesResponse
	// var result models.BaseResponseModel

	db := config.DB
	query := `SELECT
id_module, module, module_judul, module_deskripsi
FROM modules
WHERE kelas = ? AND id_mapel = ?
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
