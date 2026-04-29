package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"cryplio/pkg/database"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string
	JWTSecret  string
	JWTExpiry  time.Duration

	// Cookie settings
	CookieName     string
	CookieMaxAge   int
	CookieHTTPOnly bool
	CookieSecure   bool
	CookieSameSite string

	// Database
	Database *database.Config

	// Service URLs (kept for compatibility but not used in monolith)
	AuthServiceURL     string
	UserServiceURL     string
	WalletServiceURL   string
	TradeEngineURL     string
	KYCEngineURL       string
	NotificationURL    string
	DisputeServiceURL  string
	MerchantServiceURL string

	// OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	OAuthRedirectURL   string
	FrontendURL        string

	// Rate limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CORS
	CorsAllowedOrigins []string
}

func Load() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	rateLimitEnabled, _ := strconv.ParseBool(getEnvCompat("RATE_LIMIT_ENABLED", "true"))
	rateLimitRequests, _ := strconv.Atoi(getEnvCompat("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := time.ParseDuration(getEnvCompat("RATE_LIMIT_WINDOW", "1m"))
	jwtExpiry, _ := time.ParseDuration(getEnvCompat("JWT_EXPIRY", "24h"))

	// Database configuration
	dbCfg := database.DefaultConfig()
	dbCfg.Host = getEnvCompat("DB_HOST", dbCfg.Host)
	dbCfg.Port, _ = strconv.Atoi(getEnvCompat("DB_PORT", strconv.Itoa(dbCfg.Port)))
	dbCfg.User = getEnvCompat("DB_USER", dbCfg.User)
	dbCfg.Password = getEnvCompat("DB_PASSWORD", dbCfg.Password)
	dbCfg.DBName = getEnvCompat("DB_NAME", dbCfg.DBName)

	cfg := &Config{
		AppEnv:     getEnvCompat("APP_ENV", "development"),
		ServerPort: getEnvCompat("SERVER_PORT", "8080"),
		JWTSecret:  getEnvCompat("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiry:  jwtExpiry,

		// Cookie configuration
		CookieName:     getEnvCompat("COOKIE_NAME", "auth_token"),
		CookieMaxAge:   parseJWTExpiryToSeconds(jwtExpiry),
		CookieHTTPOnly: true,
		CookieSecure:   getEnvCompat("COOKIE_SECURE", "false") == "true",
		CookieSameSite: getEnvCompat("COOKIE_SAME_SITE", "lax"),

		Database: dbCfg,

		AuthServiceURL:     getEnvCompat("AUTH_SERVICE_URL", "http://localhost:8080"),
		UserServiceURL:     getEnvCompat("USER_SERVICE_URL", "http://localhost:8080"),
		WalletServiceURL:   getEnvCompat("WALLET_SERVICE_URL", "http://localhost:8081"),
		TradeEngineURL:     getEnvCompat("TRADE_ENGINE_URL", "http://localhost:8082"),
		KYCEngineURL:       getEnvCompat("KYC_SERVICE_URL", "http://localhost:8083"),
		NotificationURL:    getEnvCompat("NOTIFICATION_URL", "http://localhost:8084"),
		DisputeServiceURL:  getEnvCompat("DISPUTE_SERVICE_URL", "http://localhost:8085"),
		MerchantServiceURL: getEnvCompat("MERCHANT_SERVICE_URL", "http://localhost:8082"),

		// OAuth
		GoogleClientID:     getEnvCompat("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnvCompat("GOOGLE_CLIENT_SECRET", ""),
		OAuthRedirectURL:   getEnvCompat("OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/google/callback"),
		FrontendURL:        getEnvCompat("FRONTEND_URL", "http://localhost:3000"),

		RateLimitEnabled:   rateLimitEnabled,
		RateLimitRequests:  rateLimitRequests,
		RateLimitWindow:    rateLimitWindow,
		CorsAllowedOrigins: []string{getEnvCompat("CORS_ALLOWED_ORIGINS", "*")},
	}

	return cfg, nil
}

func loadEnvFile() error {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	envFile := ".env"
	if appEnv == "production" {
		envFile = ".env.prod"
	}

	if err := godotenv.Load(envFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func parseJWTExpiryToSeconds(d time.Duration) int {
	if d < time.Second {
		return 86400 // default 24h
	}
	return int(d.Seconds())
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvCompat(key, defaultValue string) string {
	switch key {
	case "SERVER_PORT":
		if value, exists := os.LookupEnv("APP_PORT"); exists {
			return value
		}
	case "JWT_SECRET":
		if value, exists := os.LookupEnv("APP_SECRET"); exists {
			return value
		}
	case "JWT_EXPIRY":
		if value, exists := os.LookupEnv("JWT_EXPIRY_HOURS"); exists {
			return value + "h"
		}
	case "RATE_LIMIT_REQUESTS":
		if value, exists := os.LookupEnv("RATE_LIMIT_RPS"); exists {
			return value
		}
	case "RATE_LIMIT_WINDOW":
		if value, exists := os.LookupEnv("RATE_LIMIT_BURST"); exists {
			if _, err := strconv.Atoi(value); err == nil {
				return "1m"
			}
		}
	}

	return getEnv(key, defaultValue)
}
