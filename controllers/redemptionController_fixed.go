package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"coba1BE/config"
	"coba1BE/models/redemption"
)

// generateRedemptionCode generates a unique 12-character alphanumeric code
func generateRedemptionCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 12)

	for i := range code {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[num.Int64()]
	}

	return string(code)
}

// RedeemGifts handles the redemption of gifts by users
func RedeemGifts(c *gin.Context) {
	db := config.DB

	// Parse request body
	var input redemption.RedemptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate input
	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is required",
		})
		return
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "At least one item is required",
		})
		return
	}

	// Check if user is verified
	var userVerified redemption.UserVerified
	result := db.Where("email = ?", input.Email).First(&userVerified)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found or not registered for verification",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error checking user verification status",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Check if user is verified with status 'accept'
	if userVerified.VerifiedStatus != "accept" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "User is not verified. Current status: " + userVerified.VerifiedStatus,
		})
		return
	}

	// Get user's current diamond balance
	var userPoints redemption.UserPoints
	if err := db.Where("email = ?", input.Email).First(&userPoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving user points",
			"error":   err.Error(),
		})
		return
	}

	// Calculate total diamond cost and check item availability
	type ItemDetail struct {
		Barang       redemption.BarangInfo
		RequestedQty int
		DiamondCost  int64
	}

	itemDetails := make([]ItemDetail, 0, len(input.Items))
	totalDiamondCost := int64(0)

	for _, item := range input.Items {
		var barang redemption.BarangInfo
		if err := db.Where("id_barang = ?", item.IDBarang).First(&barang).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Item with ID %d not found", item.IDBarang),
			})
			return
		}

		// Check if enough stock is available
		if barang.Jumlah < int64(item.Jumlah) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Not enough stock for item %s. Requested: %d, Available: %d",
					barang.NamaBarang, item.Jumlah, barang.Jumlah),
			})
			return
		}

		// Calculate diamond cost for this item
		itemDiamondCost := barang.Diamond * int64(item.Jumlah)
		totalDiamondCost += itemDiamondCost

		// Add to item details for processing later
		itemDetails = append(itemDetails, ItemDetail{
			Barang:       barang,
			RequestedQty: item.Jumlah,
			DiamondCost:  itemDiamondCost,
		})
	}

	// Check if user has enough diamonds
	if int64(userPoints.Diamond) < totalDiamondCost {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Not enough diamonds. Required: %d, Available: %d",
				totalDiamondCost, userPoints.Diamond),
		})
		return
	}

	// Begin transaction
	tx := db.Begin()

	// Update user's diamond balance
	if err := tx.Model(&redemption.UserPoints{}).Where("email = ?", input.Email).Update(
		"diamond", gorm.Expr("diamond - ?", totalDiamondCost)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating user's diamond balance",
			"error":   err.Error(),
		})
		return
	}

	// Process each item
	successfulItems := make([]map[string]interface{}, 0, len(itemDetails))

	for _, itemDetail := range itemDetails {
		// Update item stock
		if err := tx.Model(&redemption.BarangInfo{}).Where("id_barang = ?", itemDetail.Barang.IDBarang).Update(
			"jumlah", gorm.Expr("jumlah - ?", itemDetail.RequestedQty)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Error updating stock for item %s", itemDetail.Barang.NamaBarang),
				"error":   err.Error(),
			})
			return
		}

		// Generate a unique redemption code
		redemptionCode := generateRedemptionCode()

		// Check if code already exists (very unlikely but good practice)
		var existingLog redemption.LogsPenukaranPoint
		for {
			result := tx.Where("kode_penukaran = ?", redemptionCode).First(&existingLog)
			if result.Error == gorm.ErrRecordNotFound {
				// Code is unique, we can use it
				break
			}
			// Generate a new code if this one already exists
			redemptionCode = generateRedemptionCode()
		}

		// Create redemption log with the unique code
		log := redemption.LogsPenukaranPoint{
			IDBarang:         itemDetail.Barang.IDBarang,
			Email:            input.Email,
			Jumlah:           itemDetail.RequestedQty,
			KodePenukaran:    redemptionCode,
			TanggalPenukaran: time.Now(),
			StatusPenukaran:  "menunggu", // Default status
		}

		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating redemption log",
				"error":   err.Error(),
			})
			return
		}

		// Add to successful items
		successfulItems = append(successfulItems, map[string]interface{}{
			"id_barang":      itemDetail.Barang.IDBarang,
			"nama_barang":    itemDetail.Barang.NamaBarang,
			"jumlah":         itemDetail.RequestedQty,
			"diamond_cost":   itemDetail.DiamondCost,
			"kode_penukaran": redemptionCode,
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error completing transaction",
			"error":   err.Error(),
		})
		return
	}

	// Get updated diamond balance
	var updatedPoints redemption.UserPoints
	if err := db.Where("email = ?", input.Email).First(&updatedPoints).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":            "Redemption successful",
			"items":              successfulItems,
			"total_diamond_cost": totalDiamondCost,
			"diamond_balance":    "Unknown (error retrieving updated balance)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Redemption successful",
		"items":              successfulItems,
		"total_diamond_cost": totalDiamondCost,
		"diamond_balance":    updatedPoints.Diamond,
	})
}

