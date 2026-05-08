package platform

import (
	"net/http"
	"strconv"

	domaplatform "cryplio/internal/domain/platform"

	"github.com/gin-gonic/gin"
)

// PlatformHandler handles platform catalogue management endpoints (admin only).
// Handler implementations are split across focused files:
//
//	crypto.go    — CryptoAsset CRUD handlers
//	fiat.go      — FiatCurrency CRUD handlers
//	catalogue.go — PaymentMethod catalogue CRUD handlers
type PlatformHandler struct {
	platformService domaplatform.PlatformService
}

// NewPlatformHandler creates a new PlatformHandler.
func NewPlatformHandler(platformService domaplatform.PlatformService) *PlatformHandler {
	return &PlatformHandler{platformService: platformService}
}

// parsePlatformPage reads page and limit query params with sensible defaults.
func parsePlatformPage(c *gin.Context) (page, limit int) {
	page, limit = 1, 50
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	return
}

// parsePlatformID parses the ":id" path param as an integer.
// Returns (0, false) and writes a 400 response on failure.
func parsePlatformID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}
