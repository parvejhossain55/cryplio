package seeder

import (
	"context"
	"encoding/json"

	domainidentity "cryplio/internal/domain/identity"

	"github.com/google/uuid"
)

func (s *Seeder) SeedNotifications(ctx context.Context, users []*domainidentity.User) error {
	for _, u := range users {
		dataJson, _ := json.Marshal(map[string]string{"info": "welcome"})
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO notifications (notification_id, user_id, type, title, message, data, is_read, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, false, NOW())`,
			uuid.New(), u.UserID, "system_announcement", "Welcome to Cryplio", "Thanks for joining the premier P2P exchange.", dataJson,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
