package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		response.Ads[i] = dto.AdResponse{
			AdID:                 ad.AdID.String(),
			UserID:               ad.UserID.String(),
			Username:             ad.Username,
			UserAvatar:           ad.UserAvatar,
			UserRating:           ad.UserRating,
			UserTrades:           ad.UserTrades,
			Type:                 string(ad.Type),
			PriceType:            string(ad.PriceType),
			Price:                ad.Price,
			MinAmount:            ad.MinAmount,
			MaxAmount:            ad.MaxAmount,
			PaymentWindowMinutes: ad.PaymentWindowMinutes,
			IsOnline:             ad.UserLastSeen != nil && ad.UserLastSeen.After(time.Now().Add(-5*time.Minute)),
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

	var req dto.InitiateTradeRequest
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

func (h *TradeHandler) CreateAdHandler(c *gin.Context) {
	var req dto.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	paymentMethods := make([]int, 0, len(req.PaymentMethods))
	for _, pm := range req.PaymentMethods {
		id, err := strconv.Atoi(pm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid payment method ID: %s", pm)})
			return
		}
		paymentMethods = append(paymentMethods, id)
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
		PaymentMethods:       paymentMethods,
		TradeTerms:           &req.TradeTerms,
		PaymentWindowMinutes: req.PaymentWindowMinutes,
		IsPublic:             true,
		IsPaused:             false,
		Timezone:             "UTC",
		Status:               trading.TradeAdStatusActive,
	}

	if err := h.tradeService.CreateAd(c.Request.Context(), ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ad)
}

func (h *TradeHandler) ListMyAdsHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	ads, total, err := h.tradeService.ListUserAds(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			PaymentMethods:       []string{}, // TODO: map IDs to names
			PaymentWindowMinutes: ad.PaymentWindowMinutes,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *TradeHandler) ToggleAdStatusHandler(c *gin.Context) {
	adID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.tradeService.ToggleAdStatus(c.Request.Context(), adID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad status updated successfully"})
}

func (h *TradeHandler) ListTradesHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	role := c.Query("role") // buyer, seller, or empty for all

	// Need to add ListTrades to TradeService first
	// For now, call repo directly if service not updated yet or update service
	trades, err := h.tradeService.ListTrades(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trades)
}

func (h *TradeHandler) GetTradeHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	trade, err := h.tradeService.GetTrade(c.Request.Context(), tradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trade not found"})
		return
	}

	c.JSON(http.StatusOK, trade)
}

func (h *TradeHandler) UpdateTradeStatusHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required,oneof=pay release cancel"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var updateErr error
	switch req.Action {
	case "pay":
		updateErr = h.tradeService.MarkAsPaid(c.Request.Context(), tradeID, userID)
	case "release":
		updateErr = h.tradeService.ReleaseEscrow(c.Request.Context(), tradeID, userID)
	case "cancel":
		updateErr = h.tradeService.CancelTrade(c.Request.Context(), tradeID, userID, req.Reason)
	}

	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Trade %sed successfully", req.Action)})
}

func (h *TradeHandler) SendMessageHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	msg, err := h.tradeService.SendMessage(c.Request.Context(), tradeID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *TradeHandler) GetChatHistoryHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	messages, err := h.tradeService.GetChatHistory(c.Request.Context(), tradeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *TradeHandler) DisputeTradeHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req dto.RaiseDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	trade, err := h.tradeService.DisputeTrade(c.Request.Context(), tradeID, userID, req.ReasonCode, req.ReasonText)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Dispute raised successfully",
		"trade_id":   trade.TradeID.String(),
		"dispute_id": trade.DisputeID,
		"status":     trade.Status,
	})
}