// GetAllRedemptions retrieves all redemption logs with optional filtering
func GetAllRedemptions(c *gin.Context) {
	db := config.DB

	// Optional filter parameters
	email := c.Query("email")
	status := c.Query("status")
	idBarang := c.Query("id_barang")
	kodePenukaran := c.Query("kode_penukaran")

	// Build the query
	query := db.Table("logs_penukaran_point lp")

	// Apply filters if provided
	if email != "" {
		query = query.Where("lp.email = ?", email)
	}

	if status != "" {
		query = query.Where("lp.status_penukaran = ?", status)
	}

	if idBarang != "" {
		query = query.Where("lp.id_barang = ?", idBarang)
	}

	if kodePenukaran != "" {
		query = query.Where("lp.kode_penukaran = ?", kodePenukaran)
	}

	// Order by newest first
	query = query.Order("lp.tanggal_penukaran DESC")

	// Get the logs
	var logs []redemption.LogsPenukaranPoint
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving redemption logs",
			"error":   err.Error(),
		})
		return
	}

	// If no logs found, return empty array
	if len(logs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No redemption logs found",
			"count":   0,
			"data":    []interface{}{},
		})
		return
	}

	// Prepare response data with additional details
	responseData := make([]map[string]interface{}, 0, len(logs))

	for _, log := range logs {
		// Get item details
		var barang redemption.BarangInfo
		db.Where("id_barang = ?", log.IDBarang).First(&barang)

		// Get student name
		var studentName string
		var student struct {
			Nama string `gorm:"column:nama"`
		}
		if err := db.Table("siswa").Select("nama").Where("email = ?", log.Email).First(&student).Error; err == nil {
			studentName = student.Nama
		}

		// Add to response data
		responseData = append(responseData, map[string]interface{}{
			"id_log":            log.IDLog,
			"id_barang":         log.IDBarang,
			"email":             log.Email,
			"jumlah":            log.Jumlah,
			"tanggal_penukaran": log.TanggalPenukaran,
			"kode_penukaran":    log.KodePenukaran,
			"status_penukaran":  log.StatusPenukaran,
			"nama_barang":       barang.NamaBarang,
			"img":               barang.Img,
			"description":       barang.Description,
			"diamond":           barang.Diamond,
			"nama_siswa":        studentName,
		})
	}

	// Return the results
	c.JSON(http.StatusOK, gin.H{
		"message": "Redemption logs retrieved successfully",
		"count":   len(responseData),
		"data":    responseData,
	})
}

