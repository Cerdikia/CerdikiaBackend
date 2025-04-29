package controllers

import (
	"coba1BE/models"
	"coba1BE/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CGenericMapels(c *gin.Context) {
	// var response models.BaseResponseModel

	kelas := c.Query("id_kelas")
	strIsReady := c.DefaultQuery("finished", "0")
	isReady, err := strconv.ParseBool(strIsReady)

	if err != nil {
		// Jika query tidak valid (misal: bukan "true" atau "false")
		// c.JSON(400, gin.H{"error": "Parameter 'finished' harus bernilai true atau false"})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter 'finished' harus bernilai true atau false",
			Data:    nil,
		})
		return
	}

	// chek apakah parameter di isi
	if kelas == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "no parmeter found",
			Data:    nil,
		})
		return
	}

	result, msg := repositories.GetGenericActivities(c, kelas, isReady)

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
	email := c.Query("email")
	idMapel := c.Query("id_mapel")
	idKelas := c.Query("id_kelas")
	strFinished := c.DefaultQuery("finished", "1")

	finished, err := strconv.ParseBool(strFinished)

	if err != nil {
		// Jika query tidak valid (misal: bukan "true" atau "false")
		// c.JSON(400, gin.H{"error": "Parameter 'finished' harus bernilai true atau false"})
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Parameter 'finished' harus bernilai true atau false",
			Data:    nil,
		})
		return
	}

	if email != "" && idKelas != "" && idMapel != "" {
		result, msg := repositories.SpesifiedModulesKelas(c, email, idKelas, idMapel, finished)
		if strings.Contains(msg, "Success") {
			c.JSON(http.StatusOK, models.BaseResponseModel{
				Message: msg,
				Data:    result,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: msg,
				Data:    nil,
			})
			return
		}
	} else if idKelas != "" && idMapel != "" {
		result, msg := repositories.GetGenericModulesKelas(c, idKelas, idMapel, finished)
		if strings.Contains(msg, "Success") {
			c.JSON(http.StatusOK, models.BaseResponseModel{
				Message: msg,
				Data:    result,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: msg,
				Data:    nil,
			})
			return
		}
	} else if idKelas != "" {
		result, msg := repositories.GetGenericModulesByKelas(c, idKelas, finished)
		if strings.Contains(msg, "Success") {
			c.JSON(http.StatusOK, models.BaseResponseModel{
				Message: msg,
				Data:    result,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: msg,
				Data:    nil,
			})
			return
		}
	} else if idMapel != "" {
		result, msg := repositories.GetGenericModules(c, idMapel, finished)
		if strings.Contains(msg, "Success") {
			c.JSON(http.StatusOK, models.BaseResponseModel{
				Message: msg,
				Data:    result,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, models.BaseResponseModel{
				Message: msg,
				Data:    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Unidetifided Query Params",
			Data:    nil,
		})
		return
	}
}

func CGenericModule(c *gin.Context) {
	// var response models.BaseResponseModel

	idModule := c.Param("id_module")

	result, msg := repositories.GetGenericModule(c, idModule)

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
