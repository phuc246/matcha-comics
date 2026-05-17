package handlers

import (
	"fmt"
	"matcha-comic-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{DB: db, Redis: rdb}
}

// GetStories returns stories with support for filtering, sorting and limiting
func (h *Handler) GetStories(c *gin.Context) {
	var stories []models.Story
	storyType := c.Query("type")
	status := c.Query("status")
	sort := c.Query("sort")
	limitStr := c.Query("limit")
	
	query := h.DB.Preload("Genres")
	if storyType != "" {
		query = query.Where("type = ?", storyType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if sort == "trending" {
		today := time.Now().Format("2006-01-02")
		trendingKey := "trending:" + today
		
		// Get top IDs from Redis
		limit := 5
		if limitStr != "" {
			fmt.Sscanf(limitStr, "%d", &limit)
		}
		
		ids, _ := h.Redis.ZRevRange(c.Request.Context(), trendingKey, 0, int64(limit-1)).Result()
		
		if len(ids) > 0 {
			// Fetch stories from DB by these IDs
			// Note: We need to maintain the order from Redis
			if err := h.DB.Preload("Genres").Where("id IN ?", ids).Find(&stories).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			// Simple manual sort to match Redis order (optional but better)
			idMap := make(map[string]models.Story)
			for _, s := range stories {
				idMap[fmt.Sprintf("%d", s.ID)] = s
			}
			
			sortedStories := make([]models.Story, 0)
			for _, id := range ids {
				if s, ok := idMap[id]; ok {
					sortedStories = append(sortedStories, s)
				}
			}
			c.JSON(http.StatusOK, sortedStories)
			return
		} else {
			// Fallback to views if trending is empty
			sort = "views"
		}
	}

	if sort == "views" {
		query = query.Order("views desc")
	} else {
		query = query.Order("id desc")
	}
	
	if limitStr != "" {
		limit := 0
		fmt.Sscanf(limitStr, "%d", &limit)
		if limit > 0 {
			query = query.Limit(limit)
		}
	}
	
	if err := query.Find(&stories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stories)
}

// GetStoryDetail returns a single story by slug and increments views
func (h *Handler) GetStoryDetail(c *gin.Context) {
	slug := c.Param("slug")
	var story models.Story
	
	if err := h.DB.Preload("Genres").Preload("Chapters").Where("slug = ?", slug).First(&story).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}

	// Self-healing: if latest_chapter is 0 but chapters exist, sync it
	if story.LatestChapter == 0 && len(story.Chapters) > 0 {
		var maxNum float64
		for _, ch := range story.Chapters {
			if ch.Number > maxNum {
				maxNum = ch.Number
			}
		}
		if maxNum > 0 {
			story.LatestChapter = maxNum
			h.DB.Model(&story).Update("latest_chapter", maxNum)
		}
	}
	
	// Increment total views in DB
	h.DB.Model(&story).Update("views", gorm.Expr("views + 1"))
	
	// Track daily trending in Redis
	ctx := c.Request.Context()
	today := time.Now().Format("2006-01-02")
	trendingKey := "trending:" + today
	h.Redis.ZIncrBy(ctx, trendingKey, 1, fmt.Sprintf("%d", story.ID))
	// Set expiration for the daily key to 48 hours to save memory
	h.Redis.Expire(ctx, trendingKey, 48*time.Hour)
	
	c.JSON(http.StatusOK, story)
}

// GetChapter returns a single chapter
func (h *Handler) GetChapter(c *gin.Context) {
	slug := c.Param("slug")
	chapterNum := c.Param("chapter")
	
	var story models.Story
	if err := h.DB.Where("slug = ?", slug).First(&story).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}
	
	var chapter models.Chapter
	if err := h.DB.Preload("Servers").Where("story_id = ? AND number = ?", story.ID, chapterNum).First(&chapter).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}
	
	c.JSON(http.StatusOK, chapter)
}

// CreateStory adds a new story (Admin only)
func (h *Handler) CreateStory(c *gin.Context) {
	var story models.Story
	if err := c.ShouldBindJSON(&story); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.DB.Create(&story).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, story)
}

// UpdateStory updates an existing story
func (h *Handler) UpdateStory(c *gin.Context) {
	id := c.Param("id")
	var story models.Story
	
	if err := h.DB.First(&story, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}
	
	type UpdatePayload struct {
		Title       string         `json:"title"`
		Slug        string         `json:"slug"`
		Status      string         `json:"status"`
		Description string         `json:"description"`
		Author      string         `json:"author"`
		Publisher   string         `json:"publisher"`
		CoverURL    string         `json:"coverUrl"`
		Genres      []models.Genre `json:"genres"`
	}
	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	story.Title = payload.Title
	story.Slug = payload.Slug
	story.Status = payload.Status
	story.Description = payload.Description
	story.Author = payload.Author
	story.Publisher = payload.Publisher
	story.CoverURL = payload.CoverURL
	
	if len(payload.Genres) > 0 {
		h.DB.Model(&story).Association("Genres").Replace(payload.Genres)
	} else {
		h.DB.Model(&story).Association("Genres").Clear()
	}
	
	h.DB.Save(&story)
	
	c.JSON(http.StatusOK, story)
}

// DeleteStory deletes a story
func (h *Handler) DeleteStory(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Story{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Story deleted successfully"})
}
