package handler

import (
	"net/http"

	"cryplio/internal/domain/identity"
	"cryplio/internal/interfaces/http/dto"
	"cryplio/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

// handleError maps domain errors to HTTP responses.
func handleError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if apperrors.IsAppError(err) {
		appErr, _ := apperrors.GetAppError(err)
		switch appErr.Code {
		case apperrors.ErrCodeNotFound, apperrors.ErrCodeInvalidInput:
			status = http.StatusBadRequest
		case apperrors.ErrCodeUnauthorized:
			status = http.StatusUnauthorized
		case apperrors.ErrCodeConflict:
			status = http.StatusConflict
		case apperrors.ErrCodeForbidden:
			status = http.StatusForbidden
		case apperrors.ErrCodeRateLimited:
			status = http.StatusTooManyRequests
		default:
			status = http.StatusInternalServerError
		}
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

// mapUser converts a domain User to a DTO UserResponse.
func mapUser(u *identity.User) dto.UserResponse {
	if u == nil {
		return dto.UserResponse{}
	}
	kycLevel := 0
	switch u.KYCLevel {
	case identity.KYCLevel1:
		kycLevel = 1
	case identity.KYCLevel2:
		kycLevel = 2
	case identity.KYCLevel3:
		kycLevel = 3
	}
	return dto.UserResponse{
		ID:       u.UserID.String(),
		Email:    u.Email,
		Username: u.Username,
		KYCLevel: kycLevel,
	}
}

// setAuthCookie sets the authentication cookie.
func setAuthCookie(c *gin.Context, cfg *Config, token string) {
	// 86400 seconds = 24h
	c.SetCookie(cfg.CookieName, token, 86400, "/", "", cfg.CookieSecure, true)
}

// clearAuthCookie removes the authentication cookie.
func clearAuthCookie(c *gin.Context, cfg *Config) {
	c.SetCookie(cfg.CookieName, "", -1, "/", "", cfg.CookieSecure, true)
}
