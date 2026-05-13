package trade

// lifecycle.go contains TradeHandler methods for trade status and lifecycle management.

import (
	"fmt"
	"net/http"

	"cryplio/internal/interfaces/http/dto"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InitiateTradeHandler starts a new trade against a given advertisement.
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

	buyerID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

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

// ListTradesHandler returns all trades for the authenticated user, optionally filtered by role.
func (h *TradeHandler) ListTradesHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}
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

// GetTradeHandler retrieves a single trade by ID.
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

// UpdateTradeStatusHandler advances a trade through its lifecycle (pay / release / cancel).
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

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

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

// ListAllTradesHandler returns all trades for admin monitoring.
func (h *TradeHandler) ListAllTradesHandler(c *gin.Context) {
	status := c.Query("status")
	trades, err := h.tradeService.ListAllTrades(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trades": trades})
}

// DisputeTradeHandler raises a dispute for an in-progress trade.
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

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	trade, err := h.tradeService.DisputeTrade(c.Request.Context(), tradeID, userID, req.ReasonCode, req.ReasonText)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Dispute raised successfully",
		"trade_id": trade.TradeID.String(),
		"status":   trade.Status,
	})
}
