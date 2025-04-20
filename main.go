package main

import (
	"coba1BE/config"
	"coba1BE/controllers"
	"coba1BE/middleware"
	"fmt"
	"os"

	"time"

	// "fmt"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("PORT")

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// Atur middleware CORS
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"}, // sesuaikan origin frontend kamu
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r.GET("/try", controllers.GetUsers)
	r.POST("/register/:role", controllers.CreateUser)

	r.POST("/login", controllers.LoginControler)
	r.POST("/refresh", controllers.RefreshToken)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/genericMapels/:kelas", controllers.CGenericMapels)                    // ambil semua mata pelajaran dan jumlah modulnya melalui kelas
	protected.GET("/genericModulesClass/:kelas/:mapel", controllers.CGenericModulesClass) // ambil semua module dengan acuan kelas dan mapel
	protected.GET("/genericModules/:id_mapel", controllers.CGenericModules)               // ambil semua module dengan acuan  mapel
	protected.GET("/genericModule/:id_module", controllers.CGenericModule)                // ambil soal dari sebuah module dengan acuan id_module

	// ========= Actor ===============
	protected.GET("/getAllUsers", controllers.GetUsers)
	// protected.GET("/siswa", controllers.GetSiswa)
	protected.GET("/getDataActor/:role", controllers.GetDataActor)
	protected.GET("/getDataUser", controllers.GetUser)
	protected.PUT("/editDataUser/:role", controllers.UpdateDataActor)

	// ========= Siawa Verified ===============
	protected.GET("/verified", controllers.Beingverified)
	r.GET("/verifiedes", controllers.Beingverifieds)
	r.PATCH("/verifiedes", controllers.UpdateUserVerifiedBatch)

	// ========= Soal ===============
	r.POST("/upload-image", controllers.UploadImage)
	r.POST("/upload-soal", controllers.UploadSoal)
	// if strings.Contains(strings.ToLower(message), "success") {
	r.GET("/getDataSoal/:id_soal", controllers.GetDataSoal)
	r.PUT("/editDataSoal/:id_soal", controllers.UpdateDataSoal)
	r.DELETE("/deleteDataSoal/:id_soal", controllers.DeleteSoal)

	// ========= Point ===============
	protected.GET("/point", controllers.GetPoint)
	protected.PUT("/point", controllers.UpdatePoint)

	// Menutup koneksi database saat aplikasi berhenti
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close() // Pastikan koneksi ditutup saat aplikasi keluar
	}

	// if err := r.Run(":4321"); err != nil {
	if err := r.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
