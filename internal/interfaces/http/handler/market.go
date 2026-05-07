package handler

import (
	"net/http"

	"cryplio/internal/domain/market"

	"github.com/gin-gonic/gin"
)

// MarketHandler handles market rate endpoints
type MarketHandler struct {
	rateService market.RateService
}

// NewMarketHandler creates a new market handler
func NewMarketHandler(rateService market.RateService) *MarketHandler {
	return &MarketHandler{rateService: rateService}
}

// GetRatesHandler returns live market rates for a given crypto asset
func (h *MarketHandler) GetRatesHandler(c *gin.Context) {
	cryptoSymbol := c.Param("crypto")
	if cryptoSymbol == "" {
		cryptoSymbol = "USDT"
	}

	rates, err := h.rateService.GetRates(c.Request.Context(), cryptoSymbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"crypto": cryptoSymbol,
		"rates":  rates,
	})
}

// GetRateHandler returns a specific crypto/fiat rate
func (h *MarketHandler) GetRateHandler(c *gin.Context) {
	cryptoSymbol := c.Param("crypto")
	if cryptoSymbol == "" {
		cryptoSymbol = "USDT"
	}
	fiatCode := c.Param("fiat")
	if fiatCode == "" {
		fiatCode = "USD"
	}

	rate, err := h.rateService.GetRate(c.Request.Context(), cryptoSymbol, fiatCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// GetAllRatesHandler returns all market rates
func (h *MarketHandler) GetAllRatesHandler(c *gin.Context) {
	cryptoSymbol := c.Query("crypto")
	fiatCode := c.Query("fiat")

	rates, err := h.rateService.GetAllRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by crypto or fiat if specified
	if cryptoSymbol != "" || fiatCode != "" {
		filteredRates := make([]interface{}, 0)
		for _, rate := range rates {
			// This would need to be adapted based on the actual rate structure
			// For now, return all rates
			filteredRates = append(filteredRates, rate)
		}
		c.JSON(http.StatusOK, gin.H{"rates": filteredRates})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rates": rates})
}
