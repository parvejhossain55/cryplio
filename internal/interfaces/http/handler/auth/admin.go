package auth

import (
	"fmt"
	"net/http"
	"time"

	basehandler "cryplio/internal/interfaces/http/handler"

	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/trading"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler handles all /admin/* routes.
// It is intentionally separate from AuthHandler to keep auth concerns
// and admin management concerns independent.
type AdminHandler struct {
	authService    identity.AuthService
	tradeService   trading.TradeService
	disputeService dispute.Service
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(
	authService identity.AuthService,
	tradeService trading.TradeService,
	disputeService dispute.Service,
) *AdminHandler {
	return &AdminHandler{
		authService:    authService,
		tradeService:   tradeService,
		disputeService: disputeService,
	}
}

// GetDashboardStatsHandler returns aggregated admin dashboard metrics.
func (h *AdminHandler) GetDashboardStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	stats := identity.DashboardStats{}

	if h.authService != nil {
		stats.TotalUsers, _ = h.authService.CountUsers(ctx)
	}
	if h.tradeService != nil {
		stats.TotalTrades, _ = h.tradeService.CountTrades(ctx, "")
		stats.PendingTrades, _ = h.tradeService.CountTrades(ctx, "pending")
		stats.ActiveTrades, _ = h.tradeService.CountTrades(ctx, "active")
		stats.PaidTrades, _ = h.tradeService.CountTrades(ctx, "paid")
		stats.CompletedTrades, _ = h.tradeService.CountTrades(ctx, "completed")
		stats.DisputedTrades, _ = h.tradeService.CountTrades(ctx, "disputed")
		stats.CancelledTrades, _ = h.tradeService.CountTrades(ctx, "cancelled")
	}
	if h.disputeService != nil {
		stats.TotalDisputes, _ = h.disputeService.CountDisputes(ctx, "")
		stats.PendingDisputes, _ = h.disputeService.CountDisputes(ctx, "pending")
	}

	c.JSON(http.StatusOK, stats)
}

// ListUsersHandler returns a paginated list of users (admin only).
func (h *AdminHandler) ListUsersHandler(c *gin.Context) {
	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := fmt.Sscanf(l, "%d", &limit); err == nil && parsed == 1 {
			if limit > 100 {
				limit = 100
			}
		}
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	users, err := h.authService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, mapUser(&u))
	}
	c.JSON(http.StatusOK, gin.H{"users": response, "limit": limit, "offset": offset})
}

// SuspendUserHandler suspends a user account (admin only).
func (h *AdminHandler) SuspendUserHandler(c *gin.Context) {
	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Reason   string `json:"reason" binding:"required"`
		Duration *int   `json:"duration_minutes,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var duration *time.Duration
	if req.Duration != nil && *req.Duration > 0 {
		d := time.Duration(*req.Duration) * time.Minute
		duration = &d
	}

	if err := h.authService.SuspendUser(c.Request.Context(), adminID, userID, req.Reason, duration); err != nil {
		basehandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user suspended successfully"})
}

// UnsuspendUserHandler lifts a user suspension (admin only).
func (h *AdminHandler) UnsuspendUserHandler(c *gin.Context) {
	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.authService.UnsuspendUser(c.Request.Context(), adminID, userID); err != nil {
		basehandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unsuspended successfully"})
}

// BanUserHandler permanently bans a user account (admin only).
func (h *AdminHandler) BanUserHandler(c *gin.Context) {
	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.BanUser(c.Request.Context(), adminID, userID, req.Reason); err != nil {
		basehandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user banned successfully"})
}

// UnbanUserHandler unbans a user account (admin only).
func (h *AdminHandler) UnbanUserHandler(c *gin.Context) {
	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.authService.UnbanUser(c.Request.Context(), adminID, userID); err != nil {
		basehandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unbanned successfully"})
}
