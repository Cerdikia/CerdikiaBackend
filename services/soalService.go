package services

import (
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

func BuatNamaFileUnik(file *multipart.FileHeader) string {
	// Dapatkan ekstensi dan buat nama file unik
	ext := filepath.Ext(file.Filename)
	uniqueFileName := uuid.New().String() + ext
	return uniqueFileName
}
