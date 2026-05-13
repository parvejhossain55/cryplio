package platform

import (
	"net/http"

	domaplatform "cryplio/internal/domain/platform"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
)

// CreateFiatCurrencyRequest represents the request to create a fiat currency.
type CreateFiatCurrencyRequest struct {
	Code   string `json:"code" binding:"required,len=3"`
	Name   string `json:"name" binding:"required"`
	Symbol string `json:"symbol" binding:"required"`
}

// UpdateFiatCurrencyRequest represents the request to update a fiat currency.
type UpdateFiatCurrencyRequest struct {
	Code     string `json:"code" binding:"required,len=3"`
	Name     string `json:"name" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// CreateFiatCurrencyHandler handles POST /admin/fiat-currencies.
func (h *PlatformHandler) CreateFiatCurrencyHandler(c *gin.Context) {
	var req CreateFiatCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	input := domaplatform.CreateFiatCurrencyInput{
		Code:   req.Code,
		Name:   req.Name,
		Symbol: req.Symbol,
	}
	currency, err := h.platformService.CreateFiatCurrency(c.Request.Context(), input)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, currency)
}

// GetFiatCurrenciesHandler handles GET /admin/fiat-currencies.
func (h *PlatformHandler) GetFiatCurrenciesHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	searchQuery := c.Query("search")
	page, limit := parsePlatformPage(c)

	currencies, total, err := h.platformService.GetFiatCurrencies(c.Request.Context(), activeOnly, searchQuery, page, limit)
	if err != nil {
		basehandler.HandleError(c, err)
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

// GetFiatCurrencyHandler handles GET /admin/fiat-currencies/:id.
func (h *PlatformHandler) GetFiatCurrencyHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	currency, err := h.platformService.GetFiatCurrency(c.Request.Context(), id)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, currency)
}

// UpdateFiatCurrencyHandler handles PUT /admin/fiat-currencies/:id.
func (h *PlatformHandler) UpdateFiatCurrencyHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	var req UpdateFiatCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	input := domaplatform.UpdateFiatCurrencyInput{
		Code:     req.Code,
		Name:     req.Name,
		Symbol:   req.Symbol,
		IsActive: req.IsActive,
	}
	currency, err := h.platformService.UpdateFiatCurrency(c.Request.Context(), id, input)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, currency)
}

// DeleteFiatCurrencyHandler handles DELETE /admin/fiat-currencies/:id.
func (h *PlatformHandler) DeleteFiatCurrencyHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	if err := h.platformService.DeleteFiatCurrency(c.Request.Context(), id); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fiat currency deleted successfully"})
}
