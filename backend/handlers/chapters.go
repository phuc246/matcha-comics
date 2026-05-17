package handlers

import (
	"matcha-comic-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"time"
)

// CreateChapter adds a new chapter to a story (Admin only)
func (h *Handler) CreateChapter(c *gin.Context) {
	type ChapterPayload struct {
		StoryID uint    `json:"story_id"`
		Number  float64 `json:"number"`
		Title   string  `json:"title"`
		Note    string  `json:"note"`
		Servers []struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"servers"`
	}
	var payload ChapterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter := models.Chapter{
		StoryID: payload.StoryID,
		Number:  payload.Number,
		Title:   payload.Title,
		Note:    payload.Note,
	}

	if err := h.DB.Create(&chapter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save servers
	for _, s := range payload.Servers {
		server := models.ChapterServer{
			ChapterID: chapter.ID,
			Name:      s.Name,
			Type:      s.Type,
			Content:   s.Content,
		}
		h.DB.Create(&server)
	}

	// Update Story's LatestChapter and UpdatedAt
	h.DB.Model(&models.Story{ID: payload.StoryID}).Updates(map[string]interface{}{
		"latest_chapter": payload.Number,
		"updated_at":     time.Now(),
	})

	// Re-load with servers
	h.DB.Preload("Servers").First(&chapter, chapter.ID)
	c.JSON(http.StatusCreated, chapter)
}

// GetChaptersByStory returns all chapters for a given story ID (Admin only)
func (h *Handler) GetChaptersByStory(c *gin.Context) {
	storyID := c.Param("id")
	var chapters []models.Chapter
	if err := h.DB.Preload("Servers").Where("story_id = ?", storyID).Order("number desc").Find(&chapters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chapters)
}

// GetChapterDetail returns a single chapter by ID (Admin only)
func (h *Handler) GetChapterDetail(c *gin.Context) {
	id := c.Param("id")
	var chapter models.Chapter
	if err := h.DB.Preload("Servers").First(&chapter, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}
	c.JSON(http.StatusOK, chapter)
}

// UpdateChapter updates an existing chapter (Admin only)
func (h *Handler) UpdateChapter(c *gin.Context) {
	id := c.Param("id")
	var chapter models.Chapter
	if err := h.DB.First(&chapter, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	type UpdatePayload struct {
		Number  float64 `json:"number"`
		Title   string  `json:"title"`
		Note    string  `json:"note"`
		Servers []struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"servers"`
	}
	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter.Number = payload.Number
	chapter.Title = payload.Title
	chapter.Note = payload.Note
	h.DB.Save(&chapter)

	// Replace servers
	if len(payload.Servers) > 0 {
		h.DB.Where("chapter_id = ?", chapter.ID).Delete(&models.ChapterServer{})
		for _, s := range payload.Servers {
			server := models.ChapterServer{ChapterID: chapter.ID, Name: s.Name, Type: s.Type, Content: s.Content}
			h.DB.Create(&server)
		}
	}

	// Update Story's LatestChapter if this was the latest one
	var story models.Story
	h.DB.First(&story, chapter.StoryID)
	if payload.Number >= story.LatestChapter {
		h.DB.Model(&story).Updates(map[string]interface{}{
			"latest_chapter": payload.Number,
			"updated_at":     time.Now(),
		})
	}

	h.DB.Preload("Servers").First(&chapter, chapter.ID)
	c.JSON(http.StatusOK, chapter)
}

// DeleteChapter deletes a chapter (Admin only)
func (h *Handler) DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	// Delete related servers first
	h.DB.Where("chapter_id = ?", id).Delete(&models.ChapterServer{})
	// Recalculate LatestChapter for story
	var chapter models.Chapter
	h.DB.First(&chapter, id)
	storyID := chapter.StoryID

	if err := h.DB.Delete(&models.Chapter{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update story's latest chapter after deletion
	var lastChapter models.Chapter
	if err := h.DB.Where("story_id = ?", storyID).Order("number desc").First(&lastChapter).Error; err == nil {
		h.DB.Model(&models.Story{}).Where("id = ?", storyID).Update("latest_chapter", lastChapter.Number)
	} else {
		h.DB.Model(&models.Story{}).Where("id = ?", storyID).Update("latest_chapter", 0)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter deleted successfully"})
}

