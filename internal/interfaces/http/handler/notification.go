package handler

import (
	"net/http"
	"strconv"

	"cryplio/internal/domain/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService notification.Service
}

func NewNotificationHandler(service notification.Service) *NotificationHandler {
	return &NotificationHandler{notificationService: service}
}

func (h *NotificationHandler) GetNotificationsHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.notificationService.GetNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkReadHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// GetPreferencesHandler gets user's notification preferences
func (h *NotificationHandler) GetPreferencesHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	prefs, err := h.notificationService.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// SavePreferencesHandler saves user's notification preferences
type SavePreferencesRequest struct {
	Email map[notification.NotificationType]bool `json:"email"`
	Push  map[notification.NotificationType]bool `json:"push"`
	SMS   map[notification.NotificationType]bool `json:"sms"`
}

func (h *NotificationHandler) SavePreferencesHandler(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req SavePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefs := &notification.NotificationPreference{
		Email: req.Email,
		Push:  req.Push,
		SMS:   req.SMS,
	}

	if err := h.notificationService.SavePreferences(c.Request.Context(), userID, prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences saved successfully"})
}
