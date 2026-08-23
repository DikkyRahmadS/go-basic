package apperror

import (
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "internal error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

func Internal(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Err:     err,
	}
}
