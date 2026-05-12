package trade

// ad.go contains TradeHandler methods for trade advertisement management.

import (
	"net/http"
	"strconv"
	"time"

	"cryplio/internal/domain/trading"
	"cryplio/internal/interfaces/http/dto"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListAdsHandler returns a paginated, filtered list of active trade ads.
func (h *TradeHandler) ListAdsHandler(c *gin.Context) {
	adType := trading.AdType(c.Query("type"))
	cryptoID, _ := strconv.Atoi(c.Query("crypto_id"))
	fiatID, _ := strconv.Atoi(c.Query("fiat_id"))
	fiatCurrency := c.Query("fiat_currency")
	paymentMethodStr := c.Query("payment_method")
	sortBy := c.DefaultQuery("sort_by", "newest")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := trading.AdFilter{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
	}

	if adType != "" {
		filter.Type = &adType
	}
	if cryptoID > 0 {
		filter.CryptoID = &cryptoID
	}
	if fiatID > 0 {
		filter.FiatID = &fiatID
	}
	if fiatCurrency != "" && fiatCurrency != "all" {
		filter.FiatCode = &fiatCurrency
	}
	if paymentMethodStr != "" && paymentMethodStr != "all" {
		// Try to parse as int first, otherwise will need payment method lookup
		if pmID, err := strconv.Atoi(paymentMethodStr); err == nil && pmID > 0 {
			filter.PaymentMethods = []int{pmID}
		}
	}

	ads, total, err := h.adManager.ListActiveAds(c.Request.Context(), filter)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	response := dto.ListAdsResponse{
		Total: total,
		Ads:   make([]dto.AdResponse, len(ads)),
	}

	for i, ad := range ads {
		isOnline := ad.UserLastSeen != nil && ad.UserLastSeen.After(time.Now().Add(-5*time.Minute))
		response.Ads[i] = dto.AdResponse{
			AdID:                 ad.AdID.String(),
			UserID:               ad.UserID.String(),
			Username:             ad.Username,
			UserAvatar:           ad.UserAvatar,
			UserRating:           ad.UserRating,
			UserTrades:           ad.UserTrades,
			Type:                 string(ad.Type),
			CryptoSymbol:         ad.CryptoSymbol,
			FiatSymbol:           ad.FiatSymbol,
			PriceType:            string(ad.PriceType),
			Price:                ad.Price,
			MinAmount:            ad.MinAmount,
			MaxAmount:            ad.MaxAmount,
			PaymentMethods:       ad.PaymentMethodNames,
			PaymentWindowMinutes: ad.PaymentWindowMinutes,
			IsOnline:             isOnline,
			TradeTerms: func() string {
				if ad.TradeTerms != nil {
					return *ad.TradeTerms
				}
				return ""
			}(),
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateAdHandler creates a new trade advertisement for the authenticated user.
func (h *TradeHandler) CreateAdHandler(c *gin.Context) {
	var req dto.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	ad := &trading.TradeAd{
		AdID:                 uuid.New(),
		UserID:               userID,
		Type:                 trading.AdType(req.Type),
		CryptoID:             req.CryptoID,
		FiatID:               req.FiatID,
		PriceType:            trading.PriceType(req.PriceType),
		Price:                req.Price,
		FloatingMarkup:       req.FloatingMarkup,
		MinAmount:            req.MinAmount,
		MaxAmount:            req.MaxAmount,
		PaymentMethods:       req.PaymentMethods,
		TradeTerms:           &req.TradeTerms,
		PaymentWindowMinutes: req.PaymentWindowMinutes,
		IsPublic:             true,
		IsPaused:             false,
		Timezone:             "UTC",
		Status:               trading.TradeAdStatusActive,
	}

	if err := h.adManager.CreateAd(c.Request.Context(), ad); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ad created successfully", "ad_id": ad.AdID})
}

// UpdateAdHandler updates an existing trade advertisement.
func (h *TradeHandler) UpdateAdHandler(c *gin.Context) {
	adID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	updates := &trading.TradeAd{
		Type:                 trading.AdType(req.Type),
		CryptoID:             req.CryptoID,
		FiatID:               req.FiatID,
		PriceType:            trading.PriceType(req.PriceType),
		Price:                req.Price,
		FloatingMarkup:       req.FloatingMarkup,
		MinAmount:            req.MinAmount,
		MaxAmount:            req.MaxAmount,
		PaymentMethods:       req.PaymentMethods,
		TradeTerms:           &req.TradeTerms,
		PaymentWindowMinutes: req.PaymentWindowMinutes,
		Timezone:             req.Timezone,
	}

	if err := h.adManager.UpdateAd(c.Request.Context(), adID, userID, updates); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad updated successfully"})
}

// DeleteAdHandler removes a trade advertisement.
func (h *TradeHandler) DeleteAdHandler(c *gin.Context) {
	adID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	if err := h.adManager.DeleteAd(c.Request.Context(), adID, userID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad deleted successfully"})
}

// ListUserAdsHandler returns all advertisements belonging to the authenticated user.
func (h *TradeHandler) ListUserAdsHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	ads, total, err := h.adManager.ListUserAds(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	// Map to DTOs
	response := dto.ListAdsResponse{
		Total: total,
		Ads:   make([]dto.AdResponse, len(ads)),
	}

	for i, ad := range ads {
		response.Ads[i] = dto.AdResponse{
			AdID:                 ad.AdID.String(),
			UserID:               ad.UserID.String(),
			Type:                 string(ad.Type),
			PriceType:            string(ad.PriceType),
			Price:                ad.Price,
			MinAmount:            ad.MinAmount,
			MaxAmount:            ad.MaxAmount,
			PaymentMethods:       ad.PaymentMethodNames,
			PaymentWindowMinutes: ad.PaymentWindowMinutes,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ToggleAdStatusHandler pauses or activates a trade advertisement.
func (h *TradeHandler) ToggleAdStatusHandler(c *gin.Context) {
	adID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	if err := h.adManager.ToggleAdStatus(c.Request.Context(), adID, userID); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad status updated successfully"})
}
