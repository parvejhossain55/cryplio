package trading

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"cryplio/internal/domain/trading"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ─── Trade Advertisement CRUD ─────────────────────────────────────────────────

func (r *tradeRepository) CreateAd(ctx context.Context, ad *trading.TradeAd) error {
	query := `
		INSERT INTO trade_ads (
			ad_id, user_id, type, crypto_id, fiat_id, price_type, price,
			floating_markup, min_amount, max_amount, payment_methods,
			trade_terms, payment_window_minutes,
			is_public, is_paused, timezone, status,
			published_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, NOW(), NOW(), NOW()
		) RETURNING created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		ad.AdID, ad.UserID, ad.Type, ad.CryptoID, ad.FiatID, ad.PriceType,
		ad.Price, ad.FloatingMarkup, ad.MinAmount, ad.MaxAmount,
		pq.Array(ad.PaymentMethods), ad.TradeTerms, ad.PaymentWindowMinutes,
		ad.IsPublic, ad.IsPaused, ad.Timezone, ad.Status,
	).Scan(&ad.CreatedAt, &ad.UpdatedAt)
}

func (r *tradeRepository) GetAdByID(ctx context.Context, id uuid.UUID) (*trading.TradeAd, error) {
	query := `
		SELECT ad_id, user_id, type, crypto_id, fiat_id, price_type, price,
		       floating_markup, min_amount, max_amount, payment_methods,
		       trade_terms, payment_window_minutes,
		       is_public, is_paused, timezone, status, published_at,
		       created_at, updated_at
		FROM trade_ads
		WHERE ad_id = $1 AND deleted_at IS NULL
	`
	var ad trading.TradeAd
	var pms []int64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ad.AdID, &ad.UserID, &ad.Type, &ad.CryptoID, &ad.FiatID, &ad.PriceType,
		&ad.Price, &ad.FloatingMarkup, &ad.MinAmount, &ad.MaxAmount,
		pq.Array(&pms), &ad.TradeTerms, &ad.PaymentWindowMinutes,
		&ad.IsPublic, &ad.IsPaused, &ad.Timezone,
		&ad.Status, &ad.PublishedAt, &ad.CreatedAt, &ad.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get ad: %w", err)
	}
	ad.PaymentMethods = toIntSlice(pms)
	return &ad, nil
}

// ListAds returns a filtered, paginated list of ads joined with user + asset info.
func (r *tradeRepository) ListAds(ctx context.Context, filter trading.AdFilter) ([]trading.TradeAd, int, error) {
	var args []interface{}
	var conds []string
	ph := 1 // placeholder counter

	addArg := func(cond string, val interface{}) {
		conds = append(conds, fmt.Sprintf(cond, ph))
		args = append(args, val)
		ph++
	}

	if filter.Type != nil {
		addArg("a.type = $%d", *filter.Type)
	}
	if filter.CryptoID != nil {
		addArg("a.crypto_id = $%d", *filter.CryptoID)
	}
	if filter.FiatID != nil {
		addArg("a.fiat_id = $%d", *filter.FiatID)
	}
	if filter.FiatCode != nil && *filter.FiatCode != "" {
		addArg("fc.code = $%d", *filter.FiatCode)
	}
	if len(filter.PaymentMethods) > 0 {
		placeholders := make([]string, len(filter.PaymentMethods))
		for i, pm := range filter.PaymentMethods {
			placeholders[i] = fmt.Sprintf("$%d", ph)
			args = append(args, pm)
			ph++
		}
		conds = append(conds, fmt.Sprintf("a.payment_methods && ARRAY[%s]::int[]", strings.Join(placeholders, ",")))
	}
	if filter.UserID != nil {
		addArg("a.user_id = $%d", *filter.UserID)
	}
	if filter.Status != nil {
		addArg("a.status = $%d", *filter.Status)
	} else {
		conds = append(conds, "a.status = 'active'")
	}
	conds = append(conds, "a.deleted_at IS NULL")
	if filter.UserID == nil {
		conds = append(conds, "a.is_public = true", "a.is_paused = false")
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	orderBy := "a.published_at DESC"
	switch filter.SortBy {
	case "best_price":
		orderBy = "a.price ASC"
	case "trade_count":
		orderBy = "COALESCE(us.total_trades, 0) DESC"
	}

	joins := `
		FROM trade_ads a
		JOIN users u ON a.user_id = u.user_id
		LEFT JOIN user_stats us ON u.user_id = us.user_id
		JOIN crypto_assets ca ON a.crypto_id = ca.id
		JOIN fiat_currencies fc ON a.fiat_id = fc.id
	`

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+joins+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	dataArgs := append(args, limit, filter.Offset)
	dataQuery := `
		SELECT a.ad_id, a.user_id, a.type, a.crypto_id, a.fiat_id, a.price_type, a.price,
		       a.floating_markup, a.min_amount, a.max_amount, a.payment_methods,
		       a.trade_terms, a.payment_window_minutes,
		       a.is_public, a.is_paused, a.timezone, a.status, a.published_at,
		       a.created_at, a.updated_at,
		       u.username, u.avatar_url, u.last_seen_at,
		       COALESCE(us.total_trades, 0), COALESCE(us.avg_rating, 0),
		       ca.symbol, fc.symbol
	` + joins + where + fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, ph, ph+1)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query ads: %w", err)
	}
	defer rows.Close()

	var ads []trading.TradeAd
	for rows.Next() {
		var ad trading.TradeAd
		var pms []int64
		var username, avatarURL, cryptoSymbol, fiatSymbol sql.NullString
		var lastSeen pq.NullTime
		var totalTrades sql.NullInt64
		var avgRating sql.NullFloat64

		if err := rows.Scan(
			&ad.AdID, &ad.UserID, &ad.Type, &ad.CryptoID, &ad.FiatID, &ad.PriceType,
			&ad.Price, &ad.FloatingMarkup, &ad.MinAmount, &ad.MaxAmount,
			pq.Array(&pms), &ad.TradeTerms, &ad.PaymentWindowMinutes,
			&ad.IsPublic, &ad.IsPaused, &ad.Timezone,
			&ad.Status, &ad.PublishedAt, &ad.CreatedAt, &ad.UpdatedAt,
			&username, &avatarURL, &lastSeen, &totalTrades, &avgRating,
			&cryptoSymbol, &fiatSymbol,
		); err != nil {
			return nil, 0, fmt.Errorf("scan ad: %w", err)
		}
		if username.Valid {
			ad.Username = username.String
		}
		if avatarURL.Valid {
			ad.UserAvatar = avatarURL.String
		}
		if lastSeen.Valid {
			ad.UserLastSeen = &lastSeen.Time
		}
		if totalTrades.Valid {
			ad.UserTrades = int(totalTrades.Int64)
		}
		if avgRating.Valid {
			ad.UserRating = avgRating.Float64
		}
		if cryptoSymbol.Valid {
			ad.CryptoSymbol = cryptoSymbol.String
		}
		if fiatSymbol.Valid {
			ad.FiatSymbol = fiatSymbol.String
		}
		ad.PaymentMethods = toIntSlice(pms)
		ads = append(ads, ad)
	}
	return ads, total, nil
}

func (r *tradeRepository) UpdateAd(ctx context.Context, ad *trading.TradeAd) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE trade_ads
		SET type = $1, crypto_id = $2, fiat_id = $3, price_type = $4,
		    price = $5, floating_markup = $6, min_amount = $7, max_amount = $8,
		    payment_methods = $9, trade_terms = $10, payment_window_minutes = $11,
		    is_public = $12, is_paused = $13, timezone = $14,
		    status = $15, updated_at = NOW()
		WHERE ad_id = $16 AND deleted_at IS NULL`,
		ad.Type, ad.CryptoID, ad.FiatID, ad.PriceType, ad.Price,
		ad.FloatingMarkup, ad.MinAmount, ad.MaxAmount, pq.Array(ad.PaymentMethods),
		ad.TradeTerms, ad.PaymentWindowMinutes, ad.IsPublic, ad.IsPaused, ad.Timezone,
		ad.Status, ad.AdID,
	)
	if err != nil {
		return fmt.Errorf("update ad: %w", err)
	}
	return nil
}

func (r *tradeRepository) DeleteAd(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE trade_ads SET deleted_at = NOW() WHERE ad_id = $1`, id)
	return err
}

// toIntSlice converts []int64 to []int (pq.Array scans as int64).
func toIntSlice(s []int64) []int {
	out := make([]int, len(s))
	for i, v := range s {
		out[i] = int(v)
	}
	return out
}
