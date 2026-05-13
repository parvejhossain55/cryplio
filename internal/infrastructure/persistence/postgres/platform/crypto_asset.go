package platform

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cryplio/internal/domain/platform"
)

// ─── Crypto Assets ────────────────────────────────────────────────────────────

func (r *platformRepository) CreateCryptoAsset(ctx context.Context, asset *platform.CryptoAsset) error {
	now := time.Now()
	return r.db.QueryRowContext(ctx, `
		INSERT INTO crypto_assets (
			symbol, name, blockchain, contract_address,
			decimals, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`,
		asset.Symbol, asset.Name, asset.Blockchain, asset.ContractAddress,
		asset.Decimals, asset.IsActive, now, now,
	).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)
}

func (r *platformRepository) GetCryptoAsset(ctx context.Context, id int) (*platform.CryptoAsset, error) {
	asset := &platform.CryptoAsset{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, symbol, name, blockchain, contract_address,
		       decimals, is_active, created_at, updated_at
		FROM crypto_assets WHERE id = $1`, id,
	).Scan(&asset.ID, &asset.Symbol, &asset.Name, &asset.Blockchain, &asset.ContractAddress,
		&asset.Decimals, &asset.IsActive, &asset.CreatedAt, &asset.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("crypto asset not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get crypto asset: %w", err)
	}
	return asset, nil
}

func (r *platformRepository) GetCryptoAssets(ctx context.Context, activeOnly bool, searchQuery string, limit, offset int) ([]*platform.CryptoAsset, int, error) {
	base, args := buildFilterQuery("crypto_assets", "is_active", activeOnly, searchQuery, []string{"symbol", "name", "blockchain"})

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+base, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count crypto assets: %w", err)
	}

	q, a := buildPagedQuery(
		`SELECT id, symbol, name, blockchain, contract_address,
		        decimals, is_active, created_at, updated_at
		 FROM `+base+" ORDER BY symbol", args, limit, offset)

	rows, err := r.db.QueryContext(ctx, q, a...)
	if err != nil {
		return nil, 0, fmt.Errorf("get crypto assets: %w", err)
	}
	defer rows.Close()

	var assets []*platform.CryptoAsset
	for rows.Next() {
		a := &platform.CryptoAsset{}
		if err := rows.Scan(&a.ID, &a.Symbol, &a.Name, &a.Blockchain, &a.ContractAddress,
			&a.Decimals, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan crypto asset: %w", err)
		}
		assets = append(assets, a)
	}
	return assets, total, rows.Err()
}

func (r *platformRepository) UpdateCryptoAsset(ctx context.Context, asset *platform.CryptoAsset) error {
	asset.UpdatedAt = time.Now()
	return r.db.QueryRowContext(ctx, `
		UPDATE crypto_assets SET
			symbol = $1, name = $2, blockchain = $3, contract_address = $4,
			decimals = $5, is_active = $6, updated_at = $7
		WHERE id = $8
		RETURNING updated_at`,
		asset.Symbol, asset.Name, asset.Blockchain, asset.ContractAddress,
		asset.Decimals, asset.IsActive, asset.UpdatedAt, asset.ID,
	).Scan(&asset.UpdatedAt)
}

func (r *platformRepository) DeleteCryptoAsset(ctx context.Context, id int) error {
	return deleteByID(r.db, ctx, "crypto_assets", id, "crypto asset")
}
