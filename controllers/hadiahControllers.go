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
