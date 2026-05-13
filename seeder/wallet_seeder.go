package seeder

import (
	"context"
	"fmt"
	"math/rand"

	domainidentity "cryplio/internal/domain/identity"

	"github.com/google/uuid"
)

func (s *Seeder) SeedWallets(ctx context.Context, users []*domainidentity.User, cryptoMap map[string]int) error {
	for _, user := range users {
		for symbol, cryptoID := range cryptoMap {
			address := fmt.Sprintf("0x%s%s%d", user.Username, symbol, rand.Intn(10000))
			if symbol == "BTC" {
				address = fmt.Sprintf("bc1q%s%d", user.Username, rand.Intn(10000))
			}

			balance := 500.0 + rand.Float64()*1000.0

			_, err := s.db.ExecContext(ctx, `
				INSERT INTO wallets (wallet_id, user_id, crypto_id, address, address_label, balance, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
				ON CONFLICT (user_id, crypto_id) DO UPDATE SET balance = EXCLUDED.balance`,
				uuid.New(), user.UserID, cryptoID, address, symbol+" Main Wallet", balance,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
