package controllers

import (
	"coba1BE/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TestLogsRetrieval is a simple endpoint to test logs retrieval with date filtering
func TestLogsRetrieval(c *gin.Context) {
	db := config.DB

	// Validasi parameter input
	var input struct {
		StartDate string `json:"start_date" form:"start_date"`
		EndDate   string `json:"end_date" form:"end_date"`
	}

	// Coba binding dari query parameters terlebih dahulu
	if err := c.ShouldBindQuery(&input); err != nil {
		// Jika gagal, coba binding dari JSON body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter tidak valid",
				"error":   err.Error(),
			})
			return
		}
	}

	// Get all logs first
	var allLogsCount int64
	db.Table("logs").Count(&allLogsCount)

	// Build query with date filtering
	query := db.Table("logs")

	if input.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal mulai tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		// Set time to beginning of day (00:00:00)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		query = query.Where("created_at >= ?", startDate)
	}

	if input.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format tanggal akhir tidak valid",
				"format":  "YYYY-MM-DD",
			})
			return
		}
		// Set time to end of day (23:59:59)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		query = query.Where("created_at <= ?", endDate)
	}

	// Get filtered logs count
	var filteredLogsCount int64
	query.Count(&filteredLogsCount)

	// Get logs data
	type LogData struct {
		Email     string    `json:"email" gorm:"column:email"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	}

	var logs []LogData
	query.Find(&logs)

	// Format dates for display
	formattedLogs := make([]map[string]interface{}, 0, len(logs))
	for _, log := range logs {
		formattedLogs = append(formattedLogs, map[string]interface{}{
			"email":      log.Email,
			"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message":             "Test logs retrieval",
		"total_logs":          allLogsCount,
		"filtered_logs":       filteredLogsCount,
		"start_date":          input.StartDate,
		"end_date":            input.EndDate,
		"current_server_time": time.Now().Format("2006-01-02 15:04:05"),
		"logs":                formattedLogs,
	})
}
