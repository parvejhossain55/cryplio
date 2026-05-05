package handler

import (
	"net/http"
	"strconv"
	"strings"

	"cryplio/internal/domain/wallet"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService wallet.Service
}

func NewWalletHandler(walletService wallet.Service) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) GetBalancesHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	balances, err := h.walletService.GetBalances(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": balances})
}

func (h *WalletHandler) GetDepositAddressHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	crypto := strings.TrimSpace(c.Param("crypto"))
	w, err := h.walletService.GetDepositAddress(c.Request.Context(), userID, crypto)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet_id": w.WalletID,
		"crypto_id": w.CryptoID,
		"address":   w.Address,
	})
}

type withdrawRequest struct {
	CryptoSymbol string  `json:"crypto_symbol" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Fee          float64 `json:"fee"`
	Memo         *string `json:"memo"`
}

func (h *WalletHandler) WithdrawHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req withdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	tx, err := h.walletService.Withdraw(
		c.Request.Context(),
		userID,
		req.CryptoSymbol,
		req.Address,
		req.Amount,
		req.Fee,
		req.Memo,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction": tx,
		"message":     "withdrawal requested",
	})
}

func (h *WalletHandler) GetTransactionsHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	limit := 20
	offset := 0
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}
	if v := c.Query("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			offset = parsed
		}
	}

	transactions, total, err := h.walletService.GetTransactions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	userID, err := uuid.Parse(userIDRaw.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return uuid.Nil, false
	}
	return userID, true
}
