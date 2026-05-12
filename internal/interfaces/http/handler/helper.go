package handler

import (
	"net/http"

	"cryplio/pkg/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// handleError maps domain AppErrors to appropriate HTTP status codes and
// writes a JSON error response. Internal errors return a generic message to
// avoid leaking implementation details.
func handleError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := err.Error()

	if apperrors.IsAppError(err) {
		appErr, _ := apperrors.GetAppError(err)
		switch appErr.Code {
		case apperrors.ErrCodeNotFound:
			status = http.StatusNotFound
		case apperrors.ErrCodeInvalidInput, apperrors.ErrCodeValidation:
			status = http.StatusBadRequest
		case apperrors.ErrCodeUnauthorized:
			status = http.StatusUnauthorized
		case apperrors.ErrCodeConflict:
			status = http.StatusConflict
		case apperrors.ErrCodeForbidden, apperrors.ErrCodePermissionDenied:
			status = http.StatusForbidden
		case apperrors.ErrCodeRateLimited:
			status = http.StatusTooManyRequests
		case apperrors.ErrCodeInternal:
			message = "an internal server error occurred"
		default:
			message = "an internal server error occurred"
		}
		if appErr.Code != apperrors.ErrCodeInternal {
			message = appErr.Message
		}
	}

	c.JSON(status, gin.H{"error": message})
}

// getUserIDFromContext extracts and parses the authenticated user ID from the
// Gin context (set by AuthMiddleware). Returns (uuid.Nil, false) and writes a
// 401 response if the ID is missing or malformed.
func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	userID, err := uuid.Parse(userIDRaw.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return uuid.Nil, false
	}
	return userID, true
}

// parseUUIDParam extracts a UUID from a URL parameter. Returns (uuid.Nil, false)
// and writes a 400 response if the parameter is missing or not a valid UUID.
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	val := c.Param(name)
	if val == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameter: " + name})
		return uuid.Nil, false
	}
	id, err := uuid.Parse(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return id, true
}
