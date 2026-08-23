package models

import "time"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindAllUsersRequest struct {
	Search string `form:"search"`
	Name   string `form:"name"`
	Email  string `form:"email"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email"`
}
