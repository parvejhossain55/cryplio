package handler

import (
	"net/http"

	"cryplio/internal/domain/dispute"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DisputeHandler struct {
	disputeService dispute.Service
}

func NewDisputeHandler(service dispute.Service) *DisputeHandler {
	return &DisputeHandler{disputeService: service}
}

func (h *DisputeHandler) GetDisputeHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	d, err := h.disputeService.GetDispute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if d == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dispute not found"})
		return
	}

	c.JSON(http.StatusOK, d)
}

func (h *DisputeHandler) AssignDisputeHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	if err := h.disputeService.AssignDispute(c.Request.Context(), id, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dispute assigned successfully"})
}

func (h *DisputeHandler) ResolveDisputeHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	var req dto.ResolveDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminIDStr, _ := c.Get("user_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	resolution := dispute.DisputeResolution(req.Resolution)
	if err := h.disputeService.ResolveDispute(c.Request.Context(), id, adminID, resolution, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Trigger trade completion/cancellation based on resolution
	// This might be better handled in the ResolveDispute method itself or via events

	c.JSON(http.StatusOK, gin.H{"message": "Dispute resolved successfully"})
}
func (h *DisputeHandler) ListDisputesHandler(c *gin.Context) {
	disputes, err := h.disputeService.ListDisputes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"disputes": disputes})
}
