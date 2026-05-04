package trading

import (
	"context"
	domainidentity "cryplio/internal/domain/identity"
	"cryplio/internal/domain/trading"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type tradeRepository struct {
	db *sql.DB
}

// NewTradeRepository creates a new postgres trade repository
func NewTradeRepository(db *sql.DB) trading.TradeRepository {
	return &tradeRepository{db: db}
}

func (r *tradeRepository) CreateAd(ctx context.Context, ad *trading.TradeAd) error {
	query := `
		INSERT INTO trade_ads (
			ad_id, user_id, type, crypto_id, fiat_id, price_type, price,
			floating_markup, min_amount, max_amount, payment_methods,
			trade_terms, payment_window_minutes, requires_kyc_level,
			is_public, is_paused, timezone, status, published_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, NOW(), NOW(), NOW()
		) RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		ad.AdID, ad.UserID, ad.Type, ad.CryptoID, ad.FiatID, ad.PriceType,
		ad.Price, ad.FloatingMarkup, ad.MinAmount, ad.MaxAmount,
		pq.Array(ad.PaymentMethods), ad.TradeTerms, ad.PaymentWindowMinutes,
		ad.RequiresKYCLevel, ad.IsPublic, ad.IsPaused, ad.Timezone, ad.Status,
	).Scan(&ad.CreatedAt, &ad.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create ad: %w", err)
	}
	return nil
}

func (r *tradeRepository) GetAdByID(ctx context.Context, id uuid.UUID) (*trading.TradeAd, error) {
	query := `
		SELECT ad_id, user_id, type, crypto_id, fiat_id, price_type, price,
		       floating_markup, min_amount, max_amount, payment_methods,
		       trade_terms, payment_window_minutes, requires_kyc_level,
		       is_public, is_paused, timezone, status, published_at,
		       created_at, updated_at
		FROM trade_ads
		WHERE ad_id = $1 AND deleted_at IS NULL
	`
	var ad trading.TradeAd
	var paymentMethods []int64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ad.AdID, &ad.UserID, &ad.Type, &ad.CryptoID, &ad.FiatID, &ad.PriceType,
		&ad.Price, &ad.FloatingMarkup, &ad.MinAmount, &ad.MaxAmount,
		pq.Array(&paymentMethods), &ad.TradeTerms, &ad.PaymentWindowMinutes,
		&ad.RequiresKYCLevel, &ad.IsPublic, &ad.IsPaused, &ad.Timezone,
		&ad.Status, &ad.PublishedAt, &ad.CreatedAt, &ad.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get ad: %w", err)
	}

	ad.PaymentMethods = make([]int, len(paymentMethods))
	for i, pm := range paymentMethods {
		ad.PaymentMethods[i] = int(pm)
	}

	return &ad, nil
}

func (r *tradeRepository) ListAds(ctx context.Context, filter trading.AdFilter) ([]trading.TradeAd, int, error) {
	var args []interface{}
	var conditions []string
	placeholder := 1

	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("a.type = $%d", placeholder))
		args = append(args, *filter.Type)
		placeholder++
	}
	if filter.CryptoID != nil {
		conditions = append(conditions, fmt.Sprintf("a.crypto_id = $%d", placeholder))
		args = append(args, *filter.CryptoID)
		placeholder++
	}
	if filter.FiatID != nil {
		conditions = append(conditions, fmt.Sprintf("a.fiat_id = $%d", placeholder))
		args = append(args, *filter.FiatID)
		placeholder++
	}
	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("a.user_id = $%d", placeholder))
		args = append(args, *filter.UserID)
		placeholder++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("a.status = $%d", placeholder))
		args = append(args, *filter.Status)
		placeholder++
	} else {
		conditions = append(conditions, "a.status = 'active'")
	}

	conditions = append(conditions, "a.deleted_at IS NULL")
	if filter.UserID == nil {
		conditions = append(conditions, "a.is_public = true")
		conditions = append(conditions, "a.is_paused = false")
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM trade_ads a " + where
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}

	// Fetch page
	query := `
		SELECT a.ad_id, a.user_id, a.type, a.crypto_id, a.fiat_id, a.price_type, a.price,
		       a.floating_markup, a.min_amount, a.max_amount, a.payment_methods,
		       a.trade_terms, a.payment_window_minutes, a.requires_kyc_level,
		       a.is_public, a.is_paused, a.timezone, a.status, a.published_at,
		       a.created_at, a.updated_at,
		       u.username, u.avatar_url, u.last_seen_at,
		       COALESCE(us.total_trades, 0), COALESCE(us.avg_rating, 0)
		FROM trade_ads a
		JOIN users u ON a.user_id = u.user_id
		LEFT JOIN user_stats us ON u.user_id = us.user_id
	` + where + fmt.Sprintf(" ORDER BY a.published_at DESC LIMIT $%d OFFSET $%d", placeholder, placeholder+1)

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	args = append(args, limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query ads: %w", err)
	}
	defer rows.Close()

	var ads []trading.TradeAd
	for rows.Next() {
		var ad trading.TradeAd
		ad.User = &domainidentity.User{}
		ad.Stats = &domainidentity.UserStats{}
		var paymentMethods []int64
		err := rows.Scan(
			&ad.AdID, &ad.UserID, &ad.Type, &ad.CryptoID, &ad.FiatID, &ad.PriceType,
			&ad.Price, &ad.FloatingMarkup, &ad.MinAmount, &ad.MaxAmount,
			pq.Array(&paymentMethods), &ad.TradeTerms, &ad.PaymentWindowMinutes,
			&ad.RequiresKYCLevel, &ad.IsPublic, &ad.IsPaused, &ad.Timezone,
			&ad.Status, &ad.PublishedAt, &ad.CreatedAt, &ad.UpdatedAt,
			&ad.User.Username, &ad.User.AvatarURL, &ad.User.LastSeenAt,
			&ad.Stats.TotalTrades, &ad.Stats.AvgRating,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan ad: %w", err)
		}

		ad.PaymentMethods = make([]int, len(paymentMethods))
		for i, v := range paymentMethods {
			ad.PaymentMethods[i] = int(v)
		}

		ads = append(ads, ad)
	}

	return ads, total, nil
}

func (r *tradeRepository) UpdateAd(ctx context.Context, ad *trading.TradeAd) error {
	query := `
		UPDATE trade_ads
		SET type = $1, crypto_id = $2, fiat_id = $3, price_type = $4,
		    price = $5, floating_markup = $6, min_amount = $7,
		    max_amount = $8, payment_methods = $9, trade_terms = $10,
		    payment_window_minutes = $11, requires_kyc_level = $12,
		    is_public = $13, is_paused = $14, timezone = $15,
		    status = $16, updated_at = NOW()
		WHERE ad_id = $17 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(
		ctx, query,
		ad.Type, ad.CryptoID, ad.FiatID, ad.PriceType, ad.Price,
		ad.FloatingMarkup, ad.MinAmount, ad.MaxAmount, pq.Array(ad.PaymentMethods),
		ad.TradeTerms, ad.PaymentWindowMinutes, ad.RequiresKYCLevel,
		ad.IsPublic, ad.IsPaused, ad.Timezone, ad.Status, ad.AdID,
	)
	if err != nil {
		return fmt.Errorf("update ad: %w", err)
	}
	return nil
}

func (r *tradeRepository) DeleteAd(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE trade_ads SET deleted_at = NOW() WHERE ad_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *tradeRepository) CreateTrade(ctx context.Context, trade *trading.Trade) error {
	query := `
		INSERT INTO trades (
			trade_id, ad_id, buyer_id, seller_id, crypto_id, fiat_id,
			crypto_amount, fiat_amount, exchange_rate, payment_method,
			price_type, agreed_price, status, payment_window_minutes,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			NOW(), NOW()
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		trade.TradeID, trade.AdID, trade.BuyerID, trade.SellerID,
		trade.CryptoID, trade.FiatID, trade.CryptoAmount, trade.FiatAmount,
		trade.ExchangeRate, trade.PaymentMethod, trade.PriceType,
		trade.AgreedPrice, trade.Status, trade.PaymentWindowMinutes,
	).Scan(&trade.CreatedAt, &trade.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create trade: %w", err)
	}
	return nil
}

func (r *tradeRepository) GetTradeByID(ctx context.Context, id uuid.UUID) (*trading.Trade, error) {
	query := `
		SELECT trade_id, ad_id, buyer_id, seller_id, crypto_id, fiat_id,
		       crypto_amount, fiat_amount, exchange_rate, payment_method,
		       price_type, agreed_price, status, payment_window_minutes,
		       created_at, updated_at
		FROM trades
		WHERE trade_id = $1 AND deleted_at IS NULL
	`
	var t trading.Trade
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.TradeID, &t.AdID, &t.BuyerID, &t.SellerID, &t.CryptoID, &t.FiatID,
		&t.CryptoAmount, &t.FiatAmount, &t.ExchangeRate, &t.PaymentMethod,
		&t.PriceType, &t.AgreedPrice, &t.Status, &t.PaymentWindowMinutes,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get trade: %w", err)
	}
	return &t, nil
}

func (r *tradeRepository) ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]trading.Trade, error) {
	var query string
	if role == "buyer" {
		query = `SELECT * FROM trades WHERE buyer_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	} else if role == "seller" {
		query = `SELECT * FROM trades WHERE seller_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	} else {
		query = `SELECT * FROM trades WHERE (buyer_id = $1 OR seller_id = $1) AND deleted_at IS NULL ORDER BY created_at DESC`
	}

	// For simplicity, using SELECT * but in production should list columns
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query trades: %w", err)
	}
	defer rows.Close()

	var trades []trading.Trade
	for rows.Next() {
		var t trading.Trade
		// Manual scan to avoid issues with NULLs or count mismatches
		err := rows.Scan(
			&t.TradeID, &t.AdID, &t.BuyerID, &t.SellerID, &t.CryptoID, &t.FiatID,
			&t.CryptoAmount, &t.FiatAmount, &t.ExchangeRate, &t.PaymentMethod,
			&t.PriceType, &t.AgreedPrice, &t.Status, &t.DisputeID, &t.ChatRoomID,
			&t.StartedAt, &t.PaymentMarkedAt, &t.ReleasedAt, &t.CancelledAt,
			&t.CompletedAt, &t.ExpiredAt, &t.PaymentWindowMinutes,
			&t.IsAutoDisputeTriggered, &t.CancelReason, &t.EscrowTxnHash,
			&t.EscrowContractAddress, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan trade: %w", err)
		}
		trades = append(trades, t)
	}
	return trades, nil
}

func (r *tradeRepository) UpdateTrade(ctx context.Context, t *trading.Trade) error {
	query := `
		UPDATE trades
		SET status = $1, payment_marked_at = $2, released_at = $3,
		    cancelled_at = $4, completed_at = $5, expired_at = $6,
		    cancel_reason = $7, dispute_id = $8, updated_at = NOW()
		WHERE trade_id = $9 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(
		ctx, query,
		t.Status, t.PaymentMarkedAt, t.ReleasedAt, t.CancelledAt,
		t.CompletedAt, t.ExpiredAt, t.CancelReason, t.DisputeID, t.TradeID,
	)
	if err != nil {
		return fmt.Errorf("update trade: %w", err)
	}
	return nil
}

func (r *tradeRepository) CreateTradeMessage(ctx context.Context, msg *trading.TradeMessage) error {
	query := `
		INSERT INTO trade_messages (
			message_id, trade_id, sender_id, message_type, content,
			file_url, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW()
		) RETURNING created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		msg.MessageID, msg.TradeID, msg.SenderID, msg.MessageType,
		msg.Content, msg.FileURL,
	).Scan(&msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func (r *tradeRepository) ListTradeMessages(ctx context.Context, tradeID uuid.UUID) ([]trading.TradeMessage, error) {
	query := `
		SELECT message_id, trade_id, sender_id, message_type, content,
		       file_url, is_read, read_at, created_at
		FROM trade_messages
		WHERE trade_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, tradeID)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer rows.Close()

	var messages []trading.TradeMessage
	for rows.Next() {
		var m trading.TradeMessage
		err := rows.Scan(
			&m.MessageID, &m.TradeID, &m.SenderID, &m.MessageType,
			&m.Content, &m.FileURL, &m.IsRead, &m.ReadAt, &m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
