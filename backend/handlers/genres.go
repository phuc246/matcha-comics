package handlers

import (
	"matcha-comic-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetGenres returns all genres
func (h *Handler) GetGenres(c *gin.Context) {
	var genres []models.Genre
	if err := h.DB.Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// CreateGenre adds a new genre (Admin only)
func (h *Handler) CreateGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// UpdateGenre updates an existing genre (Admin only)
func (h *Handler) UpdateGenre(c *gin.Context) {
	id := c.Param("id")
	var genre models.Genre
	if err := h.DB.First(&genre, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Save(&genre)
	c.JSON(http.StatusOK, genre)
}

// DeleteGenre deletes a genre (Admin only)
func (h *Handler) DeleteGenre(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Genre{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted"})
}
