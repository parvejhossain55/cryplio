package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"cryplio/internal/domain/identity"
	"cryplio/internal/infrastructure/storage"
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
	storage     storage.ObjectStorage
}

// NewAuthHandler creates new auth handler
func NewAuthHandler(authService identity.AuthService, cfg *Config, storage storage.ObjectStorage) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		storage:     storage,
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

// UploadAvatarHandler handles avatar upload
func (h *AuthHandler) UploadAvatarHandler(c *gin.Context) {
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

	// Log receipt of upload
	fmt.Printf("[Avatar Upload] Received upload request for user %s\n", userID.String())

	// Get the file from form data
	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Printf("[Avatar Upload] Error getting file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Validate file type
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG files are allowed"})
		return
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 2MB"})
		return
	}

	// Open the file
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileReader.Close()

	// Read file content
	fileContent, err := io.ReadAll(fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Upload to storage
	uploadInput := storage.UploadInput{
		Key:         fmt.Sprintf("avatars/%s/%s", userID.String(), file.Filename),
		ContentType: file.Header.Get("Content-Type"),
		Body:        fileContent,
	}

	uploadResult, err := h.storage.Upload(c.Request.Context(), uploadInput)
	if err != nil {
		fmt.Printf("[Avatar Upload] Storage upload error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar", "details": err.Error()})
		return
	}

	fmt.Printf("[Avatar Upload] Uploaded to storage, URL: %s\n", uploadResult.URL)

	// Update user with avatar URL
	user, err := h.authService.UpdateAvatar(c.Request.Context(), userID, uploadResult.URL)
	if err != nil {
		fmt.Printf("[Avatar Upload] UpdateAvatar error: %v\n", err)
		handleError(c, err)
		return
	}

	fmt.Printf("[Avatar Upload] Successfully updated avatar for user %s\n", userID.String())
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

// GetUserByUsernameHandler returns a public user profile by username
func (h *AuthHandler) GetUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	user, stats, err := h.authService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		handleError(c, err)
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update last seen for the viewer (if authenticated)
	if viewerIDStr, exists := c.Get("user_id"); exists {
		if viewerID, err := uuid.Parse(viewerIDStr.(string)); err == nil {
			_ = h.authService.UpdateLastSeen(c.Request.Context(), viewerID)
		}
	}

	userResp := mapUser(user)
	userResp.Stats = mapUserStats(stats)

	c.JSON(http.StatusOK, gin.H{"user": userResp})
}

// BlockUserHandler blocks another user
func (h *AuthHandler) BlockUserHandler(c *gin.Context) {
	blockerIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	blockerID, err := uuid.Parse(blockerIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		BlockedID   string  `json:"blocked_id" binding:"required"`
		Reason      *string `json:"reason,omitempty"`
		IsPermanent bool    `json:"is_permanent"`
		ExpiresAt   *string `json:"expires_at,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	blockedID, err := uuid.Parse(req.BlockedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blocked user ID"})
		return
	}
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expires_at format"})
			return
		}
		expiresAt = &t
	}

	err = h.authService.BlockUser(c.Request.Context(), blockerID, blockedID, *req.Reason, req.IsPermanent, expiresAt)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// UnblockUserHandler unblocks another user
func (h *AuthHandler) UnblockUserHandler(c *gin.Context) {
	blockerIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	blockerID, err := uuid.Parse(blockerIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	blockedIDStr := c.Param("blocked_id")
	if blockedIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blocked user ID is required"})
		return
	}
	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blocked user ID"})
		return
	}

	err = h.authService.UnblockUser(c.Request.Context(), blockerID, blockedID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}

// ListBlocksHandler lists all users blocked by the current user
func (h *AuthHandler) ListBlocksHandler(c *gin.Context) {
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

	blocks, err := h.authService.ListBlocks(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocks": blocks})
}
