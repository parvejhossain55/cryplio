package auth

import (
	"errors"
	"net/http"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/domain/identity"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
	sharedjwt "cryplio/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests.
// It is composed of several focused sub-files in the same package:
//
//   - auth.go                — core auth lifecycle (this file)
//   - profile.go             — user profile operations
//   - password.go            — email verification & password reset
//   - twofactor_session.go   — 2FA management & session CRUD
//   - payment_method.go      — user payment methods
type AuthHandler struct {
	authService identity.AuthService
	cfg         *Config
	storage     storage.ObjectStorage
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService identity.AuthService, cfg *Config, storage storage.ObjectStorage) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		storage:     storage,
	}
}

// ─── Registration ────────────────────────────────────────────────────────────

// RegisterHandler handles user registration.
// On success it automatically logs the new user in and returns auth tokens.
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeRegisterRequest(&req)

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	// Fire-and-forget: send verification email after account is created.
	_ = h.authService.RequestEmailVerification(c.Request.Context(), user.UserID)

	// Auto-login. New users cannot have 2FA, so Login will not return
	// TwoFactorRequiredError. We guard anyway for safety.
	access, refresh, _, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// User was created — return 201 without tokens so the client can log in manually.
		c.JSON(http.StatusCreated, gin.H{
			"message": "Registration successful. Please log in.",
			"user":    mapUser(user),
		})
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusCreated, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// ─── Login / Logout ───────────────────────────────────────────────────────────

// LoginHandler authenticates a user and issues access + refresh tokens.
// If the account has 2FA enabled, it returns a temporary token for the
// second-factor step instead.
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeLoginRequest(&req)

	access, refresh, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		var twoFactorErr identity.TwoFactorRequiredError
		if errors.As(err, &twoFactorErr) {
			c.JSON(http.StatusOK, gin.H{
				"requires_2fa": true,
				"temp_token":   twoFactorErr.TempToken,
				"user":         mapUser(twoFactorErr.User),
			})
			return
		}
		basehandler.HandleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// LogoutHandler invalidates the user's refresh-token session and clears auth cookies.
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	refreshToken, _ := c.Cookie(h.cfg.CookieName + "_refresh")
	if refreshToken != "" {
		if claims, err := sharedjwt.Parse(h.cfg.JWTSecret, refreshToken); err == nil {
			if jti, ok := claims["jti"].(string); ok {
				_ = h.authService.Logout(c.Request.Context(), jti)
			}
		}
	}
	clearAuthCookie(c, h.cfg)
	h.clearRefreshCookie(c, h.cfg)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ─── Token Refresh ────────────────────────────────────────────────────────────

// RefreshTokenHandler rotates the refresh token and issues a new access token.
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie(h.cfg.CookieName + "_refresh")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing refresh token"})
		return
	}

	access, refresh, user, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// ─── Two-Factor Login Completion ─────────────────────────────────────────────

// Complete2FALoginHandler finalises a pending 2FA login challenge.
// The client must supply the short-lived temp_token received from LoginHandler
// together with the current TOTP code.
func (h *AuthHandler) Complete2FALoginHandler(c *gin.Context) {
	var req dto.TwoFactorLoginCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	access, refresh, user, err := h.authService.Complete2FALogin(c.Request.Context(), req.TempToken, req.Code)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// ─── Google OAuth ─────────────────────────────────────────────────────────────

// GoogleAuthHandler initiates the Google OAuth 2.0 authorization flow.
func (h *AuthHandler) GoogleAuthHandler(c *gin.Context) {
	oauthURL := h.authService.GoogleOAuthURL()
	if oauthURL == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "OAuth not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, oauthURL)
}

// GoogleCallbackHandler receives the OAuth authorization code, exchanges it for
// tokens, and redirects the user to the frontend dashboard.
func (h *AuthHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	accessToken, refreshToken, _, err := h.authService.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, accessToken)
	h.setRefreshCookie(c, refreshToken)
	c.Redirect(http.StatusFound, h.cfg.FrontendURL+"/user/dashboard")
}

// ─── Cookie helpers ───────────────────────────────────────────────────────────

// setRefreshCookie writes the refresh token into an HttpOnly cookie.
func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	if token == "" {
		return
	}
	maxAge := 604800 // default: 7 days
	if h.cfg.RefreshTokenExpiry > 0 {
		maxAge = int(h.cfg.RefreshTokenExpiry.Seconds())
	}
	c.SetCookie(h.cfg.CookieName+"_refresh", token, maxAge, "/", "", h.cfg.CookieSecure, true)
}

// clearRefreshCookie removes the refresh token cookie.
func (h *AuthHandler) clearRefreshCookie(c *gin.Context, cfg *Config) {
	c.SetCookie(cfg.CookieName+"_refresh", "", -1, "/", "", cfg.CookieSecure, true)
}
