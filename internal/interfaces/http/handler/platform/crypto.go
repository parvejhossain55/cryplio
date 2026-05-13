package platform

import (
	"net/http"

	domaplatform "cryplio/internal/domain/platform"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
)

// CreateCryptoAssetRequest represents the request to create a crypto asset.
type CreateCryptoAssetRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Blockchain      string  `json:"blockchain" binding:"required"`
	ContractAddress *string `json:"contract_address,omitempty"`
	Decimals        int     `json:"decimals" binding:"required,min=0"`
}

// UpdateCryptoAssetRequest represents the request to update a crypto asset.
type UpdateCryptoAssetRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Blockchain      string  `json:"blockchain" binding:"required"`
	ContractAddress *string `json:"contract_address,omitempty"`
	Decimals        int     `json:"decimals" binding:"required,min=0"`
	IsActive        bool    `json:"is_active"`
}

// CreateCryptoAssetHandler handles POST /admin/crypto-assets.
func (h *PlatformHandler) CreateCryptoAssetHandler(c *gin.Context) {
	var req CreateCryptoAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := domaplatform.CreateCryptoAssetInput{
		Symbol:          req.Symbol,
		Name:            req.Name,
		Blockchain:      req.Blockchain,
		ContractAddress: req.ContractAddress,
		Decimals:        req.Decimals,
	}
	asset, err := h.platformService.CreateCryptoAsset(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, asset)
}

// GetCryptoAssetsHandler handles GET /admin/crypto-assets.
func (h *PlatformHandler) GetCryptoAssetsHandler(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	page, limit := parsePlatformPage(c)

	assets, total, err := h.platformService.GetCryptoAssets(c.Request.Context(), activeOnly, page, limit)
	if err != nil {
		basehandler.HandleError(c, err)
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

// GetCryptoAssetHandler handles GET /admin/crypto-assets/:id.
func (h *PlatformHandler) GetCryptoAssetHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	asset, err := h.platformService.GetCryptoAsset(c.Request.Context(), id)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, asset)
}

// UpdateCryptoAssetHandler handles PUT /admin/crypto-assets/:id.
func (h *PlatformHandler) UpdateCryptoAssetHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	var req UpdateCryptoAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	input := domaplatform.UpdateCryptoAssetInput{
		Symbol:          req.Symbol,
		Name:            req.Name,
		Blockchain:      req.Blockchain,
		ContractAddress: req.ContractAddress,
		Decimals:        req.Decimals,
		IsActive:        req.IsActive,
	}
	asset, err := h.platformService.UpdateCryptoAsset(c.Request.Context(), id, input)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, asset)
}

// DeleteCryptoAssetHandler handles DELETE /admin/crypto-assets/:id.
func (h *PlatformHandler) DeleteCryptoAssetHandler(c *gin.Context) {
	id, ok := parsePlatformID(c)
	if !ok {
		return
	}

	if err := h.platformService.DeleteCryptoAsset(c.Request.Context(), id); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crypto asset deleted successfully"})
}
