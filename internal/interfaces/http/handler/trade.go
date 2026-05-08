package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cryplio/internal/domain/trading"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	"cryplio/internal/interfaces/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TradeHandler struct {
	tradeService trading.TradeService
	storage      storage.ObjectStorage
	wsService    websocket.Service
}

func NewTradeHandler(service trading.TradeService, storage storage.ObjectStorage, wsService websocket.Service) *TradeHandler {
	return &TradeHandler{tradeService: service, storage: storage, wsService: wsService}
}

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

func (h *TradeHandler) UpdateAdHandler(c *gin.Context) {
	adID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}

	var req dto.UpdateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

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

	if err := h.tradeService.UpdateAd(c.Request.Context(), adID, userID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad updated successfully"})
}

func (h *TradeHandler) DeleteAdHandler(c *gin.Context) {
	adID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.tradeService.DeleteAd(c.Request.Context(), adID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad deleted successfully"})
}

func (h *TradeHandler) CreateAdHandler(c *gin.Context) {
	var req dto.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

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

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	// Check if this is a file upload request
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		// Handle file upload
		defer file.Close()

		// Validate file size (max 5MB)
		if header.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file size must be less than 5MB"})
			return
		}

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		if !isAllowedFileType(contentType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			return
		}

		// Read file content
		fileContent := make([]byte, header.Size)
		_, err := file.Read(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		// Upload to storage
		fileName := header.Filename
		uploadInput := storage.UploadInput{
			Key:         fmt.Sprintf("trade-files/%s/%s/%s", tradeID.String(), userID.String(), fileName),
			ContentType: contentType,
			Body:        fileContent,
		}

		uploadResult, err := h.storage.Upload(c.Request.Context(), uploadInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		// Send file message
		msg, err := h.tradeService.SendFileMessage(c.Request.Context(), tradeID, userID, uploadResult.URL, contentType, int(header.Size))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Broadcast message to trade participants via WebSocket
		if h.wsService != nil && msg.MessageID != uuid.Nil {
			var fileURL, mimeType string
			var fileSize int

			if msg.FileURL != nil {
				fileURL = *msg.FileURL
			}
			if msg.FileMimeType != nil {
				mimeType = *msg.FileMimeType
			}
			if msg.FileSize != nil {
				fileSize = *msg.FileSize
			}

			chatMessage := websocket.ChatMessage{
				ID:        msg.MessageID.String(),
				TradeID:   msg.TradeID.String(),
				SenderID:  msg.SenderID.String(),
				FileURL:   fileURL,
				MimeType:  mimeType,
				FileSize:  fileSize,
				CreatedAt: msg.CreatedAt.Format(time.RFC3339),
			}

			h.wsService.BroadcastMessage("chat_message", chatMessage, tradeID.String())
		}

		c.JSON(http.StatusCreated, msg)
		return
	}

	// Handle text message
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required for text messages"})
		return
	}

	msg, err := h.tradeService.SendMessage(c.Request.Context(), tradeID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast message to trade participants via WebSocket
	if h.wsService != nil && msg.MessageID != uuid.Nil {
		var content string
		if msg.Content != nil {
			content = *msg.Content
		}

		chatMessage := websocket.ChatMessage{
			ID:        msg.MessageID.String(),
			TradeID:   msg.TradeID.String(),
			SenderID:  msg.SenderID.String(),
			Content:   content,
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		}

		h.wsService.BroadcastMessage("chat_message", chatMessage, tradeID.String())
	}

	c.JSON(http.StatusCreated, msg)
}

func isAllowedFileType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg", "image/png", "image/gif",
		"application/pdf",
	}
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
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

// ListAllTradesHandler returns all trades for admin monitoring
func (h *TradeHandler) ListAllTradesHandler(c *gin.Context) {
	status := c.Query("status")
	trades, err := h.tradeService.ListAllTrades(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trades": trades})
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

func (h *TradeHandler) LeaveFeedbackHandler(c *gin.Context) {
	tradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req dto.LeaveFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	rating := trading.FeedbackRating(req.Rating)
	err = h.tradeService.LeaveFeedback(c.Request.Context(), tradeID, userID, rating, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback submitted successfully"})
}
