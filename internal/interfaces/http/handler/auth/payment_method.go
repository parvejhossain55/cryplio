package auth

import (
	"net/http"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/domain/identity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ─── User Payment Methods ─────────────────────────────────────────────────────
//
// These endpoints manage the payment method profiles a user can attach to their
// account (e.g. bKash account, bank transfer details). They are distinct from
// the platform-level payment method catalogue managed by admins in platform.go.

// ListPaymentMethodsHandler returns all payment methods saved by the current user.
func (h *AuthHandler) ListPaymentMethodsHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	methods, err := h.paymentManager.GetPaymentMethods(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}
	if methods == nil {
		methods = []identity.UserPaymentMethod{}
	}

	c.JSON(http.StatusOK, gin.H{"payment_methods": methods})
}

// CreatePaymentMethodHandler adds a new payment method for the current user.
func (h *AuthHandler) CreatePaymentMethodHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var pm identity.UserPaymentMethod
	if err := c.ShouldBindJSON(&pm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if pm.PaymentMethodCode == "" || pm.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_method_code and display_name are required"})
		return
	}
	pm.IsActive = true

	result, err := h.paymentManager.AddPaymentMethod(c.Request.Context(), userID, &pm)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payment_method": result})
}

// UpdatePaymentMethodHandler updates a payment method owned by the current user.
func (h *AuthHandler) UpdatePaymentMethodHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}

	var pm identity.UserPaymentMethod
	if err := c.ShouldBindJSON(&pm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	pm.ID = pmID
	pm.UserID = userID

	result, err := h.paymentManager.UpdatePaymentMethod(c.Request.Context(), userID, &pm)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_method": result})
}

// DeletePaymentMethodHandler removes a payment method owned by the current user.
func (h *AuthHandler) DeletePaymentMethodHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}

	if err := h.paymentManager.RemovePaymentMethod(c.Request.Context(), userID, pmID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method removed"})
}

// SetDefaultPaymentMethodHandler marks one of the user's payment methods as
// their default for new trades.
func (h *AuthHandler) SetDefaultPaymentMethodHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	pmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}

	if err := h.paymentManager.SetDefaultPaymentMethod(c.Request.Context(), userID, pmID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default payment method updated"})
}
