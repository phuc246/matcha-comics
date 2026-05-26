package main

import (
	"log"
	"matcha-comic-backend/config"
	"matcha-comic-backend/models"
	"matcha-comic-backend/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to Database
	db := config.ConnectDB()
	
	// Auto Migrate Models
	db.AutoMigrate(&models.User{}, &models.Genre{}, &models.Category{}, &models.Story{}, &models.Chapter{}, &models.ChapterServer{}, &models.Media{})

	// Connect to Redis
	redisClient := config.ConnectRedis()

	// Initialize Gin
	r := gin.Default()
	r.MaxMultipartMemory = 100 << 20 // 100 MB

	// Serve uploaded images statically
	os.MkdirAll("uploads", 0755)
	r.Static("/uploads", "./uploads")

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://matcha-comics.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register Routes
	routes.RegisterRoutes(r, db, redisClient)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
