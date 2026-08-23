package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"erp-internal/internal/modules/health/models"
)

type HealthService interface {
	Check(ctx context.Context) (*models.HealthResponse, bool)
}

type healthService struct {
	db *gorm.DB
}

func NewHealthService(db *gorm.DB) HealthService {
	return &healthService{
		db: db,
	}
}

func (s *healthService) Check(ctx context.Context) (*models.HealthResponse, bool) {
	dbStatus := "connected"
	isHealthy := true

	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		dbStatus = "disconnected"
		isHealthy = false
	}

	status := "healthy"
	if !isHealthy {
		status = "unhealthy"
	}

	return &models.HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services: map[string]string{
			"database": dbStatus,
		},
	}, isHealthy
}
