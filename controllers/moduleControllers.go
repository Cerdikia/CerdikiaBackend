package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	genericactivities "coba1BE/models/genericActivities"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /modules
func GetAllModules(c *gin.Context) {
	var modules []genericactivities.Module
	if err := config.DB.Find(&modules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, modules)
}

// GET /modules/:id
func GetModuleByID(c *gin.Context) {
	id := c.Param("id")
	var module genericactivities.Module
	if err := config.DB.First(&module, "id_module = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}
	c.JSON(http.StatusOK, module)
}

// POST /modules
func CreateModule(c *gin.Context) {
	var input genericactivities.Module
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
		Message: "Data Created Successfully",
		Data:    input,
	})
}

// PUT /modules/:id
func UpdateModule(c *gin.Context) {
	id := c.Param("id")
	var module genericactivities.Module
	if err := config.DB.First(&module, "id_module = ?", id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Module not found",
			Data:    nil,
		})
		return
	}

	var input genericactivities.Module
	if err := c.ShouldBindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := config.DB.Model(&module).Updates(input).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Module updated successfully"})
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data Updated Successfuly",
		Data:    input,
	})
}

// DELETE /modules/:id
func DeleteModule(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&genericactivities.Module{}, "id_module = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Module deleted successfully"})
}

func ToggleModuleReady(c *gin.Context) {
	db := config.DB
	id := c.Param("id_module") // id_module dari URL

	var module genericactivities.Module
	if err := db.First(&module, "id_module = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	// Toggle is_ready
	module.IsReady = !module.IsReady

	if err := db.Save(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Module readiness toggled",
		"is_ready":     module.IsReady,
		"id_module":    module.IDModule,
		"modeule_name": module.ModuleJudul,
	})
}
