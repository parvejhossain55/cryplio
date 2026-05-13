package httpapi

import (
	"context"
	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/market"
	"cryplio/internal/domain/notification"
	"cryplio/internal/domain/platform"
	"cryplio/internal/domain/trading"
	"cryplio/internal/domain/wallet"
	"cryplio/internal/infrastructure/persistence/redis"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/handler"
	authh "cryplio/internal/interfaces/http/handler/auth"
	disputeh "cryplio/internal/interfaces/http/handler/dispute"
	marketh "cryplio/internal/interfaces/http/handler/market"
	notificationh "cryplio/internal/interfaces/http/handler/notification"
	platformh "cryplio/internal/interfaces/http/handler/platform"
	tradeh "cryplio/internal/interfaces/http/handler/trade"
	walleth "cryplio/internal/interfaces/http/handler/wallet"
	"cryplio/internal/interfaces/http/middleware"
	"cryplio/internal/interfaces/websocket"
	"cryplio/pkg/config"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the gin engine.
// rateLimiter is optional: pass a non-nil redis.RateLimiter for distributed
// rate limiting (recommended in production). If nil, an in-memory limiter is used.
func SetupRouter(
	cfg *config.Config,
	authService identity.AuthService,
	tradeService trading.TradeService,
	platformService platform.PlatformService,
	walletService wallet.Service,
	disputeService dispute.Service,
	notificationService notification.Service,
	storage storage.ObjectStorage,
	rateService market.RateService,
	wsService websocket.Service,
	rateLimiter redis.RateLimiter,
	pingFn handler.PingFunc,
) *gin.Engine {
	// Set Gin mode
	if cfg.AppEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CorsAllowedOrigins))
	r.Use(middleware.LoggingMiddleware())
	if cfg.RateLimitEnabled {
		r.Use(middleware.RateLimitMiddleware(rateLimiter))
	}

	// Health endpoints
	healthHandler := handler.NewHealthHandler(pingFn)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/live", healthHandler.Liveness)
	r.GET("/ready", healthHandler.Readiness)

	// API v1
	v1 := r.Group("/api/v1")
	{
		authHandler := authh.NewAuthHandler(authService, &authh.Config{
			CookieName:         cfg.CookieName,
			CookieSecure:       cfg.CookieSecure,
			CookieSameSite:     cfg.CookieSameSite,
			FrontendURL:        cfg.FrontendURL,
			RefreshTokenExpiry: cfg.RefreshTokenExpiry,
			JWTSecret:          cfg.JWTSecret,
		}, storage)
		adminHandler := authh.NewAdminHandler(authService, tradeService, disputeService, walletService)

		tradeHandler := tradeh.NewTradeHandler(tradeService, storage, wsService)
		platformHandler := platformh.NewPlatformHandler(platformService)
		walletHandler := walleth.NewWalletHandler(walletService, authService)
		disputeHandler := disputeh.NewDisputeHandler(disputeService, storage)
		notificationHandler := notificationh.NewNotificationHandler(notificationService)

		// Public routes
		v1.POST("/auth/register", authHandler.RegisterHandler)
		v1.POST("/auth/login", authHandler.LoginHandler)
		v1.POST("/auth/logout", authHandler.LogoutHandler)
		v1.POST("/auth/refresh", authHandler.RefreshTokenHandler)
		v1.GET("/auth/oauth/google", authHandler.GoogleAuthHandler)
		v1.GET("/auth/oauth/google/callback", authHandler.GoogleCallbackHandler)
		v1.POST("/auth/email/request", authHandler.RequestEmailVerificationHandler)
		v1.POST("/auth/email/verify", authHandler.VerifyEmailHandler)
		v1.POST("/auth/password/reset-request", authHandler.RequestPasswordResetHandler)
		v1.POST("/auth/password/reset", authHandler.ResetPasswordHandler)
		v1.POST("/auth/2fa/complete-login", authHandler.Complete2FALoginHandler)
		v1.GET("/users/username/:username", authHandler.GetUserByUsernameHandler)

		// Marketplace (Public)
		v1.GET("/marketplace/ads", tradeHandler.ListAdsHandler)

		// Authenticated Routes
		// Market Rates (Public)
		marketHandler := marketh.NewMarketHandler(rateService)
		v1.GET("/market/rates", marketHandler.GetAllRatesHandler)
		v1.GET("/market/rates/:crypto", marketHandler.GetRatesHandler)
		v1.GET("/market/rates/:crypto/:fiat", marketHandler.GetRateHandler)

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			auth.GET("/users/me", authHandler.GetUserHandler)
			auth.PUT("/users/me", authHandler.UpdateUserHandler)
			auth.POST("/users/me/avatar", authHandler.UploadAvatarHandler)

			// 2FA management
			auth.POST("/auth/2fa/setup", authHandler.Setup2FAHandler)
			auth.POST("/auth/2fa/verify", authHandler.Verify2FAHandler)
			auth.POST("/auth/2fa/disable", authHandler.Disable2FAHandler)

			// Session management
			auth.GET("/sessions", authHandler.GetSessionsHandler)
			auth.DELETE("/sessions/:tokenId", authHandler.DeleteSessionHandler)

			// User Payment Methods
			auth.GET("/users/me/payment-methods", authHandler.ListPaymentMethodsHandler)
			auth.POST("/users/me/payment-methods", authHandler.CreatePaymentMethodHandler)
			auth.PUT("/users/me/payment-methods/:id", authHandler.UpdatePaymentMethodHandler)
			auth.DELETE("/users/me/payment-methods/:id", authHandler.DeletePaymentMethodHandler)
			auth.PATCH("/users/me/payment-methods/:id/default", authHandler.SetDefaultPaymentMethodHandler)

			// Trading (Authenticated)
			auth.POST("/marketplace/ads", tradeHandler.CreateAdHandler)
			auth.GET("/marketplace/my-ads", tradeHandler.ListMyAdsHandler)
			auth.PUT("/marketplace/ads/:id", tradeHandler.UpdateAdHandler)
			auth.DELETE("/marketplace/ads/:id", tradeHandler.DeleteAdHandler)
			auth.PATCH("/marketplace/ads/:id/status", tradeHandler.ToggleAdStatusHandler)
			auth.POST("/marketplace/ads/:id/trades", tradeHandler.InitiateTradeHandler)
			auth.GET("/marketplace/trades", tradeHandler.ListTradesHandler)
			auth.GET("/marketplace/trades/:id", tradeHandler.GetTradeHandler)
			auth.PATCH("/marketplace/trades/:id/status", tradeHandler.UpdateTradeStatusHandler)
			auth.POST("/marketplace/trades/:id/dispute", tradeHandler.DisputeTradeHandler)
			auth.POST("/marketplace/trades/:id/feedback", tradeHandler.LeaveFeedbackHandler)
			auth.POST("/disputes/:id/evidence", disputeHandler.UploadEvidenceHandler)
			auth.GET("/marketplace/trades/:id/messages", tradeHandler.GetChatHistoryHandler)
			auth.POST("/marketplace/trades/:id/messages", tradeHandler.SendMessageHandler)

			// Wallet
			auth.GET("/wallet/balance", walletHandler.GetBalancesHandler)
			auth.GET("/wallet/deposit/:crypto", walletHandler.GetDepositAddressHandler)
			auth.POST("/wallet/withdraw", walletHandler.WithdrawHandler)
			auth.GET("/wallet/transactions", walletHandler.GetTransactionsHandler)

			// Notifications
			auth.GET("/notifications", notificationHandler.GetNotificationsHandler)
			auth.PATCH("/notifications/:id/read", notificationHandler.MarkReadHandler)
			auth.GET("/notifications/preferences", notificationHandler.GetPreferencesHandler)
			auth.POST("/notifications/preferences", notificationHandler.SavePreferencesHandler)

			// Admin Routes
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminRoleMiddleware())
			{
				// Dashboard Stats
				admin.GET("/dashboard/stats", adminHandler.GetDashboardStatsHandler)
				admin.GET("/activity", adminHandler.GetActivityHandler)
				admin.GET("/alerts", adminHandler.GetAlertsHandler)

				// User Management
				admin.GET("/users", adminHandler.ListUsersHandler)
				admin.POST("/users/:id/suspend", adminHandler.SuspendUserHandler)
				admin.POST("/users/:id/unsuspend", adminHandler.UnsuspendUserHandler)
				admin.POST("/users/:id/ban", adminHandler.BanUserHandler)
				admin.POST("/users/:id/unban", adminHandler.UnbanUserHandler)

				// Trade Monitoring
				admin.GET("/trades", tradeHandler.ListAllTradesHandler)

				// Crypto Assets
				admin.POST("/crypto-assets", platformHandler.CreateCryptoAssetHandler)
				admin.GET("/crypto-assets", platformHandler.GetCryptoAssetsHandler)
				admin.GET("/crypto-assets/:id", platformHandler.GetCryptoAssetHandler)
				admin.PUT("/crypto-assets/:id", platformHandler.UpdateCryptoAssetHandler)
				admin.DELETE("/crypto-assets/:id", platformHandler.DeleteCryptoAssetHandler)

				// Fiat Currencies
				admin.POST("/fiat-currencies", platformHandler.CreateFiatCurrencyHandler)
				admin.GET("/fiat-currencies", platformHandler.GetFiatCurrenciesHandler)
				admin.GET("/fiat-currencies/:id", platformHandler.GetFiatCurrencyHandler)
				admin.PUT("/fiat-currencies/:id", platformHandler.UpdateFiatCurrencyHandler)
				admin.DELETE("/fiat-currencies/:id", platformHandler.DeleteFiatCurrencyHandler)

				// Payment Methods
				admin.POST("/payment-methods", platformHandler.CreatePaymentMethodHandler)
				admin.GET("/payment-methods", platformHandler.GetPaymentMethodsHandler)
				admin.GET("/payment-methods/:id", platformHandler.GetPaymentMethodHandler)
				admin.PUT("/payment-methods/:id", platformHandler.UpdatePaymentMethodHandler)
				admin.DELETE("/payment-methods/:id", platformHandler.DeletePaymentMethodHandler)

				// Withdrawal Approval
				admin.GET("/withdrawals/pending", walletHandler.ListPendingWithdrawalsHandler)
				admin.POST("/withdrawals/:id/approve", walletHandler.ApproveWithdrawalHandler)
				admin.POST("/withdrawals/:id/reject", walletHandler.RejectWithdrawalHandler)

				// Disputes management
				admin.GET("/disputes", disputeHandler.ListDisputesHandler)
				admin.GET("/disputes/:id", disputeHandler.GetDisputeHandler)
				admin.POST("/disputes/:id/assign", disputeHandler.AssignDisputeHandler)
				admin.POST("/disputes/:id/resolve", disputeHandler.ResolveDisputeHandler)
			}
		}
	}

	// Register WebSocket endpoint with Gin
	wsService.Start(context.Background())
	r.GET("/ws", func(c *gin.Context) {
		wsService.HandleWebSocket(c.Writer, c.Request)
	})

	return r
}
