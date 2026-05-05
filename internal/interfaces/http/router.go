package httpapi

import (
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/platform"
	"cryplio/internal/domain/trading"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/handler"
	"cryplio/internal/interfaces/http/middleware"
	"cryplio/pkg/config"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the gin engine
func SetupRouter(
	cfg *config.Config,
	authService identity.AuthService,
	tradeService trading.TradeService,
	platformService platform.PlatformService,
	storage storage.ObjectStorage,
) *gin.Engine {
	// Set Gin mode
	if cfg.AppEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())
	if cfg.RateLimitEnabled {
		r.Use(middleware.RateLimitMiddleware())
	}

	// Health endpoints
	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/live", healthHandler.Liveness)
	r.GET("/ready", healthHandler.Readiness)

	// API v1
	v1 := r.Group("/api/v1")
	{
		authHandler := handler.NewAuthHandler(authService, &handler.Config{
			CookieName:         cfg.CookieName,
			CookieSecure:       cfg.CookieSecure,
			CookieSameSite:     cfg.CookieSameSite,
			FrontendURL:        cfg.FrontendURL,
			RefreshTokenExpiry: cfg.RefreshTokenExpiry,
			JWTSecret:          cfg.JWTSecret,
		}, storage)

		tradeHandler := handler.NewTradeHandler(tradeService)
		platformHandler := handler.NewPlatformHandler(platformService)

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
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			auth.GET("/users/me", authHandler.GetUserHandler)
			auth.PUT("/users/me", authHandler.UpdateUserHandler)
			auth.POST("/users/me/avatar", authHandler.UploadAvatarHandler)

			// User block management
			auth.POST("/users/me/block", authHandler.BlockUserHandler)
			auth.DELETE("/users/me/block/:blocked_id", authHandler.UnblockUserHandler)
			auth.GET("/users/me/block", authHandler.ListBlocksHandler)

			// 2FA management
			auth.POST("/auth/2fa/setup", authHandler.Setup2FAHandler)
			auth.POST("/auth/2fa/verify", authHandler.Verify2FAHandler)
			auth.POST("/auth/2fa/disable", authHandler.Disable2FAHandler)

			// Session management
			auth.GET("/sessions", authHandler.GetSessionsHandler)
			auth.DELETE("/sessions/:tokenId", authHandler.DeleteSessionHandler)

			// Trading (Authenticated)
			auth.POST("/marketplace/ads", tradeHandler.CreateAdHandler)
			auth.GET("/marketplace/my-ads", tradeHandler.ListMyAdsHandler)
			auth.PATCH("/marketplace/ads/:id/status", tradeHandler.ToggleAdStatusHandler)
			auth.POST("/marketplace/ads/:id/trades", tradeHandler.InitiateTradeHandler)
			auth.GET("/marketplace/trades", tradeHandler.ListTradesHandler)
			auth.GET("/marketplace/trades/:id", tradeHandler.GetTradeHandler)
			auth.PATCH("/marketplace/trades/:id/status", tradeHandler.UpdateTradeStatusHandler)
			auth.GET("/marketplace/trades/:id/messages", tradeHandler.GetChatHistoryHandler)
			auth.POST("/marketplace/trades/:id/messages", tradeHandler.SendMessageHandler)

			// Admin Routes
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminRoleMiddleware())
			{
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
			}
		}
	}

	return r
}
