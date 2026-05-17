package routes

import (
	"matcha-comic-backend/handlers"
	"matcha-comic-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	h := handlers.NewHandler(db, rdb)

	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Public Routes
		api.GET("/stories", h.GetStories)
		api.GET("/stories/:slug", h.GetStoryDetail)
		api.GET("/stories/:slug/chapters/:chapter", h.GetChapter)

		// Public Genre Routes
		api.GET("/genres", h.GetGenres)
		
		api.POST("/auth/login", h.Login)
		
		// Admin Routes (Protected)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			// Stories
			admin.POST("/stories", h.CreateStory)
			admin.PUT("/stories/:id", h.UpdateStory)
			admin.DELETE("/stories/:id", h.DeleteStory)

			// Chapters
			admin.POST("/chapters", h.CreateChapter)
			admin.GET("/stories/:id/chapters", h.GetChaptersByStory)
			admin.GET("/chapters/:id", h.GetChapterDetail)
			admin.PUT("/chapters/:id", h.UpdateChapter)
			admin.DELETE("/chapters/:id", h.DeleteChapter)

			// Genres
			admin.POST("/genres", h.CreateGenre)
			admin.PUT("/genres/:id", h.UpdateGenre)
			admin.DELETE("/genres/:id", h.DeleteGenre)

			// Image Upload & Media Library
			admin.POST("/upload-images", h.UploadImages)
			admin.GET("/media", h.ListMedia)
			admin.DELETE("/media/:id", h.DeleteMedia)
			admin.GET("/storage-stats", h.GetStorageStats)
		}
	}
}
