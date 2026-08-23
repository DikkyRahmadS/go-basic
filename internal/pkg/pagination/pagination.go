package pagination

import (
	"gorm.io/gorm"

	"go-basic/internal/pkg/response"
)

type PageRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func CalculateMeta(page, limit int, totalItems int64) *response.Metadata {
	if limit <= 0 {
		return nil
	}

	if page < 1 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)
	if totalPages < 1 && totalItems == 0 {
		totalPages = 0
	}

	return &response.Metadata{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func GetOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit <= 0 {
			return db
		}
		if limit > 100 {
			limit = 100
		}
		offset := GetOffset(page, limit)
		return db.Limit(limit).Offset(offset)
	}
}
