package handler

import (
	"context"
	"errors"
	"net/http"

	"cryplio/internal/domain/identity"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
	sharedjwt "cryplio/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles auth HTTP requests
type AuthHandler struct {
	authService identity.AuthService
	cfg         *Config
}

// NewAuthHandler creates new auth handler
func NewAuthHandler(authService identity.AuthService, cfg *Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeRegisterRequest(&req)

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	if err := h.authService.RequestEmailVerification(c.Request.Context(), user.UserID); err != nil {
		handleError(c, err)
		return
	}

	// Login to get tokens
	access, refresh, _, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)

	c.JSON(http.StatusCreated, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeLoginRequest(&req)

	access, refresh, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Check for 2FA required error
		var twoFactorErr identity.TwoFactorRequiredError
		if errors.As(err, &twoFactorErr) {
			c.JSON(http.StatusOK, gin.H{
				"requires_2fa": true,
				"temp_token":   twoFactorErr.TempToken,
				"user":         mapUser(twoFactorErr.User),
			})
			return
		}
		handleError(c, err)
		return
	}

	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)

	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// LogoutHandler handles logout
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	// Get refresh token from cookie to identify session
	refreshToken, _ := c.Cookie(h.cfg.CookieName + "_refresh")
	if refreshToken != "" {
		// Parse to get jti
		claims, err := sharedjwt.Parse(h.cfg.JWTSecret, refreshToken)
		if err == nil {
			if jti, ok := claims["jti"].(string); ok {
				_ = h.authService.Logout(c.Request.Context(), jti)
			}
		}
	}
	clearAuthCookie(c, h.cfg)
	h.clearRefreshCookie(c, h.cfg)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetUserHandler returns current user
func (h *AuthHandler) GetUserHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUser(user)})
}

// UpdateUserHandler updates user profile
func (h *AuthHandler) UpdateUserHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeUpdateProfileRequest(&req)

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, req.Username, req.Bio)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUser(user)})
}

// GoogleAuthHandler initiates Google OAuth
func (h *AuthHandler) GoogleAuthHandler(c *gin.Context) {
	url := h.authService.GoogleOAuthURL()
	if url == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "OAuth not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallbackHandler handles Google OAuth callback
func (h *AuthHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	token, _, err := h.authService.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		handleError(c, err)
		return
	}
	setAuthCookie(c, h.cfg, token)
	c.Redirect(http.StatusFound, h.cfg.FrontendURL+"/user/dashboard")
}

// RequestEmailVerificationHandler requests email verification
func (h *AuthHandler) RequestEmailVerificationHandler(c *gin.Context) {
	var req dto.EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	if err := h.authService.RequestEmailVerification(c.Request.Context(), req.UserID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// VerifyEmailHandler verifies email with token
func (h *AuthHandler) VerifyEmailHandler(c *gin.Context) {
	var req dto.EmailVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	user, err := h.authService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": mapUser(user), "message": "Email verified"})
}

// RequestPasswordResetHandler requests password reset
func (h *AuthHandler) RequestPasswordResetHandler(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an account exists, a reset link was sent"})
}

// ResetPasswordHandler resets password using token
func (h *AuthHandler) ResetPasswordHandler(c *gin.Context) {
	var req dto.PasswordResetConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	user, err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": mapUser(user), "message": "Password reset successful"})
}

// RefreshTokenHandler rotates refresh token
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie(h.cfg.CookieName + "_refresh")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing refresh token"})
		return
	}

	access, refresh, user, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		handleError(c, err)
		return
	}
	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// Setup2FAHandler initiates 2FA setup
func (h *AuthHandler) Setup2FAHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	secret, uri, err := h.authService.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.TwoFactorSetupResponse{Secret: secret, ProvisioningURI: uri})
}

// Verify2FAHandler confirms 2FA setup with code
func (h *AuthHandler) Verify2FAHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req dto.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	err = h.authService.Verify2FA(c.Request.Context(), userID, req.Code)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled"})
}

// Disable2FAHandler disables 2FA after password confirmation
func (h *AuthHandler) Disable2FAHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req dto.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	err = h.authService.Disable2FA(c.Request.Context(), userID, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}

// Complete2FALoginHandler finishes the 2FA login step
func (h *AuthHandler) Complete2FALoginHandler(c *gin.Context) {
	var req dto.TwoFactorLoginCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	access, refresh, user, err := h.authService.Complete2FALogin(c.Request.Context(), req.TempToken, req.Code)
	if err != nil {
		handleError(c, err)
		return
	}
	setAuthCookie(c, h.cfg, access)
	h.setRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: access, RefreshToken: refresh, User: mapUser(user)})
}

// GetSessionsHandler lists active sessions for current user
func (h *AuthHandler) GetSessionsHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	sessions, err := h.authService.GetSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// DeleteSessionHandler revokes a specific session
func (h *AuthHandler) DeleteSessionHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	_, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	tokenID := c.Param("tokenId")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token ID required"})
		return
	}
	err = h.authService.DeleteSession(c.Request.Context(), tokenID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// ============================================================
// Cookie helpers
// ============================================================

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	if token == "" {
		return
	}
	refreshCookieName := h.cfg.CookieName + "_refresh"
	// Use refresh token expiry to set max age; approximate 7 days
	maxAge := 604800 // 7 days
	if h.cfg.RefreshTokenExpiry > 0 {
		maxAge = int(h.cfg.RefreshTokenExpiry.Seconds())
	}
	c.SetCookie(refreshCookieName, token, maxAge, "/", "", h.cfg.CookieSecure, true)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context, cfg *Config) {
	refreshCookieName := cfg.CookieName + "_refresh"
	c.SetCookie(refreshCookieName, "", -1, "/", "", cfg.CookieSecure, true)
}

// loginAfterRegister returns access token after registration (login)
func (h *AuthHandler) loginAfterRegister(ctx context.Context, email, password string) (string, string, error) {
	access, refresh, _, err := h.authService.Login(ctx, email, password)
	return access, refresh, err
}
