package repositories

import (
	"coba1BE/config"
	"coba1BE/models/soal"
	"coba1BE/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReciveAndStoreImage(c *gin.Context) (url *soal.UploadImageResponse, mesage string) {
	var result soal.UploadImageResponse
	const maxUploadSize = 2 << 20 // 2MB

	// Batasi ukuran request body
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Gagal upload file:")
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload file: " + err.Error()})
		return nil, fmt.Sprintf("error : %e", err.Error)
	}

	// Tambahan: cek apakah file benar-benar di bawah 2MB (opsional karena sudah dibatasi di atas)
	if file.Size > maxUploadSize {
		fmt.Println("Ukuran file melebihi batas maksimum 2MB")
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file melebihi batas maksimum 2MB"})
		return nil, "error : Ukuran file melebihi batas maksimum 2MB"
	}

	uniqueFileName := services.BuatNamaFileUnik(file)

	// Simpan file
	if err := c.SaveUploadedFile(file, fmt.Sprintf("./uploads/%s", uniqueFileName)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan file"})
		return
	}
	fmt.Println("url : /uploads/" + uniqueFileName)

	result = soal.UploadImageResponse{
		Url: fmt.Sprintf("http://localhost:81/uploads/%s", uniqueFileName),
	}

	return &result, "Success"
}

func ReciveAndStoreSoal(c *gin.Context) (*soal.UploadSoal, string) {
	var req soal.UploadSoal

	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Data tidak valid:")
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return nil, fmt.Sprintf("error Data tidak valid : %e", err)
	}

	// Validasi dasar
	fmt.Println("id module : " + strconv.Itoa(int(req.IDModule)))
	fmt.Println("title : " + req.Soal)
	fmt.Println("soal : " + req.Jenis)
	fmt.Println("jawaban 1: " + req.OpsiA)
	fmt.Println("jawaban 2: " + req.OpsiB)
	fmt.Println("jawaban 3: " + req.OpsiC)
	fmt.Println("jawaban 4: " + req.OpsiD)
	fmt.Println("jawaban : " + req.Jawaban)
	if req.Soal == "" || req.OpsiA == "" || req.Jawaban == "" {
		fmt.Println("Semua field harus diisi (judul, soal, opsi, jawaban)")
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field harus diisi (judul, soal, opsi, jawaban)"})
		return nil, fmt.Sprintf("error Semua field harus diisi")
	}

	// Simpan ke database
	db := config.DB
	if err := db.Create(&req).Error; err != nil {
		fmt.Println("Gagal menyimpan ke database")
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return nil, fmt.Sprintf("gagal menyimpan ke database : %e", err.Error())
	}

	return &req, "Success"
}

func GetSoalById(c *gin.Context, idSoal int) (*soal.GenericSoalModelResponse, string) {
	// var users users.Siswa
	var soal soal.GenericSoalModelResponse
	var query string
	// fmt.Println("loginRequest role : " + role)

	db := config.DB
	// query := `SELECT email, nama, kelas, date_created FROM siswa WHERE email = ?`

	query = `SELECT id_soal, id_module, soal, jenis, opsi_a, opsi_b, opsi_c, opsi_d, jawaban FROM soal WHERE id_soal = ?`

	// tmpResult := db.Raw(query, role, email).Scan(&soal)
	tmpResult := db.Raw(query, idSoal).Scan(&soal)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error query data : %e", tmpResult.Error)

	} else if tmpResult.RowsAffected == 0 {
		return nil, "no data found"
	} else {
		// Tambahkan asal
		// soal.Role = role
		// fmt.Println("email : " + soal[0].Email)
		return &soal, "Data retrieved successfully"
	}
}

func UpdateDataSoal(soal soal.UploadSoal, idSoal int) (*soal.UploadSoal, string) {
	var query string
	// var result users.UserProfile
	db := config.DB
	query = "UPDATE soal SET id_module = ?, soal = ?, jenis = ?, opsi_a = ? , opsi_b = ?, opsi_c = ?, opsi_d = ?, jawaban = ? WHERE id_soal = ?"

	tempResult := db.Exec(query,
		soal.IDModule,
		soal.Soal,
		soal.Jenis,
		soal.OpsiA,
		soal.OpsiB,
		soal.OpsiC,
		soal.OpsiD,
		soal.Jawaban,
		idSoal)

	if tempResult.Error != nil {
		return nil, fmt.Sprintf("error query : %e", tempResult.Error)
	} else {
		rowsAffected := tempResult.RowsAffected
		if rowsAffected == 0 {
			return nil, "Tidak ada data yang di temukan atau data tidak berubah"
		} else {
			return &soal, fmt.Sprintf("Update soal id : %d Success", idSoal)
		}
	}
}

func DeleteDataSoal(idSoal int) (*soal.DeletDataSoalResponse, string) {
	var query string
	// var deleteDataSoalResponse soal.DeletDataSoalResponse
	// var result users.UserProfile
	db := config.DB
	query = "DELETE FROM soal WHERE id_soal = ?;"

	tempResult := db.Exec(query, idSoal)

	if tempResult.Error != nil {
		return nil, fmt.Sprintf("error query : %e", tempResult.Error)
	} else {
		rowsAffected := tempResult.RowsAffected
		if rowsAffected == 0 {
			return nil, "Tidak ada data yang di temukan atau data tidak berubah"
		} else {
			return &soal.DeletDataSoalResponse{
				Message: fmt.Sprintf("Soal With id : %d Successfully", idSoal),
			}, "Delete soal success"
		}
	}
}
