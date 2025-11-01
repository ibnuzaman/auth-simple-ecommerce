package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("username", validateUsername)
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return validate
}

// validatePhone validates Indonesian phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Indonesian phone number: starts with 08 or +62, 10-15 digits
	matched, _ := regexp.MatchString(`^(\+62|62|0)8[1-9][0-9]{6,11}$`, phone)
	return matched
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username: alphanumeric, underscore, 3-20 characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return matched
}

// ValidateStruct validates a struct using the global validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
