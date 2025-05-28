package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse struktur response yang konsisten untuk semua API
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *MetaInfo   `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo detail informasi error yang user-friendly
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MetaInfo informasi metadata untuk pagination dll
type MetaInfo struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// SuccessResponse mengembalikan response sukses
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

// CreatedResponse mengembalikan response untuk resource yang baru dibuat
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse mengembalikan response error
func ErrorResponse(c *gin.Context, statusCode int, errorCode string, message string, details string) {
	response := APIResponse{
		Success: false,
		Message: "Request failed",
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
	}
	c.JSON(statusCode, response)
}

// BadRequestResponse untuk error 400
func BadRequestResponse(c *gin.Context, message string, details string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// UnauthorizedResponse untuk error 401
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// ForbiddenResponse untuk error 403
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message, "")
}

// NotFoundResponse untuk error 404
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message, "")
}

// ConflictResponse untuk error 409
func ConflictResponse(c *gin.Context, message string, details string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, details)
}

// InternalServerErrorResponse untuk error 500
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, "")
}

// ValidationErrorResponse untuk error validasi input
func ValidationErrorResponse(c *gin.Context, validationErrors []string) {
	response := APIResponse{
		Success: false,
		Message: "Validation failed",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Input validation failed",
			Details: joinStrings(validationErrors, "; "),
		},
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusBadRequest, response)
}

// PaginatedResponse untuk response dengan pagination
func PaginatedResponse(c *gin.Context, message string, data interface{}, page, limit, total int) {
	totalPages := (total + limit - 1) / limit // Ceiling division

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &MetaInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

// NoContentResponse untuk response 204 (No Content)
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Helper function untuk join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}

// GetPaginationParams helper untuk extract pagination parameters
func GetPaginationParams(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

// CalculateOffset helper untuk menghitung offset untuk database query
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// HealthCheckResponse response khusus untuk health check
func HealthCheckResponse(c *gin.Context, serviceName string, status string, checks map[string]interface{}) {
	response := gin.H{
		"service":   serviceName,
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	}

	statusCode := http.StatusOK
	if status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
