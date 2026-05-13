package trading

import (
	"context"
	"fmt"
	"time"

	"cryplio/internal/domain/trading"

	"github.com/google/uuid"
)

// ─── Trade Lifecycle CRUD ─────────────────────────────────────────────────────

func (r *tradeRepository) CreateTrade(ctx context.Context, t *trading.Trade) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO trades (
			trade_id, ad_id, buyer_id, seller_id, crypto_id, fiat_id,
			crypto_amount, fiat_amount, rate, status, payment_method_code,
			escrow_wallet_id, tx_hash, payment_window_minutes,
			expires_at, paid_at, released_at, completed_at, cancelled_at,
			disputed_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20, NOW(), NOW())
		RETURNING created_at, updated_at`,
		t.TradeID, t.AdID, t.BuyerID, t.SellerID,
		t.CryptoID, t.FiatID, t.CryptoAmount, t.FiatAmount,
		t.Rate, t.Status, t.PaymentMethodCode, t.EscrowWalletID,
		t.TxHash, t.PaymentWindowMinutes, t.ExpiresAt, t.PaidAt,
		t.ReleasedAt, t.CompletedAt, t.CancelledAt, t.DisputedAt,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *tradeRepository) GetTradeByID(ctx context.Context, id uuid.UUID) (*trading.Trade, error) {
	var t trading.Trade
	err := scanTrade(
		r.db.QueryRowContext(ctx,
			`SELECT `+tradeColumns+` FROM trades WHERE trade_id = $1 AND deleted_at IS NULL`, id),
		&t,
	)
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get trade: %w", err)
	}
	return &t, nil
}

func (r *tradeRepository) ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]trading.Trade, error) {
	var where string
	switch role {
	case "buyer":
		where = "buyer_id = $1"
	case "seller":
		where = "seller_id = $1"
	default:
		where = "(buyer_id = $1 OR seller_id = $1)"
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+tradeColumns+` FROM trades WHERE `+where+` ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("list trades: %w", err)
	}
	defer rows.Close()
	return scanTrades(rows)
}

func (r *tradeRepository) ListAllTrades(ctx context.Context, status string) ([]trading.Trade, error) {
	var query string
	var args []interface{}
	if status != "" && status != "all" {
		query = `SELECT ` + tradeColumns + ` FROM trades WHERE status = $1 ORDER BY created_at DESC`
		args = append(args, status)
	} else {
		query = `SELECT ` + tradeColumns + ` FROM trades ORDER BY created_at DESC`
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list all trades: %w", err)
	}
	defer rows.Close()
	return scanTrades(rows)
}

func (r *tradeRepository) CountTrades(ctx context.Context, status string) (int, error) {
	var count int
	var err error
	if status != "" && status != "all" {
		err = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM trades WHERE status = $1`, status).Scan(&count)
	} else {
		err = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM trades`).Scan(&count)
	}
	if err != nil {
		return 0, fmt.Errorf("count trades: %w", err)
	}
	return count, nil
}

func (r *tradeRepository) ListExpiredPendingTrades(ctx context.Context, now time.Time) ([]trading.Trade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+tradeColumns+`
		FROM trades
		WHERE status IN ('pending', 'active')
		  AND paid_at IS NULL
		  AND (created_at + make_interval(mins => payment_window_minutes)) <= $1`, now)
	if err != nil {
		return nil, fmt.Errorf("list expired pending trades: %w", err)
	}
	defer rows.Close()
	return scanTrades(rows)
}

func (r *tradeRepository) UpdateTrade(ctx context.Context, t *trading.Trade) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE trades
		SET status = $1, paid_at = $2, released_at = $3,
		    cancelled_at = $4, completed_at = $5, expires_at = $6,
		    disputed_at = $7, updated_at = NOW()
		WHERE trade_id = $8`,
		t.Status, t.PaidAt, t.ReleasedAt, t.CancelledAt,
		t.CompletedAt, t.ExpiresAt, t.DisputedAt, t.TradeID,
	)
	if err != nil {
		return fmt.Errorf("update trade: %w", err)
	}
	return nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// scanTrades iterates a rows cursor and scans each row into a Trade.
func scanTrades(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]trading.Trade, error) {
	var trades []trading.Trade
	for rows.Next() {
		var t trading.Trade
		if err := scanTrade(rows, &t); err != nil {
			return nil, fmt.Errorf("scan trade: %w", err)
		}
		trades = append(trades, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate trades: %w", err)
	}
	return trades, nil
}

// isNoRows reports whether err is sql.ErrNoRows.
func isNoRows(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "sql: no rows in result set"
}
