// Package derr provides domain-level error definitions
package derr

import "errors"

var (
	// ErrNotFound indicates that a requested entity was not found
	ErrNotFound = errors.New("entity not found")

	// ErrConflict indicates a conflict in the state (e.g. duplicate key)
	ErrConflict = errors.New("entity conflict")

	// ErrValidation indicates a data validation failure
	ErrValidation = errors.New("validation failed")

	// ErrPermission indicates insufficient permissions
	ErrPermission = errors.New("permission denied")

	// ErrUnauthorized indicates authentication failure or lack of credentials
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInternal indicates an unexpected internal error
	ErrInternal = errors.New("internal error")
)
