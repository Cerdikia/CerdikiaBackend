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
	protected.GET("/genericMapels/:kelas", controllers.CGenericMapels)
	protected.GET("/genericModules/:kelas/:mapel", controllers.CGenericModules)

	// ========= Actor ===============
	protected.GET("/getAllUsers", controllers.GetUsers)
	// protected.GET("/siswa", controllers.GetSiswa)
	protected.GET("/getDataActor/:role", controllers.GetDataActor)
	protected.GET("/getDataUser", controllers.GetUser)
	protected.PUT("/editDataUser/:role", controllers.UpdateDataActor)

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
