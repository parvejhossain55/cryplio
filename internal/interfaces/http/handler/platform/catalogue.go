package platform

import (
	"net/http"

	domaplatform "cryplio/internal/domain/platform"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
)

// CreatePaymentMethodRequest represents the request to create a platform payment method catalogue entry.
type CreatePaymentMethodRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required,oneof=mobile_money bank_transfer online_wallet crypto cash"`
	IconURL     *string `json:"icon_url,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   int     `json:"sort_order"`
}

// UpdatePaymentMethodRequest represents the request to update a platform payment method catalogue entry.
type UpdatePaymentMethodRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required,oneof=mobile_money bank_transfer online_wallet crypto cash"`
	IconURL     *string `json:"icon_url,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

// CreatePaymentMethodHandler handles POST /admin/payment-methods.
func (h *PlatformHandler) CreatePaymentMethodHandler(c *gin.Context) {
	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	input := domaplatform.CreatePaymentMethodInput{
		Code:        req.Code,
		Name:        req.Name,
		Category:    domaplatform.PaymentCategory(req.Category),
		IconURL:     req.IconURL,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	method, err := h.platformService.CreatePaymentMethod(c.Request.Context(), input)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, method)
}

// GetPaymentMethodsHandler handles GET /admin/payment-methods.
func (h *PlatformHandler) GetPaymentMethodsHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	searchQuery := c.Query("search")
	page, limit := parsePlatformPage(c)

	methods, total, err := h.platformService.GetPaymentMethods(c.Request.Context(), activeOnly, searchQuery, page, limit)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_methods": methods,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetPaymentMethodHandler handles GET /admin/payment-methods/:id.
func (h *PlatformHandler) GetPaymentMethodHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	method, err := h.platformService.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, method)
}

// UpdatePaymentMethodHandler handles PUT /admin/payment-methods/:id.
func (h *PlatformHandler) UpdatePaymentMethodHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	input := domaplatform.UpdatePaymentMethodInput{
		Code:        req.Code,
		Name:        req.Name,
		Category:    domaplatform.PaymentCategory(req.Category),
		IconURL:     req.IconURL,
		Description: req.Description,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}
	method, err := h.platformService.UpdatePaymentMethod(c.Request.Context(), id, input)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, method)
}

// DeletePaymentMethodHandler handles DELETE /admin/payment-methods/:id.
func (h *PlatformHandler) DeletePaymentMethodHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	if err := h.platformService.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}
