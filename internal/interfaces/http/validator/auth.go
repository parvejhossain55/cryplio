package validator

import (
	"strings"

	"cryplio/internal/interfaces/http/dto"
)

// NormalizeRegisterRequest applies delivery-layer normalization before the domain call.
func NormalizeRegisterRequest(req *dto.RegisterRequest) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)
}

// NormalizeLoginRequest applies delivery-layer normalization before the domain call.
func NormalizeLoginRequest(req *dto.LoginRequest) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
}

// NormalizeUpdateProfileRequest trims optional profile fields.
func NormalizeUpdateProfileRequest(req *dto.UpdateProfileRequest) {
	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		req.Username = &trimmed
	}
	if req.Bio != nil {
		trimmed := strings.TrimSpace(*req.Bio)
		req.Bio = &trimmed
	}
}
