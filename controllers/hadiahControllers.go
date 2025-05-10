package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/hadiah"
	"coba1BE/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllBarang(c *gin.Context) {
	// var response models.BaseResponseModel
	db := config.DB

	var barang []hadiah.TableBarang
	if err := db.Find(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrieved successfuly",
		Data:    barang,
	})
	// c.JSON(http.StatusOK, barang)
}

func GetBarangByID(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var barang hadiah.TableBarang
	if err := db.First(&barang, "id_barang = ?", id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrieved successfuly",
		Data:    barang,
	})
	// c.JSON(http.StatusOK, barang)
}

func CreateBarang(c *gin.Context) {
	db := config.DB

	var barang hadiah.TableBarang
	if err := c.ShouldBindJSON(&barang); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := db.Create(&barang).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Data retrieved successfuly",
		Data:    &barang,
	})
	// c.JSON(http.StatusCreated, barang)
}

func UpdateBarang(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var barang hadiah.TableBarang
	if err := db.First(&barang, "id_barang = ?", id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := c.ShouldBindJSON(&barang); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	db.Save(&barang)
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "New Data",
		Data:    barang,
	})
	// c.JSON(http.StatusOK, barang)
}

func DeleteBarang(c *gin.Context) {
	var hadiah hadiah.TableBarang
	db := config.DB

	id := c.Param("id")
	if err := db.Delete(&hadiah, "id_barang = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
		// 	Message: err.Error(),
		// 	Data:    nil,
		// })
		return
	}
	// c.JSON(http.StatusOK, models.BaseResponseModel{
	// 	Message: "Data deleted successfully",
	// 	Data:    &hadiah,
	// })
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus"})
}

// CreateGift handles the creation of a gift with image upload using FormData
func CreateGift(c *gin.Context) {
	db := config.DB

	// Bind form data
	var form hadiah.BarangForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid form data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Handle image upload
	imageUrl, message := repositories.ReciveAndStoreImage(c)
	if message != "Success" {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to upload image: " + message,
			Data:    nil,
		})
		return
	}

	// Create the gift record
	barang := hadiah.TableBarang{
		NamaBarang:  form.Name,
		Jumlah:      form.Quantity,
		Diamond:     form.DiamondValue,
		Description: form.Description,
		Img:         imageUrl.Url,
	}

	if err := db.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to create gift: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.BaseResponseModel{
		Message: "Gift created successfully",
		Data:    barang,
	})
}

// GetAllGifts retrieves all gifts from the database
func GetAllGifts(c *gin.Context) {
	db := config.DB

	var gifts []hadiah.TableBarang
	if err := db.Order("created_at DESC").Find(&gifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to retrieve gifts: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Gifts retrieved successfully",
		Data:    gifts,
	})
}

// GetGiftByID retrieves a specific gift by its ID
func GetGiftByID(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	var gift hadiah.TableBarang
	if err := db.First(&gift, "id_barang = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Gift not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Gift retrieved successfully",
		Data:    gift,
	})
}

// DeleteGift deletes a gift by its ID
func DeleteGift(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	// First check if the gift exists
	var gift hadiah.TableBarang
	if err := db.First(&gift, "id_barang = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Gift not found",
			Data:    nil,
		})
		return
	}

	// Delete the gift
	if err := db.Delete(&gift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to delete gift: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Gift deleted successfully",
		Data:    nil,
	})
}

// UpdateGift handles updating a gift with FormData
func UpdateGift(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	// First check if the gift exists
	var existingGift hadiah.TableBarang
	if err := db.First(&existingGift, "id_barang = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Gift not found",
			Data:    nil,
		})
		return
	}

	// Create a form for binding
	var form struct {
		NamaBarang  string `form:"nama_barang"`
		Jumlah      int    `form:"jumlah"`
		Diamond     int    `form:"diamond"`
		Description string `form:"description"`
	}

	// Bind form data
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid form data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Update gift fields if provided
	if form.NamaBarang != "" {
		existingGift.NamaBarang = form.NamaBarang
	}
	if form.Jumlah != 0 {
		existingGift.Jumlah = form.Jumlah
	}
	if form.Diamond != 0 {
		existingGift.Diamond = form.Diamond
	}
	if form.Description != "" {
		existingGift.Description = form.Description
	}

	// Check if a new image was uploaded
	_, fileHeader, err := c.Request.FormFile("image")
	if err == nil && fileHeader != nil {
		// Handle image upload
		imageUrl, message := repositories.ReciveAndStoreImage(c)
		if message != "Success" {
			c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
				Message: "Failed to upload image: " + message,
				Data:    nil,
			})
			return
		}
		existingGift.Img = imageUrl.Url
	}

	// Save the updated gift
	if err := db.Save(&existingGift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to update gift: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Gift updated successfully",
		Data:    existingGift,
	})
}
