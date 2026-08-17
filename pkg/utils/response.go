package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents the standard JSON envelope for all API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondJSON writes a standard JSON response with the given status code
func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// Fallback if encoding fails
			http.Error(w, `{"success":false,"error":"Internal server error encoding response"}`, http.StatusInternalServerError)
		}
	}
}

// ==========================================
// SUCCESS RESPONSES
// ==========================================

// RespondSuccess sends a success envelope with optional message and data
func RespondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	RespondJSON(w, statusCode, response)
}

// SuccessResponse is an alias for RespondSuccess for naming flexibility
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	RespondSuccess(w, statusCode, message, data)
}

// ==========================================
// ERROR RESPONSES
// ==========================================

// RespondError sends a standardized error response envelope
func RespondError(w http.ResponseWriter, statusCode int, errMessage string) {
	response := APIResponse{
		Success: false,
		Error:   errMessage,
	}
	RespondJSON(w, statusCode, response)
}

// ErrorResponse is an alias for RespondError for naming flexibility
func ErrorResponse(w http.ResponseWriter, statusCode int, errMessage string) {
	RespondError(w, statusCode, errMessage)
}

// ==========================================
// VALIDATION ERROR RESPONSES
// ==========================================

// RespondValidationError sends details when request payload binding or validation fails
func RespondValidationError(w http.ResponseWriter, details interface{}) {
	response := APIResponse{
		Success: false,
		Message: "Validation failed for request payload",
		Error:   details,
	}
	RespondJSON(w, http.StatusBadRequest, response)
}

// ValidationErrorResponse is an alias for RespondValidationError
func ValidationErrorResponse(w http.ResponseWriter, details interface{}) {
	RespondValidationError(w, details)
}
