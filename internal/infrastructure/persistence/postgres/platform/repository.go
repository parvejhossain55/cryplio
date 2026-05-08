package platform

// repository.go is the entry point for the platform Postgres repository.
// Method implementations are split across focused files:
//
//   crypto_asset.go   — CryptoAsset CRUD
//   fiat_currency.go  — FiatCurrency CRUD
//   payment_method.go — PaymentMethod CRUD

import (
	"database/sql"

	"cryplio/internal/domain/platform"
)

// platformRepository implements platform.PlatformRepository on top of PostgreSQL.
type platformRepository struct {
	db *sql.DB
}

// NewPlatformRepository constructs a platformRepository backed by the given *sql.DB.
func NewPlatformRepository(db *sql.DB) platform.PlatformRepository {
	return &platformRepository{db: db}
}
