package notification

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type EmailClient interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type Service interface {
	Notify(ctx context.Context, userID uuid.UUID, nType NotificationType, title, message string, data *string) error
	GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}

type notificationService struct {
	repo        Repository
	emailClient EmailClient
}

func NewService(repo Repository, emailClient EmailClient) Service {
	return &notificationService{
		repo:        repo,
		emailClient: emailClient,
	}
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

	// In a real app, we would check user preferences here
	// and send email/push/sms if enabled.
	// For MVP, we'll try to send an email for critical notifications if an email client is provided.
	_ = s.emailClient // Placeholder for email logic

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error) {
	return s.repo.ListByUserID(ctx, userID, limit, offset)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkRead(ctx, id)
}

func ValidateNotification(notification *Notification) error {
	if notification == nil {
		return errors.New("notification is required")
	}
	return nil
}
