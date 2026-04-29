package handler

// Config holds configuration for HTTP handlers.
type Config struct {
	CookieName     string
	CookieSecure   bool
	CookieSameSite string
	FrontendURL    string
}
