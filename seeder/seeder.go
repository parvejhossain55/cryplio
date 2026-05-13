package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainidentity "cryplio/internal/domain/identity"
	persistenceidentity "cryplio/internal/infrastructure/persistence/postgres/identity"
	persistencetrading "cryplio/internal/infrastructure/persistence/postgres/trading"
	"cryplio/pkg/config"

	domaintrading "cryplio/internal/domain/trading"
)

// Seeder manages the database seeding process
type Seeder struct {
	db        *sql.DB
	cfg       *config.Config
	userRepo  domainidentity.UserRepository
	tradeRepo domaintrading.TradeRepository
}

// NewSeeder creates a new database seeder
func NewSeeder(db *sql.DB, cfg *config.Config) *Seeder {
	return &Seeder{
		db:        db,
		cfg:       cfg,
		userRepo:  persistenceidentity.NewUserRepository(db),
		tradeRepo: persistencetrading.NewTradeRepository(db),
	}
}

// SeedAll runs all registered seeders
func (s *Seeder) SeedAll(ctx context.Context) error {
	fmt.Println("🌱 Starting comprehensive database seeding...")
	start := time.Now()

	// 0. Environment Setup
	cryptoMap, fiatMap, pmMap, err := s.getLookupData()
	if err != nil {
		return fmt.Errorf("lookup data: %w", err)
	}

	// 1. Seed Users
	users, err := s.SeedUsers(ctx)
	if err != nil {
		return fmt.Errorf("users: %w", err)
	}
	fmt.Printf("✅ Seeded %d users\n", len(users))

	// 2. Wallets
	if err := s.SeedWallets(ctx, users, cryptoMap); err != nil {
		return fmt.Errorf("wallets: %w", err)
	}
	fmt.Println("✅ Seeded wallets and balances")

	// 3. Ads
	ads, err := s.SeedTradeAds(ctx, users, cryptoMap, fiatMap, pmMap)
	if err != nil {
		return fmt.Errorf("ads: %w", err)
	}
	fmt.Printf("✅ Seeded %d trade ads\n", len(ads))

	// 4. Trades & Feedback
	trades, err := s.SeedTrades(ctx, users, ads)
	if err != nil {
		return fmt.Errorf("trades: %w", err)
	}
	fmt.Printf("✅ Seeded %d trades and feedback\n", len(trades))

	// 5. Disputes
	if err := s.SeedDisputes(ctx, trades); err != nil {
		return fmt.Errorf("disputes: %w", err)
	}
	fmt.Println("✅ Seeded disputes")

	// 6. Notifications
	if err := s.SeedNotifications(ctx, users); err != nil {
		return fmt.Errorf("notifications: %w", err)
	}
	fmt.Println("✅ Seeded notifications")

	fmt.Printf("✨ Seeding completed successfully in %v!\n", time.Since(start))
	return nil
}

func (s *Seeder) getLookupData() (map[string]int, map[string]int, map[string]int, error) {
	cryptoMap := make(map[string]int)
	rows, err := s.db.Query("SELECT id, symbol FROM crypto_assets")
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var symbol string
		if err := rows.Scan(&id, &symbol); err != nil {
			return nil, nil, nil, err
		}
		cryptoMap[symbol] = id
	}

	fiatMap := make(map[string]int)
	rows, err = s.db.Query("SELECT id, code FROM fiat_currencies")
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			return nil, nil, nil, err
		}
		fiatMap[code] = id
	}

	pmMap := make(map[string]int)
	rows, err = s.db.Query("SELECT id, code FROM payment_methods")
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			return nil, nil, nil, err
		}
		pmMap[code] = id
	}

	return cryptoMap, fiatMap, pmMap, nil
}
