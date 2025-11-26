// Package validator provides validation utilities for domain models
package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct using validation tags
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// formatValidationErrors formats validation errors into a readable string
func formatValidationErrors(errs validator.ValidationErrors) error {
	var errMsg string
	for i, err := range errs {
		if i > 0 {
			errMsg += "; "
		}
		errMsg += fmt.Sprintf("%s: %s", err.Field(), getErrorMsg(err))
	}
	return fmt.Errorf("validation error: %s", errMsg)
}

// getErrorMsg returns a human-readable error message for a validation error
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", err.Param())
	case "email":
		return "must be a valid email address"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", err.Param())
	default:
		return fmt.Sprintf("validation failed on '%s' tag", err.Tag())
	}
}
