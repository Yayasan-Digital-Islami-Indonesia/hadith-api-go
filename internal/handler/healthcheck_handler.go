package handler

import (
	"net/http"

	"hadith-api-go/internal/domain/healthcheck"
	"hadith-api-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	service healthcheck.HealthCheckService
}

func NewHealthCheckHandler(service healthcheck.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{service: service}
}

func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	health, err := h.service.HealthCheck(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	response.Success(c, health)
}

func (h *HealthCheckHandler) ReadyCheck(c *gin.Context) {
	ready, err := h.service.ReadyCheck(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	response.Success(c, ready)
}
