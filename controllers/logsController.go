package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/logsprogres"
	"coba1BE/services"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllLogs(c *gin.Context) {
	var logs []logsprogres.Log
	db := config.DB

	if err := db.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrived Successfuly",
		Data:    logs,
	})
	// c.JSON(http.StatusOK, logs)
}

func GetLogByID(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var log logsprogres.Log
	if err := db.First(&log, id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrived Successfuly",
		Data:    log,
	})
	// c.JSON(http.StatusOK, log)
}

func CreateLog(c *gin.Context) {
	db := config.DB

	var log logsprogres.Log
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrived Successfuly",
		Data:    log,
	})
	// c.JSON(http.StatusCreated, log)
}

func UpdateLog(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var log logsprogres.Log
	if err := db.First(&log, id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := c.ShouldBindJSON(&log); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := db.Save(&log).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
		Message: "Data retrived Successfuly",
		Data:    log,
	})
	// c.JSON(http.StatusOK, log)
}

func DeleteLog(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	if err := db.Delete(&logsprogres.Log{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Log deleted"})
}

// ================== Get Logs Using Email With pagination Start =====================

// func GetLogsByEmail(c *gin.Context) {
// 	email := c.Param("email")

// 	// Ambil query params
// 	page := c.DefaultQuery("page", "1")
// 	limit := c.DefaultQuery("limit", "10")

// 	// Konversi ke int
// 	pageInt, err := strconv.Atoi(page)
// 	if err != nil || pageInt < 1 {
// 		pageInt = 1
// 	}

// 	limitInt, err := strconv.Atoi(limit)
// 	if err != nil || limitInt < 1 {
// 		limitInt = 10
// 	}

// 	offset := (pageInt - 1) * limitInt

// 	var logs []models.Log
// 	var total int64

// 	// Hitung total data
// 	models.DB.Model(&models.Log{}).Where("email = ?", email).Count(&total)

// 	// Ambil data dengan limit & offset
// 	err = models.DB.Where("email = ?", email).
// 		Limit(limitInt).
// 		Offset(offset).
// 		Order("created_at DESC").
// 		Find(&logs).Error

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":       logs,
// 		"page":       pageInt,
// 		"limit":      limitInt,
// 		"total_data": total,
// 		"total_page": int((total + int64(limitInt) - 1) / int64(limitInt)),
// 	})
// }

// ================== Get Logs Using Email With pagination End =====================

// GET logs 1 orang (WHERE email = ?)

func GetLogsByEmail(c *gin.Context) {
	var logs []logsprogres.Log
	db := config.DB

	// ================== get email using token start =====================

	// email, errmail := services.GetUserEmailFromToken(c)
	// if errmail != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error di adminContollrs": errmail.Error()})
	// 	return
	// }
	// fmt.Println(errmail)

	// // chek apakah parameter di isi
	// if email == "" {
	// 	c.JSON(http.StatusBadRequest, models.BaseResponseModel{
	// 		Message: "no parmeter found",
	// 		Data:    nil,
	// 	})
	// 	return
	// }

	// // chek apakah parameter berformat email
	// _, err := mail.ParseAddress(email)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, models.BaseResponseModel{
	// 		Message: "bad parameter",
	// 		Data:    nil,
	// 	})
	// 	return
	// }
	// ================== get email using token end =====================

	email := c.Param("email")

	if err := db.Where("email = ?", email).
		// Limit(2).
		// Offset(0).
		// Order("id_logs DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetLogsBydEmailWithToken(c *gin.Context) {
	var logs []logsprogres.Log
	db := config.DB

	// ================== get email using token start =====================

	email, errmail := services.GetUserEmailFromToken(c)
	if errmail != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error di adminContollrs": errmail.Error()})
		return
	}
	fmt.Println(errmail)

	// chek apakah parameter di isi
	if email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	// chek apakah parameter berformat email
	_, err := mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "bad parameter",
			Data:    nil,
		})
		return
	}
	// ================== get email using token end =====================

	// email := c.Param("email")

	if err := db.Where("email = ?", email).
		// Limit(2).
		// Offset(0).
		// Order("id_logs DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(logs) == 0 {
		// c.JSON(http.StatusNotFound, gin.H{"message": "No logs found for this email"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "No logs found for this email",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrived Successfuly",
		Data:    logs,
	})
	// c.JSON(http.StatusOK, logs)
}

func GetLogsBy(c *gin.Context) {
	var logs []logsprogres.Log
	db := config.DB
	email := c.Query("email")
	idKelas := c.Query("id_kelas")
	idMapel := c.Query("id_mapel")

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

	// GET logs mapel 1 kelas (WHERE kelas = ? and id_mepel = ?)
	if idKelas != "" && idMapel != "" {
		if err := db.Where("id_kelas = ? AND id_mapel = ?", idKelas, idMapel).
			Limit(limitInt).
			Offset(offset).
			Order("created_at DESC").
			Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if email != "" && idMapel != "" {
		// GET logs 1 orang (WHERE email = ? and id_mepel = ?)
		if err := db.Where("email = ? AND id_mapel = ?", email, idMapel).
			Limit(limitInt).
			Offset(offset).
			Order("created_at DESC").
			Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if idKelas != "" {
		// GET logs 1 kelas (WHERE kelas = ?)
		if err := db.Where("id_kelas = ?", idKelas).
			Limit(limitInt).
			Offset(offset).
			Order("created_at DESC").
			Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if email != "" {
		// GET logs 1 orang (WHERE email = ?)
		if err := db.Where("email = ?", email).
			Limit(limitInt).
			Offset(offset).
			Order("created_at DESC").
			Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(logs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No logs found for this user and module"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// ============================= GET LOGS PERIODE ====================================

func GetLogsByPeriod(c *gin.Context) {
	var logs []logsprogres.Log
	db := config.DB

	period := c.Query("periode") // contoh: "today", "week", "month", "semester", "year"
	now := time.Now()
	end := now // End time is always the current time
	var start time.Time

	switch period {
	case "today":
		// Today: from current time to 24 hours before
		start = now.Add(-24 * time.Hour)

	case "week":
		// Week: from current time to 7 days before
		start = now.AddDate(0, 0, -7)

	case "month":
		// Month: from current time to 30 days before
		start = now.AddDate(0, 0, -30)

	case "semester":
		// Semester: from current time to 6 months before
		start = now.AddDate(0, -6, 0)

	case "year":
		// Year: from current time to 1 year before
		start = now.AddDate(-1, 0, 0)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Periode tidak valid. Gunakan: today, week, month, semester, year",
		})
		return
	}

	if err := db.Where("created_at >= ? AND created_at < ?", start, end).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: fmt.Sprintf("Data logs berdasarkan periode : %s", period),
		Data:    logs,
	})
}
