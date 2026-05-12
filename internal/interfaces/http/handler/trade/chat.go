package trade

// chat.go contains TradeHandler methods for trade chat messages and feedback,
// plus the package-level isAllowedFileType helper.

import (
	"fmt"
	"net/http"
	"time"

	"cryplio/internal/domain/trading"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/http/dto"
	basehandler "cryplio/internal/interfaces/http/handler"
	"cryplio/internal/interfaces/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendMessageHandler sends a text or file message in a trade's chat.
func (h *TradeHandler) SendMessageHandler(c *gin.Context) {
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

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
		msg, err := h.communication.SendFileMessage(c.Request.Context(), tradeID, userID, uploadResult.URL, contentType, int(header.Size))
		if err != nil {
			basehandler.HandleError(c, err)
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

	msg, err := h.communication.SendMessage(c.Request.Context(), tradeID, userID, req.Content)
	if err != nil {
		basehandler.HandleError(c, err)
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

// GetChatHistoryHandler returns all chat messages for a given trade.
func (h *TradeHandler) GetChatHistoryHandler(c *gin.Context) {
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	messages, err := h.communication.GetChatHistory(c.Request.Context(), tradeID, userID)
	if err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// LeaveFeedbackHandler submits a rating and optional comment for a completed trade.
func (h *TradeHandler) LeaveFeedbackHandler(c *gin.Context) {
	tradeID, ok := basehandler.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.LeaveFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	rating := trading.FeedbackRating(req.Rating)
	if err := h.communication.LeaveFeedback(c.Request.Context(), tradeID, userID, rating, req.Comment); err != nil {
		basehandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback submitted successfully"})
}

// isAllowedFileType reports whether the given MIME type is permitted for trade chat uploads.
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
