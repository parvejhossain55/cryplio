package application

import (
	"database/sql"
	"fmt"

	domainidentity "cryplio/internal/domain/identity"
	"cryplio/internal/infrastructure/notification"
	identitypostgres "cryplio/internal/infrastructure/persistence/postgres/identity"
	httpapi "cryplio/internal/interfaces/http"
	"cryplio/pkg/config"
	"cryplio/pkg/database"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *gin.Engine
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := database.Open(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	userRepo := identitypostgres.NewUserRepository(db)
	emailClient := notification.NewSMTPClient(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.EmailFrom,
		cfg.FrontendURL,
	)
	passwordResetMailer := emailClient
	authService := domainidentity.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		cfg.RefreshTokenExpiry,
		cfg.CookieName,
		cfg.CookieSecure,
		cfg.CookieSameSite,
		cfg.IssuerName,
	).WithGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.OAuthRedirectURL,
	).WithPasswordResetMailer(passwordResetMailer)

	router := httpapi.SetupRouter(cfg, authService)

	return &App{
		Config: cfg,
		DB:     db,
		Router: router,
	}, nil
}
