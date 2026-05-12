package auth

import (
	"fmt"
	"io"
	"net/http"

	"cryplio/internal/domain/identity"
	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
	"cryplio/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ─── Own Profile ─────────────────────────────────────────────────────────────

// GetUserHandler returns the authenticated user's full profile including stats and header fields.
func (h *AuthHandler) GetUserHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	user, err := h.profileManager.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	stats, _ := h.profileManager.GetUserStats(c.Request.Context(), userID)

	// Get unread notification count
	unreadCount, _ := h.notificationService.GetUnreadCount(c.Request.Context(), userID)

	// Calculate trader badge and security status
	traderBadge := h.determineTraderBadge(user, stats)
	accountHealth, accountSecurity, twoFactorStatus, loginNotifications := h.calculateSecurityStatus(user)

	// Map user with stats and header fields
	userDTO := mapUserWithStats(user, stats)
	userDTO.TraderBadge = traderBadge
	userDTO.UnreadNotificationCount = unreadCount
	userDTO.AccountHealth = accountHealth
	userDTO.AccountSecurity = accountSecurity
	userDTO.TwoFactorStatus = twoFactorStatus
	userDTO.LoginNotifications = loginNotifications

	c.JSON(http.StatusOK, gin.H{"user": userDTO})
}

// UpdateUserHandler updates the authenticated user's profile fields.
func (h *AuthHandler) UpdateUserHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	httpvalidator.NormalizeUpdateProfileRequest(&req)

	user, err := h.profileManager.UpdateProfile(c.Request.Context(), userID, req.Username, req.Bio)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUser(user)})
}

// ─── Avatar ───────────────────────────────────────────────────────────────────

// UploadAvatarHandler accepts a multipart JPEG or PNG file (≤ 2 MB) and stores
// it in object storage, updating the user's avatar URL on success.
func (h *AuthHandler) UploadAvatarHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	logger.Info("avatar upload received", logger.Fields{"user_id": userID.String()})

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG files are allowed"})
		return
	}

	const maxSize = 2 * 1024 * 1024 // 2 MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 2 MB"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileReader.Close()

	fileContent, err := io.ReadAll(fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	uploadResult, err := h.storage.Upload(c.Request.Context(), storage.UploadInput{
		Key:         fmt.Sprintf("avatars/%s/%s", userID.String(), file.Filename),
		ContentType: contentType,
		Body:        fileContent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar", "details": err.Error()})
		return
	}

	user, err := h.profileManager.UpdateAvatar(c.Request.Context(), userID, uploadResult.URL)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUser(user)})
}

// ─── Public Profile ───────────────────────────────────────────────────────────

// GetUserByUsernameHandler returns a user's public profile including trade stats.
// Also refreshes last-seen for the authenticated viewer (if any).
func (h *AuthHandler) GetUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	user, stats, err := h.profileManager.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Refresh last-seen for the authenticated viewer (best-effort).
	if viewerIDRaw, exists := c.Get("user_id"); exists {
		if viewerID, err := uuid.Parse(viewerIDRaw.(string)); err == nil {
			_ = h.profileManager.UpdateLastSeen(c.Request.Context(), viewerID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUserWithStats(user, stats)})
}

// ─── Header Profile ───────────────────────────────────────────────────────────

// GetHeaderProfileHandler returns the user profile data needed for the header component.
// This includes username, avatar, online status, trader badge, unread notification count, and security status.
func (h *AuthHandler) GetHeaderProfileHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	user, err := h.profileManager.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	// Get unread notification count
	unreadCount, _ := h.notificationService.GetUnreadCount(c.Request.Context(), userID)

	// Determine trader badge based on user stats and role
	stats, _ := h.profileManager.GetUserStats(c.Request.Context(), userID)
	traderBadge := h.determineTraderBadge(user, stats)

	// Calculate account health and security status
	accountHealth, accountSecurity, twoFactorStatus, loginNotifications := h.calculateSecurityStatus(user)

	response := dto.HeaderProfileResponse{
		Username:                user.Username,
		AvatarURL:               user.AvatarURL,
		IsOnline:                user.IsOnline(),
		TraderBadge:             traderBadge,
		UnreadNotificationCount: unreadCount,
		AccountHealth:           accountHealth,
		AccountSecurity:         accountSecurity,
		TwoFactorStatus:         twoFactorStatus,
		LoginNotifications:      loginNotifications,
	}

	c.JSON(http.StatusOK, response)
}

// determineTraderBadge determines the trader badge based on user role and stats
func (h *AuthHandler) determineTraderBadge(user *identity.User, stats *identity.UserStats) string {
	// Admin badge
	if user.Role == identity.UserRoleAdmin {
		return "ADMIN"
	}

	// Merchant badge
	if user.IsMerchant || user.Role == identity.UserRoleMerchant {
		return "PRO TRADER"
	}

	// Verified trader badge based on successful trades
	if stats != nil && stats.SuccessfulTrades >= 10 {
		return "VERIFIED"
	}

	// New trader badge
	if stats != nil && stats.SuccessfulTrades >= 5 {
		return "TRADER"
	}

	return ""
}

// calculateSecurityStatus calculates the account health and security status
func (h *AuthHandler) calculateSecurityStatus(user *identity.User) (accountHealth, accountSecurity, twoFactorStatus, loginNotifications string) {
	// Account Security: based on email verification
	if user.EmailVerified {
		accountSecurity = "VERIFIED"
	} else {
		accountSecurity = "UNVERIFIED"
	}

	// 2FA Status: based on two factor authentication (check if secret is set)
	if user.TwoFASecret != nil && *user.TwoFASecret != "" {
		twoFactorStatus = "ENABLED"
	} else {
		twoFactorStatus = "DISABLED"
	}

	// Login Notifications: assume active for now (can be enhanced with user preferences)
	loginNotifications = "ACTIVE"

	// Calculate overall account health based on security factors
	securityScore := 0
	if user.EmailVerified {
		securityScore++
	}
	if user.TwoFASecret != nil && *user.TwoFASecret != "" {
		securityScore++
	}

	// Determine health based on score
	switch securityScore {
	case 2:
		accountHealth = "EXCELLENT"
	case 1:
		accountHealth = "GOOD"
	default:
		accountHealth = "FAIR"
	}

	return
}
