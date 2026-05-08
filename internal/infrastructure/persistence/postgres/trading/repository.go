package trading

// repository.go is the entry point for the trading Postgres repository.
// Method implementations are split across focused files:
//
//   ad.go        — TradeAd CRUD + listing
//   trade.go     — Trade lifecycle CRUD + expiry queries
//   message.go   — Trade chat messages
//   feedback.go  — Post-trade feedback

import (
	"database/sql"

	"cryplio/internal/domain/trading"
)

// tradeRepository implements trading.TradeRepository on top of PostgreSQL.
type tradeRepository struct {
	db *sql.DB
}

// NewTradeRepository constructs a tradeRepository backed by the given *sql.DB.
func NewTradeRepository(db *sql.DB) trading.TradeRepository {
	return &tradeRepository{db: db}
}
