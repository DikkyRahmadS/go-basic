package repository

import (
	"context"

	"gorm.io/gorm"

	"go-basic/internal/modules/user/models"
	"go-basic/internal/pkg/database"
	"go-basic/internal/pkg/pagination"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindAll(ctx context.Context, req *models.FindAllUsersRequest) ([]models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	db := database.GetDB(ctx, r.db)
	return db.Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	db := database.GetDB(ctx, r.db)
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	db := database.GetDB(ctx, r.db)
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	db := database.GetDB(ctx, r.db)
	return db.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	db := database.GetDB(ctx, r.db)
	return db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *userRepository) FindAll(ctx context.Context, req *models.FindAllUsersRequest) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := database.GetDB(ctx, r.db).Model(&models.User{})
	db = r.applyFilters(db, req)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := db.Order("name DESC")
	if req != nil {
		query = query.Scopes(pagination.Paginate(req.Page, req.Limit))
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) applyFilters(db *gorm.DB, req *models.FindAllUsersRequest) *gorm.DB {
	if req == nil {
		return db
	}

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		db = db.Where("name LIKE ? OR email LIKE ?", searchTerm, searchTerm)
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}

	return db
}
