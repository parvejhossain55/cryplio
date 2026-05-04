package platform

import (
	"context"
	"cryplio/internal/domain/platform"
	"database/sql"
	"fmt"
	"time"
)

type platformRepository struct {
	db *sql.DB
}

// NewPlatformRepository creates a new postgres platform repository
func NewPlatformRepository(db *sql.DB) platform.PlatformRepository {
	return &platformRepository{db: db}
}

// Crypto Assets
func (r *platformRepository) CreateCryptoAsset(ctx context.Context, asset *platform.CryptoAsset) error {
	query := `
		INSERT INTO crypto_assets (
			symbol, name, blockchain, contract_address, decimals, min_confirmation, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		asset.Symbol, asset.Name, asset.Blockchain, asset.ContractAddress,
		asset.Decimals, asset.MinConfirmation, asset.IsActive, now, now,
	).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create crypto asset: %w", err)
	}
	return nil
}

func (r *platformRepository) GetCryptoAsset(ctx context.Context, id int) (*platform.CryptoAsset, error) {
	query := `
		SELECT id, symbol, name, blockchain, contract_address, decimals, min_confirmation, is_active, created_at, updated_at
		FROM crypto_assets WHERE id = $1
	`

	asset := &platform.CryptoAsset{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&asset.ID, &asset.Symbol, &asset.Name, &asset.Blockchain, &asset.ContractAddress,
		&asset.Decimals, &asset.MinConfirmation, &asset.IsActive, &asset.CreatedAt, &asset.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("crypto asset not found")
		}
		return nil, fmt.Errorf("get crypto asset: %w", err)
	}
	return asset, nil
}

func (r *platformRepository) GetCryptoAssets(ctx context.Context, activeOnly bool, limit, offset int) ([]*platform.CryptoAsset, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM crypto_assets`
	countArgs := []interface{}{}

	if activeOnly {
		countQuery += " WHERE is_active = $1"
		countArgs = append(countArgs, true)
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("get crypto assets count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, symbol, name, blockchain, contract_address, decimals, min_confirmation, is_active, created_at, updated_at
		FROM crypto_assets
	`
	args := []interface{}{}

	if activeOnly {
		query += " WHERE is_active = $1"
		args = append(args, true)
	}

	query += " ORDER BY symbol"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get crypto assets: %w", err)
	}
	defer rows.Close()

	var assets []*platform.CryptoAsset
	for rows.Next() {
		asset := &platform.CryptoAsset{}
		err := rows.Scan(
			&asset.ID, &asset.Symbol, &asset.Name, &asset.Blockchain, &asset.ContractAddress,
			&asset.Decimals, &asset.MinConfirmation, &asset.IsActive, &asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan crypto asset: %w", err)
		}
		assets = append(assets, asset)
	}

	return assets, total, rows.Err()
}

func (r *platformRepository) UpdateCryptoAsset(ctx context.Context, asset *platform.CryptoAsset) error {
	query := `
		UPDATE crypto_assets SET
			symbol = $1, name = $2, blockchain = $3, contract_address = $4,
			decimals = $5, min_confirmation = $6, is_active = $7, updated_at = $8
		WHERE id = $9
		RETURNING updated_at
	`

	asset.UpdatedAt = time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		asset.Symbol, asset.Name, asset.Blockchain, asset.ContractAddress,
		asset.Decimals, asset.MinConfirmation, asset.IsActive, asset.UpdatedAt, asset.ID,
	).Scan(&asset.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update crypto asset: %w", err)
	}
	return nil
}

func (r *platformRepository) DeleteCryptoAsset(ctx context.Context, id int) error {
	query := `DELETE FROM crypto_assets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete crypto asset: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("crypto asset not found")
	}

	return nil
}

// Fiat Currencies
func (r *platformRepository) CreateFiatCurrency(ctx context.Context, currency *platform.FiatCurrency) error {
	query := `
		INSERT INTO fiat_currencies (code, name, symbol, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		currency.Code, currency.Name, currency.Symbol, currency.IsActive, now,
	).Scan(&currency.ID, &currency.CreatedAt)

	if err != nil {
		return fmt.Errorf("create fiat currency: %w", err)
	}
	return nil
}

func (r *platformRepository) GetFiatCurrency(ctx context.Context, id int) (*platform.FiatCurrency, error) {
	query := `
		SELECT id, code, name, symbol, is_active, created_at
		FROM fiat_currencies WHERE id = $1
	`

	currency := &platform.FiatCurrency{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&currency.ID, &currency.Code, &currency.Name, &currency.Symbol,
		&currency.IsActive, &currency.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("fiat currency not found")
		}
		return nil, fmt.Errorf("get fiat currency: %w", err)
	}
	return currency, nil
}

func (r *platformRepository) GetFiatCurrencies(ctx context.Context, activeOnly bool, limit, offset int) ([]*platform.FiatCurrency, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM fiat_currencies`
	countArgs := []interface{}{}

	if activeOnly {
		countQuery += " WHERE is_active = $1"
		countArgs = append(countArgs, true)
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("get fiat currencies count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, code, name, symbol, is_active, created_at
		FROM fiat_currencies
	`
	args := []interface{}{}

	if activeOnly {
		query += " WHERE is_active = $1"
		args = append(args, true)
	}

	query += " ORDER BY code"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get fiat currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*platform.FiatCurrency
	for rows.Next() {
		currency := &platform.FiatCurrency{}
		err := rows.Scan(
			&currency.ID, &currency.Code, &currency.Name, &currency.Symbol,
			&currency.IsActive, &currency.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan fiat currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, total, rows.Err()
}

func (r *platformRepository) UpdateFiatCurrency(ctx context.Context, currency *platform.FiatCurrency) error {
	query := `
		UPDATE fiat_currencies SET
			code = $1, name = $2, symbol = $3, is_active = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx, query,
		currency.Code, currency.Name, currency.Symbol, currency.IsActive, currency.ID,
	)

	if err != nil {
		return fmt.Errorf("update fiat currency: %w", err)
	}
	return nil
}

func (r *platformRepository) DeleteFiatCurrency(ctx context.Context, id int) error {
	query := `DELETE FROM fiat_currencies WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete fiat currency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("fiat currency not found")
	}

	return nil
}

// Payment Methods
func (r *platformRepository) CreatePaymentMethod(ctx context.Context, method *platform.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods (
			code, name, category, icon_url, description, is_active, sort_order, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		method.Code, method.Name, method.Category, method.IconURL, method.Description,
		method.IsActive, method.SortOrder, now,
	).Scan(&method.ID, &method.CreatedAt)

	if err != nil {
		return fmt.Errorf("create payment method: %w", err)
	}
	return nil
}

func (r *platformRepository) GetPaymentMethod(ctx context.Context, id int) (*platform.PaymentMethod, error) {
	query := `
		SELECT id, code, name, category, icon_url, description, is_active, sort_order, created_at
		FROM payment_methods WHERE id = $1
	`

	method := &platform.PaymentMethod{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&method.ID, &method.Code, &method.Name, &method.Category, &method.IconURL,
		&method.Description, &method.IsActive, &method.SortOrder, &method.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("get payment method: %w", err)
	}
	return method, nil
}

func (r *platformRepository) GetPaymentMethods(ctx context.Context, activeOnly bool, limit, offset int) ([]*platform.PaymentMethod, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM payment_methods`
	countArgs := []interface{}{}

	if activeOnly {
		countQuery += " WHERE is_active = $1"
		countArgs = append(countArgs, true)
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("get payment methods count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, code, name, category, icon_url, description, is_active, sort_order, created_at
		FROM payment_methods
	`
	args := []interface{}{}

	if activeOnly {
		query += " WHERE is_active = $1"
		args = append(args, true)
	}

	query += " ORDER BY sort_order, name"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get payment methods: %w", err)
	}
	defer rows.Close()

	var methods []*platform.PaymentMethod
	for rows.Next() {
		method := &platform.PaymentMethod{}
		err := rows.Scan(
			&method.ID, &method.Code, &method.Name, &method.Category, &method.IconURL,
			&method.Description, &method.IsActive, &method.SortOrder, &method.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan payment method: %w", err)
		}
		methods = append(methods, method)
	}

	return methods, total, rows.Err()
}

func (r *platformRepository) UpdatePaymentMethod(ctx context.Context, method *platform.PaymentMethod) error {
	query := `
		UPDATE payment_methods SET
			code = $1, name = $2, category = $3, icon_url = $4, description = $5,
			is_active = $6, sort_order = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx, query,
		method.Code, method.Name, method.Category, method.IconURL, method.Description,
		method.IsActive, method.SortOrder, method.ID,
	)

	if err != nil {
		return fmt.Errorf("update payment method: %w", err)
	}
	return nil
}

func (r *platformRepository) DeletePaymentMethod(ctx context.Context, id int) error {
	query := `DELETE FROM payment_methods WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete payment method: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}
