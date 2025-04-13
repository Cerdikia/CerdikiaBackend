package controllers

import (
	"coba1BE/models"
	"coba1BE/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CGenericMapels(c *gin.Context) {
	// var response models.BaseResponseModel

	kelas := c.Param("kelas")

	// chek apakah parameter di isi
	if kelas == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	result, msg := repositories.GetGenericActivities(c, kelas)

	if strings.Contains(msg, "error fetching data") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: msg,
			Data:    nil,
		})
		return
	}

	if strings.Contains(msg, "maybe wrong in query") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: msg,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: msg,
		Data:    result,
	})
}

func CGenericModules(c *gin.Context) {
	// var response models.BaseResponseModel

	kelas := c.Param("kelas")
	mapel := c.Param("mapel")

	// chek apakah parameter di isi
	if kelas == "" || mapel == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	result, msg := repositories.GetGenericModules(c, kelas, mapel)

	if strings.Contains(msg, "error fetching data") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: msg,
			Data:    nil,
		})
		return
	}

	if strings.Contains(msg, "maybe wrong in query") {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: msg,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: msg,
		Data:    result,
	})
}
