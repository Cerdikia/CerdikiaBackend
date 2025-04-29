package controllers

import (
	"coba1BE/config"
	"coba1BE/models/ranking"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetModulesWithCompletion(c *gin.Context) {
	db := config.DB
	// Ambil email dari query parameter atau context
	email := c.Query("email")      // Misal ?email=user@example.com
	idKelas := c.Query("id_kelas") // Misal ?email=user@example.com
	idMapel := c.Query("id_mapel") // Misal ?email=user@example.com

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	var modules []ranking.ModuleResponse

	query := `
			SELECT 
    m.id_module,
    m.id_mapel,
    m.id_kelas,
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

	if err := db.Raw(query, email, idKelas, idMapel).Scan(&modules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": modules})
}

func GetRankingByKelas(c *gin.Context) {
	var query string
	var rankings []ranking.RankingResponse

	db := config.DB
	// Ambil ID kelas dari query param
	idKelasStr := c.Query("id_kelas")

	if idKelasStr == "" {
		query := `
		SELECT 
				RANK() OVER (ORDER BY up.exp DESC) AS ranking,
				s.nama,
				k.kelas,
				up.exp
		FROM siswa s
		JOIN kelas k ON s.id_kelas = k.id_kelas
		JOIN user_points up ON s.email = up.email
`

		if err := db.Raw(query).Scan(&rankings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		idKelas, err := strconv.Atoi(idKelasStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_kelas"})
			return
		}

		query = `
			SELECT 
					RANK() OVER (ORDER BY up.exp DESC) AS ranking,
					s.nama,
					k.kelas,
					up.exp
			FROM siswa s
			JOIN kelas k ON s.id_kelas = k.id_kelas
			JOIN user_points up ON s.email = up.email
			WHERE s.id_kelas = ?
	`
		if err := db.Raw(query, idKelas).Scan(&rankings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": rankings})
}
