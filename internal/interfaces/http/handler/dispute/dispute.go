package dispute

import (
	"fmt"
	"net/http"
	"strconv"

	basehandler "cryplio/internal/interfaces/http/handler"

	disputedomain "cryplio/internal/domain/dispute"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DisputeHandler struct {
	disputeService disputedomain.Service
	storage        storage.ObjectStorage
}

func NewDisputeHandler(service disputedomain.Service, s storage.ObjectStorage) *DisputeHandler {
	return &DisputeHandler{disputeService: service, storage: s}
}

func (h *DisputeHandler) GetDisputeHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	d, err := h.disputeService.GetDispute(c.Request.Context(), id)
	if err != nil {
		basehandler.HandleError(c, err)
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

	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

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

	adminID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	resolution := disputedomain.DisputeResolution(req.Resolution)
	if err := h.disputeService.ResolveDispute(c.Request.Context(), id, adminID, resolution, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Trigger trade completion/cancellation based on resolution
	// This might be better handled in the ResolveDispute method itself or via events

	c.JSON(http.StatusOK, gin.H{"message": "Dispute resolved successfully"})
}

func (h *DisputeHandler) UploadEvidenceHandler(c *gin.Context) {
	disputeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	file, err := c.FormFile("evidence")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "evidence file is required"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer openedFile.Close()

	data := make([]byte, file.Size)
	_, err = openedFile.Read(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	key := fmt.Sprintf("disputes/%s/%s_%s", disputeID.String(), uuid.New().String(), file.Filename)
	result, err := h.storage.Upload(c.Request.Context(), storage.UploadInput{
		Key:         key,
		ContentType: file.Header.Get("Content-Type"),
		Body:        data,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("upload failed: %v", err)})
		return
	}

	if err := h.disputeService.UploadEvidence(c.Request.Context(), disputeID, userID, result.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "evidence uploaded successfully", "url": result.URL})
}

func (h *DisputeHandler) ListDisputesHandler(c *gin.Context) {
	limit := 50
	offset := 0
	status := c.Query("status")
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	disputes, total, err := h.disputeService.ListDisputes(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"disputes": disputes, "total": total, "limit": limit, "offset": offset})
}
