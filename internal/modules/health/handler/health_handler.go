package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"erp-internal/internal/modules/health/service"
	"erp-internal/internal/pkg/response"
)

type HealthHandler struct {
	service service.HealthService
}

func NewHealthHandler(service service.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	result, isHealthy := h.service.Check(c.Request.Context())

	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	response.SuccessResponse(c, statusCode, "Health status", result)
}
