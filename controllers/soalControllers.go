package controllers

import (
	"coba1BE/models"
	"coba1BE/models/soal"
	"coba1BE/repositories"
	"fmt"
	"strconv"
	"strings"

	// "coba1BE/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// var response models.BaseResponseModel
	// var imageUrl soal.UploadImageResponse

	imageUrl, message := repositories.ReciveAndStoreImage(c)
	if message == "Success" {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: message,
			Data:    imageUrl,
		})
	} else {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
	}
}

func UploadSoal(c *gin.Context) {

	result, message := repositories.ReciveAndStoreSoal(c)
	if message == "Success" {
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: message,
			Data:    result,
		})
	} else {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
	}
}

func GetDataSoal(c *gin.Context) {
	idSoal := c.Param("id_soal")

	fmt.Printf("id soal : " + idSoal)
	// chek apakah parameter di isi
	if idSoal == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	idSoalInt, err := strconv.Atoi(idSoal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	result, message := repositories.GetSoalById(c, idSoalInt)
	if strings.Contains(strings.ToLower(message), "success") {
		fmt.Println("status success")
		c.JSON(http.StatusOK, models.BaseResponseModel{
			Message: message,
			Data:    result,
		})
	} else {
		fmt.Println("status error")
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
	}
}

func UpdateDataSoal(c *gin.Context) {
	// var result users.UserProfileReq
	var response models.BaseResponseModel
	var soal soal.UploadSoal

	idSoal := c.Param("id_soal")

	// chek apakah parameter di isi
	if idSoal == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	idSoalInt, err := strconv.Atoi(idSoal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&soal); err != nil {
		response = models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, message := repositories.UpdateDataSoal(soal, idSoalInt)
	if strings.Contains(strings.ToLower(message), "success") {
		response = models.BaseResponseModel{
			Message: message,
			Data:    result,
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		response = models.BaseResponseModel{
			Message: message,
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
}

func DeleteSoal(c *gin.Context) {
	// DeleteDataSoal

	// var result users.UserProfileReq
	var response models.BaseResponseModel

	idSoal := c.Param("id_soal")

	// chek apakah parameter di isi
	if idSoal == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	idSoalInt, err := strconv.Atoi(idSoal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	result, message := repositories.DeleteDataSoal(idSoalInt)
	if strings.Contains(strings.ToLower(message), "success") {
		response = models.BaseResponseModel{
			Message: message,
			Data:    result,
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		response = models.BaseResponseModel{
			Message: message,
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
}
