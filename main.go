package main

import (
	"coba1BE/config"
	"coba1BE/controllers"
	"coba1BE/middleware"

	// "fmt"

	"log"

	"github.com/gin-gonic/gin"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
)

func main() {
	r := gin.Default()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r.GET("/users", controllers.GetUsers)
	r.GET("/user/:email", controllers.GetUser)
	r.POST("/users", controllers.CreateUser)

	r.POST("/login", controllers.Login)

	// protected := r.Group("/api")
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", controllers.GetUsers)

	// Menutup koneksi database saat aplikasi berhenti
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close() // Pastikan koneksi ditutup saat aplikasi keluar
	}

	if err := r.Run(":4321"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
