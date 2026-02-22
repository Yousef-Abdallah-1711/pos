package http

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationErrorResponse represents a validation error response with field-level details
type ValidationErrorResponse struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    string                  `json:"code"`
	Status  int                     `json:"status"`
	Details []ValidationErrorDetail `json:"details,omitempty"`
}

// ErrorCode represents standard error codes
type ErrorCode string

const (
	// Validation errors
	ErrValidationFailed     ErrorCode = "VALIDATION_ERROR"
	ErrInvalidHexColor      ErrorCode = "INVALID_HEX_COLOR"
	ErrInvalidURL           ErrorCode = "INVALID_URL"
	ErrInvalidFont          ErrorCode = "INVALID_FONT"
	ErrOutOfRange           ErrorCode = "OUT_OF_RANGE"
	ErrRequiredField        ErrorCode = "REQUIRED_FIELD"
	ErrSQLInjectionDetected ErrorCode = "SQL_INJECTION_DETECTED"

	// Resource errors
	ErrNotFound  ErrorCode = "NOT_FOUND"
	ErrConflict  ErrorCode = "CONFLICT"
	ErrDuplicate ErrorCode = "DUPLICATE"

	// Authentication/Authorization errors
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"

	// Request errors
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"

	// Server errors
	ErrInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// SuccessResponse writes a successful response
func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorResponseWithMessage writes an error response with a message
func ErrorResponseWithMessage(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Error:   message,
		Status:  statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

// CreatedResponse writes a 201 Created response
func CreatedResponse(w http.ResponseWriter, data interface{}) {
	SuccessResponse(w, data, http.StatusCreated)
}

// BadRequestResponse writes a 400 Bad Request response
func BadRequestResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusBadRequest)
}

// UnauthorizedResponse writes a 401 Unauthorized response
func UnauthorizedResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusUnauthorized)
}

// ForbiddenResponse writes a 403 Forbidden response
func ForbiddenResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusForbidden)
}

// NotFoundResponse writes a 404 Not Found response
func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusNotFound)
}

// ConflictResponse writes a 409 Conflict response
func ConflictResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusConflict)
}

// InternalErrorResponse writes a 500 Internal Server Error response
func InternalErrorResponse(w http.ResponseWriter, message string) {
	ErrorResponseWithMessage(w, message, http.StatusInternalServerError)
}

// OKResponse writes a 200 OK response
func OKResponse(w http.ResponseWriter, data interface{}) {
	SuccessResponse(w, data, http.StatusOK)
}

// NoContentResponse writes a 204 No Content response
func NoContentResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// DecodeJSONBody decodes a JSON request body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if r.Body == nil {
		BadRequestResponse(w, "request body cannot be empty")
		return false
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		BadRequestResponse(w, "invalid request body: "+err.Error())
		return false
	}

	return true
}

// GetPathParam extracts a path parameter from the URL
func GetPathParam(r *http.Request, paramName string) string {
	// This is a helper function - actual implementation depends on router
	// For chi router: use chi.URLParam(r, paramName)
	// For gorilla/mux: use mux.Vars(r)[paramName]
	return ""
}

// GetQueryParam gets a query parameter with default value
func GetQueryParam(r *http.Request, paramName string, defaultValue string) string {
	value := r.URL.Query().Get(paramName)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt gets a query parameter as integer
func GetQueryParamInt(r *http.Request, paramName string, defaultValue int) int {
	value := r.URL.Query().Get(paramName)
	if value == "" {
		return defaultValue
	}

	var result int
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return defaultValue
	}

	return result
}

// GetQueryParamBool gets a query parameter as boolean
func GetQueryParamBool(r *http.Request, paramName string) bool {
	value := r.URL.Query().Get(paramName)
	return value == "true" || value == "1" || value == "yes"
}

// ValidationErrorResponse writes a validation error response with field-level details
func ValidationErrorResponseWithDetails(w http.ResponseWriter, errorMsg string, details []ValidationErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Success: false,
		Error:   errorMsg,
		Code:    string(ErrValidationFailed),
		Status:  http.StatusBadRequest,
		Details: details,
	}

	json.NewEncoder(w).Encode(response)
}

// ConflictResponseWithDetails writes a 409 Conflict response with additional details
func ConflictResponseWithDetails(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusConflict,
	}

	json.NewEncoder(w).Encode(response)
}

// NotFoundResponseWithCode writes a 404 Not Found response with error code
func NotFoundResponseWithCode(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusNotFound,
	}

	json.NewEncoder(w).Encode(response)
}

// BadRequestResponseWithCode writes a 400 Bad Request response with error code
func BadRequestResponseWithCode(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusBadRequest,
	}

	json.NewEncoder(w).Encode(response)
}

// InternalErrorResponseWithCode writes a 500 Internal Server Error response with error code
func InternalErrorResponseWithCode(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusInternalServerError,
	}

	json.NewEncoder(w).Encode(response)
}

// UnauthorizedResponseWithCode writes a 401 Unauthorized response with error code
func UnauthorizedResponseWithCode(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusUnauthorized,
	}

	json.NewEncoder(w).Encode(response)
}

// writeJSON is a general-purpose helper for writing JSON responses with any status code
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ForbiddenResponseWithCode writes a 403 Forbidden response with error code
func ForbiddenResponseWithCode(w http.ResponseWriter, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Error:   message,
		Code:    string(code),
		Status:  http.StatusForbidden,
	}

	json.NewEncoder(w).Encode(response)
}
