package auth

import (
	"time"

	"cryplio/internal/domain/identity"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

// mapUserStats converts UserStats to DTO.
func mapUserStats(s *identity.UserStats) dto.UserStatsDTO {
	if s == nil {
		return dto.UserStatsDTO{}
	}
	var lastTradeAt string
	if s.LastTradeAt != nil {
		lastTradeAt = s.LastTradeAt.Format(time.RFC3339)
	}
	return dto.UserStatsDTO{
		TotalTrades:           s.TotalTrades,
		SuccessfulTrades:      s.SuccessfulTrades,
		DisputeRate:           s.DisputeRate,
		AvgRating:             s.AvgRating,
		PositiveFeedbackCount: s.PositiveFeedbackCount,
		NeutralFeedbackCount:  s.NeutralFeedbackCount,
		NegativeFeedbackCount: s.NegativeFeedbackCount,
		TotalVolumeUSD:        s.TotalVolumeUSD,
		LastTradeAt:           lastTradeAt,
	}
}

// mapUser converts a domain User to a DTO UserResponse (empty stats).
func mapUser(u *identity.User) dto.UserResponse {
	if u == nil {
		return dto.UserResponse{}
	}
	var lastSeenAt, createdAt string
	if u.LastSeenAt != nil {
		lastSeenAt = u.LastSeenAt.Format(time.RFC3339)
	}
	createdAt = u.CreatedAt.Format(time.RFC3339)
	return dto.UserResponse{
		ID:            u.UserID.String(),
		Email:         u.Email,
		Username:      u.Username,
		EmailVerified: u.EmailVerified,
		TwoFAEnabled:  u.TwoFASecret != nil,
		AvatarURL:     u.AvatarURL,
		Bio:           u.Bio,
		LastSeenAt:    lastSeenAt,
		CreatedAt:     createdAt,
		Role:          string(u.Role),
		IsOnline:      u.IsOnline(),
		Stats:         dto.UserStatsDTO{},
	}
}

// mapUserWithStats converts a domain User + UserStats to a full DTO UserResponse.
func mapUserWithStats(u *identity.User, stats *identity.UserStats) dto.UserResponse {
	r := mapUser(u)
	r.Stats = mapUserStats(stats)
	return r
}

// setAuthCookie writes the access token into an HttpOnly cookie.
func setAuthCookie(c *gin.Context, cfg *Config, token string) {
	c.SetCookie(cfg.CookieName, token, 86400, "/", "", cfg.CookieSecure, true)
}

// clearAuthCookie removes the authentication cookie.
func clearAuthCookie(c *gin.Context, cfg *Config) {
	c.SetCookie(cfg.CookieName, "", -1, "/", "", cfg.CookieSecure, true)
}
