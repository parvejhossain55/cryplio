package seeder

import (
	"context"

	domainidentity "cryplio/internal/domain/identity"
)

func (s *Seeder) SeedKYC(ctx context.Context, users []*domainidentity.User) error {
	for _, user := range users {
		if user.Username == "admin" {
			continue
		}

		status := user.KYCLevel

		_, err := s.db.ExecContext(ctx, `
			INSERT INTO kyc_records (user_id, level, document_type, document_front_url, document_back_url, selfie_url, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			ON CONFLICT DO NOTHING`,
			user.UserID, user.KYCLevel, "passport", "https://example.com/front.jpg", "https://example.com/back.jpg", "https://example.com/selfie.jpg", status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
