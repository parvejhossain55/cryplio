package handler

import "time"

// Config holds configuration for HTTP handlers.
type Config struct {
	CookieName         string
	CookieSecure       bool
	CookieSameSite     string
	FrontendURL        string
	RefreshTokenExpiry time.Duration
	JWTSecret          string
}
