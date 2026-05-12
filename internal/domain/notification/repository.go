package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
	MarkRead(ctx context.Context, notificationID uuid.UUID) error
	GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreference, error)
	SavePreferences(ctx context.Context, prefs *NotificationPreference) error
}
