package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	genericactivities "coba1BE/models/genericActivities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /mapel
func GetAllMapel(c *gin.Context) {
	// 	// Ambil query params
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

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

	var mapels []genericactivities.Mapel
	if err := config.DB.Limit(limitInt).
		Offset(offset).
		Order("mapel ASC").Find(&mapels).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// c.JSON(http.StatusOK, mapels)
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "success get data mapel",
		Data:    mapels,
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
