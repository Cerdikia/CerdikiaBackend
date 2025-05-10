package controllers

import (
	"coba1BE/models"
	"coba1BE/models/genericActivities"
	"coba1BE/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CGenericMapels(c *gin.Context) {
	var isReady bool
	var err error
	var result []genericActivities.GenericActivitiesResponse
	var msg string
	// var response models.BaseResponseModel

	kelas := c.Query("id_kelas")
	strIsReady := c.Query("finished")

	if strIsReady != "" {
		isReady, err = strconv.ParseBool(strIsReady)
		if err != nil {
			// Jika query tidak valid (misal: bukan "true" atau "false")
			// c.JSON(400, gin.H{"error": "Parameter 'finished' harus bernilai true atau false"})
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: "Parameter 'finished' harus bernilai true atau false",
				Data:    nil,
			})
			return
		}
	}

	// chek apakah parameter di isi
	// if kelas == "" {
	// 	c.JSON(http.StatusBadRequest, models.BaseResponseModel{
	// 		Message: "no parmeter found",
	// 		Data:    nil,
	// 	})
	// 	return
	// }
	if strIsReady != "" {
		result, msg = repositories.GetGenericActivities(c, kelas, isReady)
	} else if kelas != "" {
		result, msg = repositories.GetGenericActivitiesAllStatusKelas(c, kelas)
	} else {
		result, msg = repositories.GetGenericActivitiesAllStatus(c)
	}

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
	email := c.Query("email")
	idMapel := c.Query("id_mapel")
	idKelas := c.Query("id_kelas")
	strFinished := c.Query("finished")

	var finished bool
	var err error

	// Parse finished jika ada
	if strFinished != "" {
		finished, err = strconv.ParseBool(strFinished)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: "Parameter 'finished' harus bernilai true atau false",
				Data:    nil,
			})
			return
		}
	}

	// Helper untuk respon JSON
	sendResponse := func(result interface{}, msg string) {
		status := http.StatusOK
		if !strings.Contains(msg, "Success") {
			status = http.StatusBadRequest
			result = nil
		}
		c.JSON(status, models.BaseResponseModel{
			Message: msg,
			Data:    result,
		})
	}

	// Logic query kombinasi
	switch {
	case email != "" && idKelas != "" && idMapel != "" && strFinished != "":
		result, msg := repositories.SpesifiedModulesKelas(c, email, idKelas, idMapel, finished)
		sendResponse(result, msg)

	case email != "" && idKelas != "" && idMapel != "":
		result, msg := repositories.SpesifiedModulesKelasAllStatus(c, email, idKelas, idMapel)
		sendResponse(result, msg)

	case idKelas != "" && idMapel != "" && strFinished != "":
		result, msg := repositories.GetGenericModulesKelas(c, idKelas, idMapel, finished)
		sendResponse(result, msg)

	case idKelas != "" && idMapel != "":
		result, msg := repositories.GetGenericModulesKelasAllStatus(c, idKelas, idMapel)
		sendResponse(result, msg)

	case idKelas != "" && strFinished != "":
		result, msg := repositories.GetGenericModulesByKelas(c, idKelas, finished)
		sendResponse(result, msg)

	case idKelas != "":
		result, msg := repositories.GetGenericModulesByKelasAllStatus(c, idKelas)
		sendResponse(result, msg)

	case idMapel != "" && strFinished != "":
		result, msg := repositories.GetGenericModules(c, idMapel, finished)
		sendResponse(result, msg)

	case idMapel != "":
		result, msg := repositories.GetGenericModulesAllStatus(c, idMapel)
		sendResponse(result, msg)

	default:
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Unidentified query params",
			Data:    nil,
		})
	}
}

func CGenericSoal(c *gin.Context) {
	// var response models.BaseResponseModel

	idModule := c.Param("id_module")

	result, msg := repositories.GetGenericSoal(c, idModule)

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
