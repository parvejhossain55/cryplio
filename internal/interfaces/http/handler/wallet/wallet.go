package wallet

import (
	"net/http"
	"strconv"
	"strings"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/domain/identity"
	walletdomain "cryplio/internal/domain/wallet"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService walletdomain.Service
	authService   identity.AuthService
}

func NewWalletHandler(walletService walletdomain.Service, authService identity.AuthService) *WalletHandler {
	return &WalletHandler{walletService: walletService, authService: authService}
}

func (h *WalletHandler) GetBalancesHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	balances, err := h.walletService.GetBalances(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": balances})
}

func (h *WalletHandler) CreateWalletHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req struct {
		CryptoSymbol string `json:"crypto_symbol" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	w, err := h.walletService.CreateWallet(c.Request.Context(), userID, req.CryptoSymbol)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"wallet":  w,
		"message": "wallet created successfully",
	})
}

func (h *WalletHandler) GetDepositAddressHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	crypto := strings.TrimSpace(c.Param("crypto"))
	w, err := h.walletService.GetDepositAddress(c.Request.Context(), userID, crypto)
	if err != nil {
		basehandler.HandleError(c, err)
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
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	// Check if 2FA is enabled for withdrawals
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}
	if !h.authService.Is2FAEnabled(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "2FA is required for withdrawals"})
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
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction": tx,
		"message":     "withdrawal requested",
	})
}

func (h *WalletHandler) GetTransactionsHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
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
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// ListPendingWithdrawalsHandler returns pending withdrawals requiring admin approval
func (h *WalletHandler) ListPendingWithdrawalsHandler(c *gin.Context) {
	limit := 50
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

	transactions, total, err := h.walletService.ListPendingWithdrawals(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"withdrawals": transactions,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
	})
}

// ApproveWithdrawalHandler approves a pending withdrawal
func (h *WalletHandler) ApproveWithdrawalHandler(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req struct {
		TxHash string `json:"tx_hash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletService.ApproveWithdrawal(c.Request.Context(), txID, adminID, req.TxHash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "withdrawal approved successfully"})
}

// RejectWithdrawalHandler rejects a pending withdrawal
func (h *WalletHandler) RejectWithdrawalHandler(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletService.RejectWithdrawal(c.Request.Context(), txID, adminID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "withdrawal rejected successfully"})
}
