package application

import (
	"database/sql"
	"fmt"

	domaindispute "cryplio/internal/domain/dispute"
	domainidentity "cryplio/internal/domain/identity"
	domainnotification "cryplio/internal/domain/notification"
	"cryplio/internal/domain/platform"
	domaintraing "cryplio/internal/domain/trading"
	domainwallet "cryplio/internal/domain/wallet"
	"cryplio/internal/infrastructure/notification"
	disputepostgres "cryplio/internal/infrastructure/persistence/postgres/dispute"
	identitypostgres "cryplio/internal/infrastructure/persistence/postgres/identity"
	notificationpostgres "cryplio/internal/infrastructure/persistence/postgres/notification"
	platformpostgres "cryplio/internal/infrastructure/persistence/postgres/platform"
	tradingpostgres "cryplio/internal/infrastructure/persistence/postgres/trading"
	walletpostgres "cryplio/internal/infrastructure/persistence/postgres/wallet"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/infrastructure/worker"
	httpapi "cryplio/internal/interfaces/http"
	"cryplio/pkg/config"
	"cryplio/pkg/database"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config    *config.Config
	DB        *sql.DB
	Router    *gin.Engine
	Worker    *worker.Worker
	Scheduler *worker.Scheduler
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
	).WithPasswordResetMailer(emailClient)

	notificationRepo := notificationpostgres.NewNotificationRepository(db)
	notificationService := domainnotification.NewService(notificationRepo, emailClient)

	disputeRepo := disputepostgres.NewDisputeRepository(db)
	disputeService := domaindispute.NewService(disputeRepo)

	tradeRepo := tradingpostgres.NewTradeRepository(db)
	tradeService := domaintraing.NewTradeService(tradeRepo, userRepo, disputeRepo)

	platformRepo := platformpostgres.NewPlatformRepository(db)
	platformService := platform.NewPlatformService(platformRepo)
	walletRepo := walletpostgres.NewWalletRepository(db)
	walletService := domainwallet.NewService(walletRepo)

	storage, err := storage.NewS3Storage(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize storage: %w", err)
	}

	router := httpapi.SetupRouter(cfg, authService, tradeService, platformService, walletService, disputeService, notificationService, storage)

	asynqWorker := worker.NewWorker(cfg, tradeService)
	asynqScheduler := worker.NewScheduler(cfg)

	return &App{
		Config:    cfg,
		DB:        db,
		Router:    router,
		Worker:    asynqWorker,
		Scheduler: asynqScheduler,
	}, nil
}
