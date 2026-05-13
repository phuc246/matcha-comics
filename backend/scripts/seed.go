package main

import (
	"log"
	"matcha-comic-backend/config"
	"matcha-comic-backend/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load(".env")
	db := config.ConnectDB()

	// Auto Migrate before seeding
	db.AutoMigrate(&models.User{}, &models.Story{}, &models.Chapter{}, &models.Genre{})

	// Create admin user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "superadmin",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("Admin user might already exist:", err)
	} else {
		log.Println("Admin user created: superadmin / password")
	}

	// Create some mock genres
	genres := []models.Genre{
		{Name: "Action", Slug: "action"},
		{Name: "Romance", Slug: "romance"},
		{Name: "Fantasy", Slug: "fantasy"},
	}

	for _, g := range genres {
		db.FirstOrCreate(&g, models.Genre{Slug: g.Slug})
	}
	log.Println("Seeding completed!")
}
