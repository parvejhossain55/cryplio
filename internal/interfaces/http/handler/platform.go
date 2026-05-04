package handler

import (
	"cryplio/internal/domain/platform"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PlatformHandler handles platform management endpoints
type PlatformHandler struct {
	platformService platform.PlatformService
}

// NewPlatformHandler creates a new platform handler
func NewPlatformHandler(platformService platform.PlatformService) *PlatformHandler {
	return &PlatformHandler{
		platformService: platformService,
	}
}

// Crypto Assets

// CreateCryptoAssetRequest represents the request to create a crypto asset
type CreateCryptoAssetRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Blockchain      string  `json:"blockchain" binding:"required"`
	ContractAddress *string `json:"contract_address,omitempty"`
	Decimals        int     `json:"decimals" binding:"required,min=0"`
	MinConfirmation int     `json:"min_confirmation" binding:"required,min=0"`
}

// CreateCryptoAssetHandler handles POST /admin/crypto-assets
func (h *PlatformHandler) CreateCryptoAssetHandler(c *gin.Context) {
	var req CreateCryptoAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	asset, err := h.platformService.CreateCryptoAsset(
		c.Request.Context(),
		req.Symbol, req.Name, req.Blockchain, req.ContractAddress, req.Decimals, req.MinConfirmation,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, asset)
}

// GetCryptoAssetsHandler handles GET /admin/crypto-assets
func (h *PlatformHandler) GetCryptoAssetsHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	page := 1
	limit := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	assets, total, err := h.platformService.GetCryptoAssets(c.Request.Context(), activeOnly, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"crypto_assets": assets,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetCryptoAssetHandler handles GET /admin/crypto-assets/:id
func (h *PlatformHandler) GetCryptoAssetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	asset, err := h.platformService.GetCryptoAsset(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, asset)
}

// UpdateCryptoAssetRequest represents the request to update a crypto asset
type UpdateCryptoAssetRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Blockchain      string  `json:"blockchain" binding:"required"`
	ContractAddress *string `json:"contract_address,omitempty"`
	Decimals        int     `json:"decimals" binding:"required,min=0"`
	MinConfirmation int     `json:"min_confirmation" binding:"required,min=0"`
	IsActive        bool    `json:"is_active"`
}

// UpdateCryptoAssetHandler handles PUT /admin/crypto-assets/:id
func (h *PlatformHandler) UpdateCryptoAssetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateCryptoAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	asset, err := h.platformService.UpdateCryptoAsset(
		c.Request.Context(),
		id, req.Symbol, req.Name, req.Blockchain, req.ContractAddress, req.Decimals, req.MinConfirmation, req.IsActive,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, asset)
}

// DeleteCryptoAssetHandler handles DELETE /admin/crypto-assets/:id
func (h *PlatformHandler) DeleteCryptoAssetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.platformService.DeleteCryptoAsset(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crypto asset deleted successfully"})
}

// Fiat Currencies

// CreateFiatCurrencyRequest represents the request to create a fiat currency
type CreateFiatCurrencyRequest struct {
	Code   string `json:"code" binding:"required,len=3"`
	Name   string `json:"name" binding:"required"`
	Symbol string `json:"symbol" binding:"required"`
}

// CreateFiatCurrencyHandler handles POST /admin/fiat-currencies
func (h *PlatformHandler) CreateFiatCurrencyHandler(c *gin.Context) {
	var req CreateFiatCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	currency, err := h.platformService.CreateFiatCurrency(c.Request.Context(), req.Code, req.Name, req.Symbol)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, currency)
}

// GetFiatCurrenciesHandler handles GET /admin/fiat-currencies
func (h *PlatformHandler) GetFiatCurrenciesHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	page := 1
	limit := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	currencies, total, err := h.platformService.GetFiatCurrencies(c.Request.Context(), activeOnly, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fiat_currencies": currencies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetFiatCurrencyHandler handles GET /admin/fiat-currencies/:id
func (h *PlatformHandler) GetFiatCurrencyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	currency, err := h.platformService.GetFiatCurrency(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, currency)
}

// UpdateFiatCurrencyRequest represents the request to update a fiat currency
type UpdateFiatCurrencyRequest struct {
	Code     string `json:"code" binding:"required,len=3"`
	Name     string `json:"name" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// UpdateFiatCurrencyHandler handles PUT /admin/fiat-currencies/:id
func (h *PlatformHandler) UpdateFiatCurrencyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateFiatCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	currency, err := h.platformService.UpdateFiatCurrency(c.Request.Context(), id, req.Code, req.Name, req.Symbol, req.IsActive)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, currency)
}

// DeleteFiatCurrencyHandler handles DELETE /admin/fiat-currencies/:id
func (h *PlatformHandler) DeleteFiatCurrencyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.platformService.DeleteFiatCurrency(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fiat currency deleted successfully"})
}

// Payment Methods

// CreatePaymentMethodRequest represents the request to create a payment method
type CreatePaymentMethodRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required,oneof=mobile_money bank_transfer online_wallet crypto cash"`
	IconURL     *string `json:"icon_url,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   int     `json:"sort_order"`
}

// CreatePaymentMethodHandler handles POST /admin/payment-methods
func (h *PlatformHandler) CreatePaymentMethodHandler(c *gin.Context) {
	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	category := platform.PaymentCategory(req.Category)
	method, err := h.platformService.CreatePaymentMethod(
		c.Request.Context(),
		req.Code, req.Name, category, req.IconURL, req.Description, req.SortOrder,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, method)
}

// GetPaymentMethodsHandler handles GET /admin/payment-methods
func (h *PlatformHandler) GetPaymentMethodsHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	page := 1
	limit := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	methods, total, err := h.platformService.GetPaymentMethods(c.Request.Context(), activeOnly, page, limit)
	if err != nil {
		handleError(c, err)
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

// GetPaymentMethodHandler handles GET /admin/payment-methods/:id
func (h *PlatformHandler) GetPaymentMethodHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	method, err := h.platformService.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, method)
}

// UpdatePaymentMethodRequest represents the request to update a payment method
type UpdatePaymentMethodRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required,oneof=mobile_money bank_transfer online_wallet crypto cash"`
	IconURL     *string `json:"icon_url,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

// UpdatePaymentMethodHandler handles PUT /admin/payment-methods/:id
func (h *PlatformHandler) UpdatePaymentMethodHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	category := platform.PaymentCategory(req.Category)
	method, err := h.platformService.UpdatePaymentMethod(
		c.Request.Context(),
		id, req.Code, req.Name, category, req.IconURL, req.Description, req.IsActive, req.SortOrder,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, method)
}

// DeletePaymentMethodHandler handles DELETE /admin/payment-methods/:id
func (h *PlatformHandler) DeletePaymentMethodHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.platformService.DeletePaymentMethod(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}
