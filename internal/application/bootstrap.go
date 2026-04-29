package application

import (
	"database/sql"
	"fmt"

	domainidentity "cryplio/internal/domain/identity"
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
	authService := domainidentity.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		cfg.CookieName,
		cfg.CookieSecure,
		cfg.CookieSameSite,
	).WithGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.OAuthRedirectURL,
	)

	router := httpapi.SetupRouter(cfg, authService)

	return &App{
		Config: cfg,
		DB:     db,
		Router: router,
	}, nil
}
