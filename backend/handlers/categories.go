package handlers

import (
	"matcha-comic-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories returns categories (either hierarchical tree or flat list)
func (h *Handler) GetCategories(c *gin.Context) {
	tree := c.Query("tree")
	var categories []models.Category

	if tree == "true" {
		if err := h.DB.Preload("Children").Where("parent_id IS NULL").Order("categories.order ASC, categories.id ASC").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.DB.Order("categories.order ASC, categories.id ASC").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory adds a new category (Admin only)
func (h *Handler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gorm handles nil ParentID perfectly
	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates an existing category (Admin only)
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Create temporary structure to read updated fields, specifically supporting nullable fields like ParentID
	type UpdatePayload struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Path     string `json:"path"`
		Icon     string `json:"icon"`
		ParentID *uint  `json:"parentId"`
		Order    int    `json:"order"`
	}

	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = payload.Name
	category.Slug = payload.Slug
	category.Path = payload.Path
	category.Icon = payload.Icon
	category.ParentID = payload.ParentID
	category.Order = payload.Order

	if err := h.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category (Admin only)
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Set children's parent_id to nil before deleting parent to avoid orphaned references
	h.DB.Model(&models.Category{}).Where("parent_id = ?", id).Update("parent_id", nil)

	if err := h.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
