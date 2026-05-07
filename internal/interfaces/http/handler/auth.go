package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/trading"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
	sharedjwt "cryplio/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles auth HTTP requests
type AuthHandler struct {
	authService    identity.AuthService
	tradeService   trading.TradeService
	disputeService dispute.Service
	cfg            *Config
	storage        storage.ObjectStorage
}

// NewAuthHandler creates new auth handler
func NewAuthHandler(authService identity.AuthService, cfg *Config, storage storage.ObjectStorage) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		storage:     storage,
	}
}

// WithTradeService sets the trade service for admin dashboard
func (h *AuthHandler) WithTradeService(s trading.TradeService) *AuthHandler {
	h.tradeService = s
	return h
}

// WithDisputeService sets the dispute service for admin dashboard
func (h *AuthHandler) WithDisputeService(s dispute.Service) *AuthHandler {
	h.disputeService = s
	return h
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

// ── User Payment Method Handlers ──────────────────────────────────────────

// ListPaymentMethodsHandler returns all payment methods for the current user
func (h *AuthHandler) ListPaymentMethodsHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	methods, err := h.authService.GetPaymentMethods(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	if methods == nil {
		methods = []identity.UserPaymentMethod{}
	}
	c.JSON(http.StatusOK, gin.H{"payment_methods": methods})
}

// CreatePaymentMethodHandler adds a new payment method for the current user
func (h *AuthHandler) CreatePaymentMethodHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var pm identity.UserPaymentMethod
	if err := c.ShouldBindJSON(&pm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if pm.PaymentMethodCode == "" || pm.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_method_code and display_name are required"})
		return
	}
	pm.IsActive = true

	result, err := h.authService.AddPaymentMethod(c.Request.Context(), userID, &pm)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"payment_method": result})
}

// UpdatePaymentMethodHandler updates a payment method owned by the current user
func (h *AuthHandler) UpdatePaymentMethodHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method id"})
		return
	}

	var pm identity.UserPaymentMethod
	if err := c.ShouldBindJSON(&pm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	pm.ID = pmID
	pm.UserID = userID

	result, err := h.authService.UpdatePaymentMethod(c.Request.Context(), userID, &pm)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_method": result})
}

// DeletePaymentMethodHandler removes a payment method owned by the current user
func (h *AuthHandler) DeletePaymentMethodHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method id"})
		return
	}

	if err := h.authService.RemovePaymentMethod(c.Request.Context(), userID, pmID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment method removed"})
}

// SetDefaultPaymentMethodHandler sets a payment method as the user's default
func (h *AuthHandler) SetDefaultPaymentMethodHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method id"})
		return
	}

	if err := h.authService.SetDefaultPaymentMethod(c.Request.Context(), userID, pmID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "default payment method updated"})
}

// Admin User Management Handlers

// ListUsersHandler returns a paginated list of users (admin only)
func (h *AuthHandler) ListUsersHandler(c *gin.Context) {
	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := fmt.Sscanf(l, "%d", &limit); err == nil && parsed == 1 {
			if limit > 100 {
				limit = 100
			}
		}
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	users, err := h.authService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, mapUser(&u))
	}
	c.JSON(http.StatusOK, gin.H{"users": response, "limit": limit, "offset": offset})
}

// SuspendUserHandler suspends a user account (admin only)
func (h *AuthHandler) SuspendUserHandler(c *gin.Context) {
	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Reason   string `json:"reason" binding:"required"`
		Duration *int   `json:"duration_minutes,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var duration *time.Duration
	if req.Duration != nil && *req.Duration > 0 {
		d := time.Duration(*req.Duration) * time.Minute
		duration = &d
	}

	if err := h.authService.SuspendUser(c.Request.Context(), adminID, userID, req.Reason, duration); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user suspended successfully"})
}

// UnsuspendUserHandler lifts a user suspension (admin only)
func (h *AuthHandler) UnsuspendUserHandler(c *gin.Context) {
	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.authService.UnsuspendUser(c.Request.Context(), adminID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unsuspended successfully"})
}

// BanUserHandler bans a user account permanently (admin only)
func (h *AuthHandler) BanUserHandler(c *gin.Context) {
	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.BanUser(c.Request.Context(), adminID, userID, req.Reason); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user banned successfully"})
}

// UnbanUserHandler unbans a user account (admin only)
func (h *AuthHandler) UnbanUserHandler(c *gin.Context) {
	adminID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.authService.UnbanUser(c.Request.Context(), adminID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unbanned successfully"})
}

// GetDashboardStatsHandler returns admin dashboard aggregated stats
func (h *AuthHandler) GetDashboardStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	stats := identity.DashboardStats{}

	if h.authService != nil {
		stats.TotalUsers, _ = h.authService.CountUsers(ctx)
	}
	if h.tradeService != nil {
		stats.TotalTrades, _ = h.tradeService.CountTrades(ctx, "")
		stats.PendingTrades, _ = h.tradeService.CountTrades(ctx, "pending")
		stats.ActiveTrades, _ = h.tradeService.CountTrades(ctx, "active")
		stats.PaidTrades, _ = h.tradeService.CountTrades(ctx, "paid")
		stats.CompletedTrades, _ = h.tradeService.CountTrades(ctx, "completed")
		stats.DisputedTrades, _ = h.tradeService.CountTrades(ctx, "disputed")
		stats.CancelledTrades, _ = h.tradeService.CountTrades(ctx, "cancelled")
	}
	if h.disputeService != nil {
		stats.TotalDisputes, _ = h.disputeService.CountDisputes(ctx, "")
		stats.PendingDisputes, _ = h.disputeService.CountDisputes(ctx, "pending")
	}

	c.JSON(http.StatusOK, stats)
}
