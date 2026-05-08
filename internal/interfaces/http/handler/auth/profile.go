package auth

import (
	"fmt"
	"io"
	"net/http"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	httpvalidator "cryplio/internal/interfaces/http/validator"
	"cryplio/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ─── Own Profile ─────────────────────────────────────────────────────────────

// GetUserHandler returns the authenticated user's full profile including stats.
func (h *AuthHandler) GetUserHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	stats, _ := h.authService.GetUserStats(c.Request.Context(), userID)
	c.JSON(http.StatusOK, gin.H{"user": mapUserWithStats(user, stats)})
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

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, req.Username, req.Bio)
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

	user, err := h.authService.UpdateAvatar(c.Request.Context(), userID, uploadResult.URL)
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

	user, stats, err := h.authService.GetUserByUsername(c.Request.Context(), username)
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
			_ = h.authService.UpdateLastSeen(c.Request.Context(), viewerID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": mapUserWithStats(user, stats)})
}
