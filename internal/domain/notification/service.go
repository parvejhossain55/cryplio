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

	// For MVP, send email for critical trade-related notifications
	criticalTypes := map[NotificationType]bool{
		NotificationTypeTradeStarted:        true,
		NotificationTypeTradePaid:           true,
		NotificationTypeTradeReleased:       true,
		NotificationTypeTradeCompleted:      true,
		NotificationTypeTradeCancelled:      true,
		NotificationTypeTradeDisputed:       true,
		NotificationTypeDisputeResolved:     true,
		NotificationTypeDepositReceived:     true,
		NotificationTypeWithdrawalCompleted: true,
	}
	if s.emailClient != nil && criticalTypes[nType] {
		_ = s.emailClient.SendEmail(ctx, userID.String()+"@cryplio.local", title, message)
	}

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
