package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"coba1BE/controllers"
	"coba1BE/models/kelas"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/kelas", controllers.CreateKelas)
	return r
}

func TestCreateKelas(t *testing.T) {
	// Setup
	r := setupRouter()

	// Data yang akan dikirim
	body := kelas.Kelas{
		Kelas: "XII IPA 1",
	}
	jsonValue, _ := json.Marshal(body)

	// Simulasi request
	req, _ := http.NewRequest("POST", "/kelas", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertion
	assert.Equal(t, http.StatusCreated, w.Code)
	var response kelas.Kelas
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, body.Kelas, response.Kelas)
}
