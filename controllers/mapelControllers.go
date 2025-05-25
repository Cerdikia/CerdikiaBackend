package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	genericactivities "coba1BE/models/genericActivities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MapelWithModuleCount represents a mapel with its module count
type MapelWithModuleCount struct {
	IDMapel     int    `json:"id_mapel"`
	Mapel       string `json:"mapel"`
	JumlahModul int    `json:"jumlah_modul"`
}

// GET /mapel
func GetAllMapel(c *gin.Context) {
	// 	// Ambil query params
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "100")

	// 	// Konversi ke int
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	// Step 1: Ambil total jumlah mapel
	var jumlahMapel int64
	if err := config.DB.Table("mapel").Count(&jumlahMapel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Gagal menghitung jumlah mapel",
			Data:    nil,
			Total:   0,
		})
		return
	}

	// Query to get mapels with module count
	var mapelsWithCount []MapelWithModuleCount
	query := `
		SELECT 
			m.id_mapel,
			m.mapel,
			COUNT(md.id_module) AS jumlah_modul
		FROM 
			mapel m
		LEFT JOIN 
			modules md ON m.id_mapel = md.id_mapel
		GROUP BY 
			m.id_mapel, m.mapel
		ORDER BY 
			m.mapel ASC
		LIMIT ? OFFSET ?
	`

	if err := config.DB.Raw(query, limitInt, offset).Scan(&mapelsWithCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "success get data mapel",
		Data:    mapelsWithCount,
		Total:   int(jumlahMapel),
	})
}

// GET /mapel/:id
func GetMapelByID(c *gin.Context) {
	id := c.Param("id")
	var mapel genericactivities.Mapel
	if err := config.DB.First(&mapel, "id_mapel = ?", id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Mapel not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Mapel not found",
			Data:    nil,
		})
		return
	}
	// c.JSON(http.StatusOK, mapel)
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "success get data mapel",
		Data:    mapel,
	})
}

// POST /mapel
func CreateMapel(c *gin.Context) {
	var input genericactivities.Mapel
	if err := c.ShouldBindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// c.JSON(http.StatusCreated, input)
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data Created Successfuly",
		Data:    input,
	})
}

// PUT /mapel/:id
func UpdateMapel(c *gin.Context) {
	id := c.Param("id")
	var mapel genericactivities.Mapel
	if err := config.DB.First(&mapel, "id_mapel = ?", id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Mapel not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var input genericactivities.Mapel
	if err := c.ShouldBindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Update hanya field Mapel
	if err := config.DB.Model(&mapel).Update("mapel", input.Mapel).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Mapel updated successfully"})
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "update Successfuly",
		Data:    &mapel,
	})
}

// DELETE /mapel/:id
func DeleteMapel(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&genericactivities.Mapel{}, "id_mapel = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mapel deleted successfully"})
}
