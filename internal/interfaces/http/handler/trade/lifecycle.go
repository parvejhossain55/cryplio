package trade

// lifecycle.go contains TradeHandler methods for trade status and lifecycle management.

import (
	"fmt"
	"net/http"

	"cryplio/internal/interfaces/http/dto"
	basehandler "cryplio/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
)

// InitiateTradeHandler starts a new trade against a given advertisement.
func (h *TradeHandler) InitiateTradeHandler(c *gin.Context) {
	adID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
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

	trade, err := h.lifecycle.InitiateTrade(c.Request.Context(), adID, buyerID, req.Amount)
	if err != nil {
		basehandler.HandleError(c, err)
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
	trades, err := h.lifecycle.ListTrades(c.Request.Context(), userID, role)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ListTradesResponse{
		Trades: func() []any {
			res := make([]any, len(trades))
			for i, t := range trades {
				res[i] = t
			}
			return res
		}(),
	})
}

// GetTradeHandler retrieves a single trade by ID.
func (h *TradeHandler) GetTradeHandler(c *gin.Context) {
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	trade, err := h.lifecycle.GetTrade(c.Request.Context(), tradeID)
	if err != nil {
		basehandler.HandleError(c, err)
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
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
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

	var actionErr error
	switch req.Action {
	case "pay":
		actionErr = h.lifecycle.MarkAsPaid(c.Request.Context(), tradeID, userID)
	case "release":
		actionErr = h.lifecycle.ReleaseEscrow(c.Request.Context(), tradeID, userID)
	case "cancel":
		actionErr = h.lifecycle.CancelTrade(c.Request.Context(), tradeID, userID, req.Reason)
	}

	if actionErr != nil {
		basehandler.HandleError(c, actionErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Trade action '%s' completed", req.Action)})
}

// DisputeTradeHandler initiates a dispute for a trade.
func (h *TradeHandler) DisputeTradeHandler(c *gin.Context) {
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req struct {
		ReasonCode string `json:"reason_code" binding:"required"`
		ReasonText string `json:"reason_text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	trade, err := h.lifecycle.DisputeTrade(c.Request.Context(), tradeID, userID, req.ReasonCode, req.ReasonText)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Dispute initiated",
		"trade_id": trade.TradeID.String(),
		"status":   trade.Status,
	})
}

// ListAllTradesHandler returns all trades in the system (admin only).
func (h *TradeHandler) ListAllTradesHandler(c *gin.Context) {
	status := c.Query("status")
	trades, err := h.lifecycle.ListAllTrades(c.Request.Context(), status)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ListTradesResponse{
		Trades: func() []any {
			res := make([]any, len(trades))
			for i, t := range trades {
				res[i] = t
			}
			return res
		}(),
	})
}
