package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// DecodeAndValidate decodes JSON request body and validates it
func DecodeAndValidate(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		if err == io.EOF {
			return fmt.Errorf("empty request body")
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}

// ValidationErrorsToMap converts validator errors to map for response
func ValidationErrorsToMap(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			errors[field] = getErrorMessage(fieldErr)
		}
	}

	return errors
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	tag := fieldErr.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fieldErr.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldErr.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fieldErr.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldErr.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fieldErr.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fieldErr.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// DecodeAndValidateRequest is a helper function that decodes, validates,
// and sends error response if validation fails
func (h *Handler) DecodeAndValidateRequest(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := DecodeAndValidate(r, v); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			RespondValidationError(w, ValidationErrorsToMap(validationErrs))
			return false
		}
		RespondBadRequest(w, err.Error())
		return false
	}
	return true
}
