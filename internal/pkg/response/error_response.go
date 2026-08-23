package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"erp-internal/internal/pkg/apperror"
)

func ErrorResponse(c *gin.Context, code int, err string, details []ErrorDetails) {
	c.JSON(code, &WebErrorResponse{
		Code:    code,
		Message: err,
		Details: details,
	})
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		ErrorResponse(c, appErr.Code, appErr.Message, nil)
		return
	}

	var valErrors validator.ValidationErrors
	if errors.As(err, &valErrors) {
		details := FormatValidationError(err)
		ErrorResponse(c, http.StatusBadRequest, "Validation failed", details)
		return
	}

	ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
}

func FormatValidationError(err error) []ErrorDetails {
	var details []ErrorDetails
	var valErrors validator.ValidationErrors

	if errors.As(err, &valErrors) {
		for _, fe := range valErrors {
			details = append(details, ErrorDetails{
				Field:   fe.Field(),
				Message: formatErrorMessage(fe),
			})
		}
	} else if err != nil {
		details = append(details, ErrorDetails{
			Field:   "general",
			Message: err.Error(),
		})
	}

	return details
}

func formatErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", fe.Field(), fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "nefield":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", fe.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	default:
		return fmt.Sprintf("%s is invalid (%s)", fe.Field(), fe.Tag())
	}
}

func (e *WebErrorResponse) Error() string {
	return e.Message
}
