package http

import (
	"encoding/json"
	"net/http"
)

// Response represents standard API response structure
type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorData `json:"error,omitempty"`
	Meta    *MetaData  `json:"meta,omitempty"`
}

// ErrorData represents error information in API response
type ErrorData struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// MetaData represents pagination and additional metadata
type MetaData struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
}

// respondJSON sends a JSON response with the given status code
func respondJSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback error response if encoding fails
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(w http.ResponseWriter, data any) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// RespondSuccessWithMeta sends a successful JSON response with metadata
func RespondSuccessWithMeta(w http.ResponseWriter, data any, meta *MetaData) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// RespondCreated sends a 201 Created response
func RespondCreated(w http.ResponseWriter, data any) {
	respondJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// RespondNoContent sends a 204 No Content response
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RespondError sends an error JSON response
func RespondError(w http.ResponseWriter, statusCode int, code, message string) {
	respondJSON(w, statusCode, Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
		},
	})
}

// RespondErrorWithDetails sends an error JSON response with validation details
func RespondErrorWithDetails(w http.ResponseWriter, statusCode int, code, message string, details map[string]string) {
	respondJSON(w, statusCode, Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// Error code constants
const (
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeTooManyRequests    = "TOO_MANY_REQUESTS"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// RespondBadRequest sends a 400 Bad Request error response
func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusBadRequest, ErrCodeBadRequest, message)
}

// RespondUnauthorized sends a 401 Unauthorized error response
func RespondUnauthorized(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusUnauthorized, ErrCodeUnauthorized, message)
}

// RespondForbidden sends a 403 Forbidden error response
func RespondForbidden(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusForbidden, ErrCodeForbidden, message)
}

// RespondNotFound sends a 404 Not Found error response
func RespondNotFound(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusNotFound, ErrCodeNotFound, message)
}

// RespondConflict sends a 409 Conflict error response
func RespondConflict(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusConflict, ErrCodeConflict, message)
}

// RespondValidationError sends a 400 Bad Request error response with validation details
func RespondValidationError(w http.ResponseWriter, details map[string]string) {
	RespondErrorWithDetails(w, http.StatusBadRequest, ErrCodeValidation, "validation failed", details)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(w http.ResponseWriter) {
	RespondError(w, http.StatusInternalServerError, ErrCodeInternalServer, "internal server error")
}

// RespondTooManyRequests sends a 429 Too Many Requests error response
func RespondTooManyRequests(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusTooManyRequests, ErrCodeTooManyRequests, message)
}
