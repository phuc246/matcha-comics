package main

import (
	"log"
	"matcha-comic-backend/config"
	"matcha-comic-backend/models"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db := config.ConnectDB()

	var stories []models.Story
	db.Find(&stories)

	for _, story := range stories {
		var lastChapter models.Chapter
		err := db.Where("story_id = ?", story.ID).Order("number desc").First(&lastChapter).Error
		if err == nil {
			log.Printf("Updating Story %s with LatestChapter %.1f", story.Title, lastChapter.Number)
			db.Model(&story).Update("latest_chapter", lastChapter.Number)
		} else {
			log.Printf("Story %s has no chapters", story.Title)
		}
	}

	log.Println("Sync completed!")
}
