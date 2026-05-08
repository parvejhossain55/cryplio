package platform

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cryplio/internal/domain/platform"
)

// ─── Payment Methods ──────────────────────────────────────────────────────────

func (r *platformRepository) CreatePaymentMethod(ctx context.Context, m *platform.PaymentMethod) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO payment_methods (
			code, name, category, icon_url, description,
			is_active, sort_order, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`,
		m.Code, m.Name, m.Category, m.IconURL, m.Description,
		m.IsActive, m.SortOrder, time.Now(),
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *platformRepository) GetPaymentMethod(ctx context.Context, id int) (*platform.PaymentMethod, error) {
	m := &platform.PaymentMethod{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, category, icon_url, description,
		       is_active, sort_order, created_at
		FROM payment_methods WHERE id = $1`, id,
	).Scan(&m.ID, &m.Code, &m.Name, &m.Category, &m.IconURL, &m.Description,
		&m.IsActive, &m.SortOrder, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment method not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get payment method: %w", err)
	}
	return m, nil
}

func (r *platformRepository) GetPaymentMethods(ctx context.Context, activeOnly bool, limit, offset int) ([]*platform.PaymentMethod, int, error) {
	base, args := buildFilterQuery("payment_methods", "is_active", activeOnly)

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+base, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count payment methods: %w", err)
	}

	q, a := buildPagedQuery(
		`SELECT id, code, name, category, icon_url, description,
		        is_active, sort_order, created_at
		 FROM `+base+" ORDER BY sort_order, name",
		args, limit, offset)

	rows, err := r.db.QueryContext(ctx, q, a...)
	if err != nil {
		return nil, 0, fmt.Errorf("get payment methods: %w", err)
	}
	defer rows.Close()

	var methods []*platform.PaymentMethod
	for rows.Next() {
		m := &platform.PaymentMethod{}
		if err := rows.Scan(&m.ID, &m.Code, &m.Name, &m.Category, &m.IconURL, &m.Description,
			&m.IsActive, &m.SortOrder, &m.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan payment method: %w", err)
		}
		methods = append(methods, m)
	}
	return methods, total, rows.Err()
}

func (r *platformRepository) UpdatePaymentMethod(ctx context.Context, m *platform.PaymentMethod) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE payment_methods SET
			code = $1, name = $2, category = $3, icon_url = $4, description = $5,
			is_active = $6, sort_order = $7
		WHERE id = $8`,
		m.Code, m.Name, m.Category, m.IconURL, m.Description,
		m.IsActive, m.SortOrder, m.ID)
	if err != nil {
		return fmt.Errorf("update payment method: %w", err)
	}
	return nil
}

func (r *platformRepository) DeletePaymentMethod(ctx context.Context, id int) error {
	return deleteByID(r.db, ctx, "payment_methods", id, "payment method")
}
