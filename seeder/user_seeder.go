package seeder

import (
	"context"

	domainidentity "cryplio/internal/domain/identity"
	"cryplio/pkg/crypto"
)

func (s *Seeder) SeedUsers(ctx context.Context) ([]*domainidentity.User, error) {
	passwordHash, _ := crypto.HashPassword("Password123!")

	userData := []struct {
		email    string
		username string
	}{
		{"admin@cryplio.com", "admin"},
		{"trader.one@example.com", "CryptoKing"},
		{"trader.two@example.com", "SwiftExchange"},
		{"trader.alice@example.com", "AliceTrader"},
		{"trader.bob@example.com", "BobCrypto"},
		{"trader.charlie@example.com", "CharlieP2P"},
		{"trader.diana@example.com", "DianaCoin"},
		{"trader.ethan@example.com", "EthanX"},
	}

	var users []*domainidentity.User
	for _, ud := range userData {
		existing, _ := s.userRepo.GetByEmail(ctx, ud.email)
		if existing != nil {
			users = append(users, existing)
			continue
		}

		user := domainidentity.NewUser(ud.email, ud.username, passwordHash)
		user.EmailVerified = true
		user.PhoneVerified = true
		user.Status = domainidentity.UserStatusActive

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
