package httpapi

import (
	"cryplio/internal/domain/identity"
	"cryplio/internal/interfaces/http/handler"
	"cryplio/internal/interfaces/http/middleware"
	"cryplio/pkg/config"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the gin engine
func SetupRouter(
	cfg *config.Config,
	authService identity.AuthService,
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
		})

		// Public
		public := v1.Group("/")
		{
			public.POST("/auth/register", authHandler.RegisterHandler)
			public.POST("/auth/login", authHandler.LoginHandler)
			public.POST("/auth/logout", authHandler.LogoutHandler)
			public.POST("/auth/refresh", authHandler.RefreshTokenHandler)
			public.GET("/auth/oauth/google", authHandler.GoogleAuthHandler)
			public.GET("/auth/oauth/google/callback", authHandler.GoogleCallbackHandler)
			public.POST("/auth/email/request", authHandler.RequestEmailVerificationHandler)
			public.POST("/auth/email/verify", authHandler.VerifyEmailHandler)
			public.POST("/auth/password/reset-request", authHandler.RequestPasswordResetHandler)
			public.POST("/auth/password/reset", authHandler.ResetPasswordHandler)
			public.POST("/auth/2fa/complete-login", authHandler.Complete2FALoginHandler)
		}

		// Authenticated
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			auth.GET("/users/me", authHandler.GetUserHandler)
			auth.PUT("/users/me", authHandler.UpdateUserHandler)
			// 2FA management
			auth.POST("/auth/2fa/setup", authHandler.Setup2FAHandler)
			auth.POST("/auth/2fa/verify", authHandler.Verify2FAHandler)
			auth.POST("/auth/2fa/disable", authHandler.Disable2FAHandler)
			// Session management
			auth.GET("/sessions", authHandler.GetSessionsHandler)
			auth.DELETE("/sessions/:tokenId", authHandler.DeleteSessionHandler)
			// TODO: other domain routes to be added
		}
	}

	return r
}
