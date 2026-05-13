package notification

import (
	"net/http"
	"strconv"

	basehandler "cryplio/internal/interfaces/http/handler"

	notificationdomain "cryplio/internal/domain/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService notificationdomain.Service
}

func NewNotificationHandler(service notificationdomain.Service) *NotificationHandler {
	return &NotificationHandler{notificationService: service}
}

func (h *NotificationHandler) GetNotificationsHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.notificationService.GetNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if notifications == nil {
		notifications = []notificationdomain.Notification{}
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
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
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	prefs, err := h.notificationService.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// SavePreferencesHandler saves user's notification preferences
type SavePreferencesRequest struct {
	Email map[notificationdomain.NotificationType]bool `json:"email"`
	Push  map[notificationdomain.NotificationType]bool `json:"push"`
	SMS   map[notificationdomain.NotificationType]bool `json:"sms"`
}

func (h *NotificationHandler) SavePreferencesHandler(c *gin.Context) {
	userID, ok := basehandler.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var req SavePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefs := &notificationdomain.NotificationPreference{
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
