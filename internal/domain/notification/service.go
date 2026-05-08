package notification

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type EmailClient interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type WebSocketNotifier interface {
	NotifyUser(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data map[string]interface{})
}

type Service interface {
	Notify(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data *string) error
	NotifyWithWebSocket(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data map[string]interface{}) error
	GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	SetWebSocketNotifier(notifier WebSocketNotifier)
	GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreference, error)
	SavePreferences(ctx context.Context, userID uuid.UUID, prefs *NotificationPreference) error
}

type notificationService struct {
	repo        Repository
	emailClient EmailClient
	wsNotifier  WebSocketNotifier
}

func NewService(repo Repository, emailClient EmailClient) Service {
	return &notificationService{
		repo:        repo,
		emailClient: emailClient,
	}
}

func (s *notificationService) SetWebSocketNotifier(notifier WebSocketNotifier) {
	s.wsNotifier = notifier
}

func (s *notificationService) Notify(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data *string) error {
	n := &Notification{
		NotificationID: uuid.New(),
		UserID:         userID,
		Type:           nType,
		Title:          title,
		Message:        message,
		Data:           data,
		IsRead:         false,
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}

	return nil
}

func (s *notificationService) NotifyWithWebSocket(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data map[string]interface{}) error {
	// Save to database first
	var dataStr *string
	if data != nil {
		d, _ := json.Marshal(data)
		ds := string(d)
		dataStr = &ds
	}

	if err := s.Notify(ctx, userID, nType, title, message, dataStr); err != nil {
		return err
	}

	// Send real-time WebSocket notification
	if s.wsNotifier != nil {
		s.wsNotifier.NotifyUser(ctx, userID, nType, title, message, data)
	}

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error) {
	return s.repo.ListByUserID(ctx, userID, limit, offset)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *notificationService) GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreference, error) {
	return s.repo.GetPreferences(ctx, userID)
}

func (s *notificationService) SavePreferences(ctx context.Context, userID uuid.UUID, prefs *NotificationPreference) error {
	prefs.UserID = userID
	return s.repo.SavePreferences(ctx, prefs)
}
