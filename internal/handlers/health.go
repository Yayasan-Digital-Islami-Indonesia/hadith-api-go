package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/models"
)

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns service status
// @Tags         system
// @Produce      json
// @Success      200 {object} models.HealthResponse
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, models.HealthResponse{Status: "ok"})
}