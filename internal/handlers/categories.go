package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"spending-tracker/db"
	"spending-tracker/templates/components"
)

// GetCategories returns the list of all categories
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

// GetCategoryOptions returns category options for a dropdown, with optional selected value
func (h *Handler) GetCategoryOptions(c *gin.Context) {
	selected, _ := strconv.ParseInt(c.Query("selected"), 10, 64)

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	var selectedPtr *int64
	if selected > 0 {
		selectedPtr = &selected
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryOptions(categories, selectedPtr).Render(c.Request.Context(), c.Writer)
}

// CreateCategory creates a new category
func (h *Handler) CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	color := c.PostForm("color")

	_, err := db.CreateCategory(c.Request.Context(), name, color)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

// UpdateCategory updates an existing category
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	name := c.PostForm("name")
	color := c.PostForm("color")

	if err := db.UpdateCategory(c.Request.Context(), id, name, color); err != nil {
		c.String(http.StatusInternalServerError, "Error updating category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

// EditCategoryName returns the inline edit form for a category name
func (h *Handler) EditCategoryName(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	cat, err := db.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, "Category not found")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryNameEdit(*cat).Render(c.Request.Context(), c.Writer)
}

// DeleteCategory deletes a category
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := db.DeleteCategory(c.Request.Context(), id); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

// CategoryModal returns the manage categories modal form
func (h *Handler) CategoryModal(c *gin.Context) {
	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryModal(categories).Render(c.Request.Context(), c.Writer)
}
