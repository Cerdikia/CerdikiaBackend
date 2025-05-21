package controllers

import (
	"coba1BE/models/soal"
	"coba1BE/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SaveImportedDataHandler handles the API endpoint for saving imported data
func SaveImportedDataHandler(c *gin.Context) {
	// Parse request body
	var importedData []soal.ImportedData
	if err := c.ShouldBindJSON(&importedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate imported data
	if len(importedData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No data to import",
		})
		return
	}

	// Process and save the imported data
	result, message := repositories.SaveImportedData(importedData)
	if result == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, result)
}
