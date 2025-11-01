package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/labstack/echo/v4"
)

// ErrorHandler handles all errors in a standardized way
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// Check if it's our custom AppError
	if appErr, ok := err.(*helpers.AppError); ok {
		_ = helpers.ResponseHttp(c, appErr.Code, appErr.Message, map[string]string{
			"details": appErr.Details,
		})
		return
	}

	// Check if it's echo.HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code := he.Code
		message := "An error occurred"

		if msg, ok := he.Message.(string); ok {
			message = msg
		}

		_ = helpers.ResponseHttp(c, code, message, nil)
		return
	}

	// Check if it's validation error
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, fieldErr := range validationErrs {
			errors[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
		}
		_ = helpers.ResponseHttp(c, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Default internal server error
	_ = helpers.ResponseHttp(c, http.StatusInternalServerError, "Internal server error", nil)
}

func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "datetime":
		return "Invalid date format, expected " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}
