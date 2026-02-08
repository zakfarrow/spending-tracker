package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health returns the health status of the application
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
