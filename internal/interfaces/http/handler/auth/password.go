package auth

import (
	"net/http"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

// ─── Email Verification ───────────────────────────────────────────────────────

// RequestEmailVerificationHandler sends a verification link to the user's email.
func (h *AuthHandler) RequestEmailVerificationHandler(c *gin.Context) {
	var req dto.EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := h.authService.RequestEmailVerification(c.Request.Context(), req.UserID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// VerifyEmailHandler confirms the user's email address using the token from the link.
func (h *AuthHandler) VerifyEmailHandler(c *gin.Context) {
	var req dto.EmailVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	user, err := h.authService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully", "user": mapUser(user)})
}

// ─── Password Reset ───────────────────────────────────────────────────────────

// RequestPasswordResetHandler sends a password-reset link to the email if it
// exists. Always responds with the same message to prevent user enumeration.
func (h *AuthHandler) RequestPasswordResetHandler(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Ignore errors to avoid leaking whether the email exists.
	_ = h.authService.RequestPasswordReset(c.Request.Context(), req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset link has been sent"})
}

// ResetPasswordHandler sets a new password using the token from the reset link.
func (h *AuthHandler) ResetPasswordHandler(c *gin.Context) {
	var req dto.PasswordResetConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	user, err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully", "user": mapUser(user)})
}
