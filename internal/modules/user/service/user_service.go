package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"erp-internal/internal/modules/user/models"
	"erp-internal/internal/modules/user/repository"
	"erp-internal/internal/pkg/apperror"
	"erp-internal/internal/pkg/database"
)

type UserService interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	FindAll(ctx context.Context, req *models.FindAllUsersRequest) ([]*models.UserResponse, int64, error)
	Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	db         *gorm.DB
	repository repository.UserRepository
	validate   *validator.Validate
}

func NewUserService(
	db *gorm.DB,
	repository repository.UserRepository,
	validate *validator.Validate,
) UserService {
	return &userService{
		db:         db,
		repository: repository,
		validate:   validate,
	}
}

func (s *userService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := database.WithTx(ctx, s.db, func(txCtx context.Context) error {
		existingUser, err := s.repository.FindByEmail(txCtx, user.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existingUser != nil {
			return apperror.Conflict(fmt.Sprintf("email %s already exists", req.Email))
		}

		if err := s.repository.Create(txCtx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return models.ToUserResponse(user), nil
}

func (s *userService) FindAll(ctx context.Context, req *models.FindAllUsersRequest) ([]*models.UserResponse, int64, error) {
	users, total, err := s.repository.FindAll(ctx, req)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}

	return models.ToUserResponses(users), total, nil
}

func (s *userService) Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	var updatedUser *models.User

	if err := database.WithTx(ctx, s.db, func(txCtx context.Context) error {
		user, err := s.repository.FindByID(txCtx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("user not found")
			}
			return err
		}

		if req.Email != user.Email {
			existingUser, err := s.repository.FindByEmail(txCtx, req.Email)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if existingUser != nil && existingUser.ID != id {
				return apperror.Conflict(fmt.Sprintf("email %s already exists", req.Email))
			}
		}

		user.Name = req.Name
		user.Email = req.Email

		if err := s.repository.Update(txCtx, user); err != nil {
			return err
		}

		updatedUser = user
		return nil
	}); err != nil {
		return nil, err
	}

	return models.ToUserResponse(updatedUser), nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if err := database.WithTx(ctx, s.db, func(txCtx context.Context) error {
		_, err := s.repository.FindByID(txCtx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("user not found")
			}
			return err
		}

		if err := s.repository.Delete(txCtx, id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
