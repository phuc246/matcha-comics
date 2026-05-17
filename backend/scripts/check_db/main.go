package main

import (
	"fmt"
	"log"
	"matcha-comic-backend/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var total int64
	db.Model(&models.Media{}).Select("SUM(size)").Scan(&total)
	fmt.Printf("Total size in DB: %d bytes\n", total)

	var count int64
	db.Model(&models.Media{}).Count(&count)
	fmt.Printf("Total records in DB: %d\n", count)
}
