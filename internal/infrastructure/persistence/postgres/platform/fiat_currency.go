package platform

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cryplio/internal/domain/platform"
)

// ─── Fiat Currencies ──────────────────────────────────────────────────────────

func (r *platformRepository) CreateFiatCurrency(ctx context.Context, c *platform.FiatCurrency) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO fiat_currencies (code, name, symbol, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		c.Code, c.Name, c.Symbol, c.IsActive, time.Now(),
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *platformRepository) GetFiatCurrency(ctx context.Context, id int) (*platform.FiatCurrency, error) {
	c := &platform.FiatCurrency{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, symbol, is_active, created_at
		FROM fiat_currencies WHERE id = $1`, id,
	).Scan(&c.ID, &c.Code, &c.Name, &c.Symbol, &c.IsActive, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fiat currency not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get fiat currency: %w", err)
	}
	return c, nil
}

func (r *platformRepository) GetFiatCurrencies(ctx context.Context, activeOnly bool, limit, offset int) ([]*platform.FiatCurrency, int, error) {
	base, args := buildFilterQuery("fiat_currencies", "is_active", activeOnly)

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+base, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count fiat currencies: %w", err)
	}

	q, a := buildPagedQuery(
		`SELECT id, code, name, symbol, is_active, created_at FROM `+base+" ORDER BY code",
		args, limit, offset)

	rows, err := r.db.QueryContext(ctx, q, a...)
	if err != nil {
		return nil, 0, fmt.Errorf("get fiat currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*platform.FiatCurrency
	for rows.Next() {
		c := &platform.FiatCurrency{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Symbol, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan fiat currency: %w", err)
		}
		currencies = append(currencies, c)
	}
	return currencies, total, rows.Err()
}

func (r *platformRepository) UpdateFiatCurrency(ctx context.Context, c *platform.FiatCurrency) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE fiat_currencies SET code = $1, name = $2, symbol = $3, is_active = $4
		WHERE id = $5`,
		c.Code, c.Name, c.Symbol, c.IsActive, c.ID)
	if err != nil {
		return fmt.Errorf("update fiat currency: %w", err)
	}
	return nil
}

func (r *platformRepository) DeleteFiatCurrency(ctx context.Context, id int) error {
	return deleteByID(r.db, ctx, "fiat_currencies", id, "fiat currency")
}