// GetRedemptionByID retrieves a specific redemption log by ID
func GetRedemptionByID(c *gin.Context) {
	db := config.DB
	idStr := c.Param("id")

	// First retrieve the basic redemption log
	var log redemption.LogsPenukaranPoint
	if err := db.Where("id_log = ?", idStr).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption log with ID %s not found", idStr),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption log",
				"error":   err.Error(),
			})
		}
		return
	}

	// Get item details
	var barang redemption.BarangInfo
	db.Where("id_barang = ?", log.IDBarang).First(&barang)

	// Get student name
	var studentName string
	var student struct {
		Nama string `gorm:"column:nama"`
	}
	if err := db.Table("siswa").Select("nama").Where("email = ?", log.Email).First(&student).Error; err == nil {
		studentName = student.Nama
	}

	// Prepare response
	response := map[string]interface{}{
		"id_log":            log.IDLog,
		"id_barang":         log.IDBarang,
		"email":             log.Email,
		"jumlah":            log.Jumlah,
		"tanggal_penukaran": log.TanggalPenukaran,
		"kode_penukaran":    log.KodePenukaran,
		"status_penukaran":  log.StatusPenukaran,
		"nama_barang":       barang.NamaBarang,
		"img":               barang.Img,
		"description":       barang.Description,
		"diamond":           barang.Diamond,
		"nama_siswa":        studentName,
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"message": "Redemption log retrieved successfully",
		"data":    response,
	})
}

// GetRedemptionByCode retrieves a specific redemption log by its unique code
func GetRedemptionByCode(c *gin.Context) {
	db := config.DB
	code := c.Param("code")

	// Retrieve the redemption log by code
	var log redemption.LogsPenukaranPoint
	if err := db.Where("kode_penukaran = ?", code).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption log with code %s not found", code),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption log",
				"error":   err.Error(),
			})
		}
		return
	}

	// Get item details
	var barang redemption.BarangInfo
	db.Where("id_barang = ?", log.IDBarang).First(&barang)

	// Get student name
	var studentName string
	var student struct {
		Nama string `gorm:"column:nama"`
	}
	if err := db.Table("siswa").Select("nama").Where("email = ?", log.Email).First(&student).Error; err == nil {
		studentName = student.Nama
	}

	// Prepare response
	response := map[string]interface{}{
		"id_log":            log.IDLog,
		"id_barang":         log.IDBarang,
		"email":             log.Email,
		"jumlah":            log.Jumlah,
		"tanggal_penukaran": log.TanggalPenukaran,
		"kode_penukaran":    log.KodePenukaran,
		"status_penukaran":  log.StatusPenukaran,
		"nama_barang":       barang.NamaBarang,
		"img":               barang.Img,
		"description":       barang.Description,
		"diamond":           barang.Diamond,
		"nama_siswa":        studentName,
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"message": "Redemption log retrieved successfully",
		"data":    response,
	})
}

// UpdateRedemptionStatus updates the status of a redemption log
func UpdateRedemptionStatus(c *gin.Context) {
	db := config.DB
	idStr := c.Param("id")

	// Parse request body
	var input struct {
		Status string `json:"status" binding:"required"` // menunggu, selesai, dibatalkan
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate status value
	if input.Status != "menunggu" && input.Status != "selesai" && input.Status != "dibatalkan" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid status value. Must be one of: menunggu, selesai, dibatalkan",
		})
		return
	}

	// Check if the redemption log exists
	var log redemption.LogsPenukaranPoint
	if err := db.Where("id_log = ?", idStr).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption log with ID %s not found", idStr),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption log",
				"error":   err.Error(),
			})
		}
		return
	}

	// Update the status
	if err := db.Model(&redemption.LogsPenukaranPoint{}).Where("id_log = ?", idStr).Update("status_penukaran", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating redemption status",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Redemption status updated successfully",
		"id_log":  log.IDLog,
		"status":  input.Status,
	})
}

// DeleteRedemption deletes a redemption log
func DeleteRedemption(c *gin.Context) {
	db := config.DB
	idStr := c.Param("id")

	// Check if the redemption log exists
	var log redemption.LogsPenukaranPoint
	if err := db.Where("id_log = ?", idStr).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption log with ID %s not found", idStr),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption log",
				"error":   err.Error(),
			})
		}
		return
	}

	// Delete the redemption log
	if err := db.Delete(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting redemption log",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Redemption log deleted successfully",
	})
}
