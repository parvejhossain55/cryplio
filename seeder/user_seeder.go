package seeder

import (
	"context"
	"fmt"

	domainidentity "cryplio/internal/domain/identity"
	"cryplio/pkg/crypto"
)

func (s *Seeder) SeedUsers(ctx context.Context) ([]*domainidentity.User, error) {
	passwordHash, _ := crypto.HashPassword("Password123!")

	userData := []struct {
		email      string
		username   string
		isMerchant bool
		kycLevel   domainidentity.KYCLevel
	}{
		{"admin@cryplio.com", "admin", false, domainidentity.KYCLevel3},
		{"merchant.one@example.com", "CryptoKing", true, domainidentity.KYCLevel3},
		{"merchant.two@example.com", "SwiftExchange", true, domainidentity.KYCLevel2},
		{"trader.alice@example.com", "AliceTrader", false, domainidentity.KYCLevel1},
		{"trader.bob@example.com", "BobCrypto", false, domainidentity.KYCLevel1},
		{"trader.charlie@example.com", "CharlieP2P", false, domainidentity.KYCLevel0},
		{"trader.diana@example.com", "DianaCoin", false, domainidentity.KYCLevel2},
		{"trader.ethan@example.com", "EthanX", false, domainidentity.KYCLevel1},
	}

	var users []*domainidentity.User
	for _, ud := range userData {
		existing, _ := s.userRepo.GetByEmail(ctx, ud.email)
		if existing != nil {
			users = append(users, existing)
			continue
		}

		user := domainidentity.NewUser(ud.email, ud.username, passwordHash)
		user.IsMerchant = ud.isMerchant
		user.KYCLevel = ud.kycLevel
		user.EmailVerified = true
		user.PhoneVerified = true
		user.Status = domainidentity.UserStatusActive

		if ud.isMerchant {
			bio := fmt.Sprintf("Professional merchant since 2024. Quick release, %s specialist.", ud.username)
			user.Bio = &bio
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
