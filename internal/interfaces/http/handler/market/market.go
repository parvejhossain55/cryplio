package market

import (
	"net/http"
	"strings"

	marketdomain "cryplio/internal/domain/market"

	"github.com/gin-gonic/gin"
)

// MarketHandler handles market rate endpoints
type MarketHandler struct {
	rateService marketdomain.RateService
}

// NewMarketHandler creates a new market handler
func NewMarketHandler(rateService marketdomain.RateService) *MarketHandler {
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

// GetAllRatesHandler returns all market rates, optionally filtered by crypto or fiat.
func (h *MarketHandler) GetAllRatesHandler(c *gin.Context) {
	cryptoFilter := strings.ToUpper(strings.TrimSpace(c.Query("crypto")))
	fiatFilter := strings.ToUpper(strings.TrimSpace(c.Query("fiat")))

	rates, err := h.rateService.GetAllRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if cryptoFilter == "" && fiatFilter == "" {
		c.JSON(http.StatusOK, gin.H{"rates": rates})
		return
	}

	filtered := make([]marketdomain.Rate, 0, len(rates))
	for _, r := range rates {
		if cryptoFilter != "" && strings.ToUpper(r.CryptoSymbol) != cryptoFilter {
			continue
		}
		if fiatFilter != "" && strings.ToUpper(r.FiatSymbol) != fiatFilter {
			continue
		}
		filtered = append(filtered, r)
	}
	c.JSON(http.StatusOK, gin.H{"rates": filtered})
}
