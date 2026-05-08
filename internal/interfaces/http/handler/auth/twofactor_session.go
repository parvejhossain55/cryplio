package auth

import (
	"net/http"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

// ─── Two-Factor Authentication Setup ─────────────────────────────────────────

// Setup2FAHandler generates a TOTP secret and provisioning URI.
// The client should display the QR code to the user and then call
// Verify2FAHandler to confirm the setup.
func (h *AuthHandler) Setup2FAHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	secret, uri, err := h.authService.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TwoFactorSetupResponse{Secret: secret, ProvisioningURI: uri})
}

// Verify2FAHandler confirms the pending 2FA setup with a valid TOTP code,
// enabling 2FA on the account.
func (h *AuthHandler) Verify2FAHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := h.authService.Verify2FA(c.Request.Context(), userID, req.Code); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Two-factor authentication enabled"})
}

// Disable2FAHandler turns off 2FA after the user confirms their password.
func (h *AuthHandler) Disable2FAHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := h.authService.Disable2FA(c.Request.Context(), userID, req.Password); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Two-factor authentication disabled"})
}

// ─── Session Management ───────────────────────────────────────────────────────

// GetSessionsHandler returns all active login sessions for the current user.
func (h *AuthHandler) GetSessionsHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	sessions, err := h.authService.GetSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// DeleteSessionHandler revokes a specific session identified by its token ID.
// Use this to force-sign-out from a particular device.
func (h *AuthHandler) DeleteSessionHandler(c *gin.Context) {
	// Verify the caller is authenticated (we don't use the userID directly, but
	// we require it to ensure the middleware ran correctly).
	_, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	tokenID := c.Param("tokenId")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token ID is required"})
		return
	}

	if err := h.authService.DeleteSession(c.Request.Context(), tokenID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}
