package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	// MYSQL DSN
	// dsn := "mrg:123123123@tcp(localhost:3306)/cerdikia?charset=utf8mb4&parseTime=True&loc=Local"

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("dsn : " + dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	DB = db
	// db.AutoMigrate(&ListMhsModel.ListMhs{}, &AbsenMhsModel.ListAbsen{})
	return db, nil
}
