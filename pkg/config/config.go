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

	// Refresh token configuration
	RefreshTokenExpiry time.Duration

	// Cookie settings
	CookieName     string
	CookieMaxAge   int
	CookieHTTPOnly bool
	CookieSecure   bool
	CookieSameSite string

	// Database
	Database *database.Config

	// OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	OAuthRedirectURL   string
	FrontendURL        string

	// Email
	EmailFrom string

	// SMTP settings
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// Rate limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CORS
	CorsAllowedOrigins []string

	// 2FA Issuer
	IssuerName string

	// S3/MinIO configuration
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3UseSSL          bool
	S3BucketName      string
	S3PublicBaseURL   string // Optional: public URL for accessing objects (if different from endpoint)

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Blockchain
	EthRPCURL             string
	EthPrivateKey         string
	EscrowContractAddress string
}

func Load() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	rateLimitEnabled, _ := strconv.ParseBool(getEnvCompat("RATE_LIMIT_ENABLED", "true"))
	rateLimitRequests, _ := strconv.Atoi(getEnvCompat("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := time.ParseDuration(getEnvCompat("RATE_LIMIT_WINDOW", "1m"))
	jwtExpiry, _ := time.ParseDuration(getEnvCompat("JWT_EXPIRY", "24h"))
	refreshTokenExpiry, _ := time.ParseDuration(getEnvCompat("REFRESH_TOKEN_EXPIRY", "168h")) // 7 days

	// Database configuration
	dbCfg := database.DefaultConfig()
	dbCfg.Host = getEnvCompat("DB_HOST", dbCfg.Host)
	dbCfg.Port, _ = strconv.Atoi(getEnvCompat("DB_PORT", strconv.Itoa(dbCfg.Port)))
	dbCfg.User = getEnvCompat("DB_USER", dbCfg.User)
	dbCfg.Password = getEnvCompat("DB_PASSWORD", dbCfg.Password)
	dbCfg.DBName = getEnvCompat("DB_NAME", dbCfg.DBName)

	cfg := &Config{
		AppEnv:             getEnvCompat("APP_ENV", "development"),
		ServerPort:         getEnvCompat("SERVER_PORT", "8080"),
		JWTSecret:          getEnvCompat("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiry:          jwtExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,

		// Cookie configuration
		CookieName:     getEnvCompat("COOKIE_NAME", "auth_token"),
		CookieMaxAge:   parseJWTExpiryToSeconds(jwtExpiry),
		CookieHTTPOnly: true,
		CookieSecure:   getEnvCompat("COOKIE_SECURE", "false") == "true",
		CookieSameSite: getEnvCompat("COOKIE_SAME_SITE", "strict"),

		Database: dbCfg,

		// OAuth
		GoogleClientID:     getEnvCompat("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnvCompat("GOOGLE_CLIENT_SECRET", ""),
		OAuthRedirectURL:   getEnvCompat("OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/google/callback"),
		FrontendURL:        getEnvCompat("FRONTEND_URL", "http://localhost:3000"),

		EmailFrom: getEnvCompat("EMAIL_FROM", ""),

		// SMTP settings
		SMTPHost: getEnvCompat("SMTP_HOST", ""),
		SMTPPort: func() int {
			if port, err := strconv.Atoi(getEnvCompat("SMTP_PORT", "587")); err == nil {
				return port
			}
			return 587 // default to TLS port
		}(),
		SMTPUsername: getEnvCompat("SMTP_USERNAME", ""),
		SMTPPassword: getEnvCompat("SMTP_PASSWORD", ""),

		RateLimitEnabled:   rateLimitEnabled,
		RateLimitRequests:  rateLimitRequests,
		RateLimitWindow:    rateLimitWindow,
		CorsAllowedOrigins: []string{getEnvCompat("CORS_ALLOWED_ORIGINS", "*")},
		IssuerName:         getEnvCompat("ISSUER_NAME", "Cryplio"),

		// S3/MinIO configuration
		S3Endpoint:        getEnvCompat("S3_ENDPOINT", "localhost:9000"),
		S3AccessKeyID:     getEnvCompat("S3_ACCESS_KEY_ID", "minioadmin"),
		S3SecretAccessKey: getEnvCompat("S3_SECRET_ACCESS_KEY", "minioadmin"),
		S3UseSSL:          getEnvCompat("S3_USE_SSL", "false") == "true",
		S3BucketName:      getEnvCompat("S3_BUCKET_NAME", "cryplio-storage"),
		S3PublicBaseURL:   getEnvCompat("S3_PUBLIC_BASE_URL", ""),

		// Redis
		RedisAddr:     getEnvCompat("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnvCompat("REDIS_PASSWORD", ""),
		RedisDB: func() int {
			db, _ := strconv.Atoi(getEnvCompat("REDIS_DB", "0"))
			return db
		}(),

		// Blockchain
		EthRPCURL:             getEnvCompat("ETH_RPC_URL", "http://localhost:8545"),
		EthPrivateKey:         getEnvCompat("ETH_PRIVATE_KEY", ""),
		EscrowContractAddress: getEnvCompat("ESCROW_CONTRACT_ADDRESS", ""),
	}

	if strings.ToLower(cfg.AppEnv) == "production" && cfg.JWTSecret == "your-secret-key-change-this" {
		return nil, errors.New("JWT_SECRET must be set to a secure value in production")
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
