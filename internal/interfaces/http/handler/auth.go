package handler

import (
	"context"
	"net/http"

	"cryplio/internal/domain/identity"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
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

	// Login to get token
	token, err := h.loginAfterRegister(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{Token: token, User: mapUser(user)})
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeLoginRequest(&req)

	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	setAuthCookie(c, h.cfg, token)

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token, User: mapUser(user)})
}

// LogoutHandler handles logout
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	_ = h.authService.Logout(c.Request.Context())
	clearAuthCookie(c, h.cfg)
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

// loginAfterRegister helper
func (h *AuthHandler) loginAfterRegister(ctx context.Context, email, password string) (string, error) {
	token, _, err := h.authService.Login(ctx, email, password)
	return token, err
}
