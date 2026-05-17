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

	var users []models.User
	db.Find(&users)
	fmt.Println("Users in DB:")
	for _, u := range users {
		fmt.Printf("- Username: %s, Role: %s\n", u.Username, u.Role)
	}
}
