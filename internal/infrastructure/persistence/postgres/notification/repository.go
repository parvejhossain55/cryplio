package notification

import (
	"context"
	"database/sql"
	"fmt"

	"cryplio/internal/domain/notification"

	"github.com/google/uuid"
)

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) notification.Repository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	query := `
		INSERT INTO notifications (
			notification_id, user_id, type, title, message, data, is_read, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW()
		) RETURNING created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		n.NotificationID, n.UserID, n.Type, n.Title, n.Message, n.Data, n.IsRead,
	).Scan(&n.CreatedAt)

	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *notificationRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]notification.Notification, error) {
	query := `
		SELECT notification_id, user_id, type, title, message, data, is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		err := rows.Scan(
			&n.NotificationID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data,
			&n.IsRead, &n.ReadAt, &n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true, read_at = NOW() WHERE notification_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*notification.NotificationPreference, error) {
	query := `SELECT user_id, email_prefs, push_prefs, sms_prefs, created_at, updated_at FROM notification_preferences WHERE user_id = $1`
	var p notification.NotificationPreference
	// Note: We'll need to handle the JSONB fields properly in a real app,
	// for now we'll just mock the preferences or return a default.
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&p.UserID, nil, nil, nil, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default preferences
			return &notification.NotificationPreference{UserID: userID}, nil
		}
		return nil, err
	}
	return &p, nil
}
