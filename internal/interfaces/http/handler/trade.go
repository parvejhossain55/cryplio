package handler

import (
	"net/http"
	"strconv"

	"cryplio/internal/domain/trading"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TradeHandler struct {
	tradeService trading.TradeService
}

func NewTradeHandler(service trading.TradeService) *TradeHandler {
	return &TradeHandler{tradeService: service}
}

func (h *TradeHandler) ListAdsHandler(c *gin.Context) {
	adType := trading.AdType(c.Query("type"))
	cryptoID, _ := strconv.Atoi(c.Query("crypto_id"))
	fiatID, _ := strconv.Atoi(c.Query("fiat_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := trading.AdFilter{
		Limit:  limit,
		Offset: offset,
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

	ads, total, err := h.tradeService.ListActiveAds(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	response := dto.ListAdsResponse{
		Total: total,
		Ads:   make([]dto.AdResponse, len(ads)),
	}

	for i, ad := range ads {
		var avatar string
		if ad.User.AvatarURL != nil {
			avatar = *ad.User.AvatarURL
		}
		response.Ads[i] = dto.AdResponse{
			AdID:                 ad.AdID.String(),
			UserID:               ad.UserID.String(),
			Username:             ad.User.Username,
			UserAvatar:           avatar,
			UserRating:           0, // Default for now
			UserTrades:           0, // Default for now
			Type:                 string(ad.Type),
			PriceType:            string(ad.PriceType),
			Price:                ad.Price,
			MinAmount:            ad.MinAmount,
			MaxAmount:            ad.MaxAmount,
			PaymentWindowMinutes: ad.PaymentWindowMinutes,
			IsOnline:             ad.User.IsOnline(),
		}
		if ad.Stats != nil {
			response.Ads[i].UserTrades = ad.Stats.TotalTrades
			if ad.Stats.AvgRating != nil {
				response.Ads[i].UserRating = *ad.Stats.AvgRating
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *TradeHandler) InitiateTradeHandler(c *gin.Context) {
	adID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get buyer ID from context (set by AuthMiddleware)
	buyerIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	buyerID, _ := uuid.Parse(buyerIDStr.(string))

	trade, err := h.tradeService.InitiateTrade(c.Request.Context(), adID, buyerID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Trade initiated successfully",
		"trade_id": trade.TradeID.String(),
		"status":   trade.Status,
	})
}
