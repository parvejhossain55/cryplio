package identity

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// ─── User Payment Methods ─────────────────────────────────────────────────────

// CreateUserPaymentMethod inserts a new payment method profile and populates
// ID, CreatedAt, and UpdatedAt.
func (r *userRepository) CreateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO user_payment_methods (
			user_id, payment_method_code, display_name,
			account_name, account_number, bank_name,
			is_active, is_default
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`,
		pm.UserID, pm.PaymentMethodCode, pm.DisplayName,
		pm.AccountName, pm.AccountNumber, pm.BankName,
		pm.IsActive, pm.IsDefault,
	).Scan(&pm.ID, &pm.CreatedAt, &pm.UpdatedAt)
}

// GetUserPaymentMethod returns a single payment method by its UUID, or nil.
func (r *userRepository) GetUserPaymentMethod(ctx context.Context, id uuid.UUID) (*UserPaymentMethod, error) {
	var pm UserPaymentMethod
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, payment_method_code, display_name,
		       account_name, account_number, bank_name,
		       is_active, is_default, created_at, updated_at
		FROM user_payment_methods
		WHERE id = $1`, id,
	).Scan(
		&pm.ID, &pm.UserID, &pm.PaymentMethodCode, &pm.DisplayName,
		&pm.AccountName, &pm.AccountNumber, &pm.BankName,
		&pm.IsActive, &pm.IsDefault, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

// GetUserPaymentMethods returns all payment methods for a user, default first.
func (r *userRepository) GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, payment_method_code, display_name,
		       account_name, account_number, bank_name,
		       is_active, is_default, created_at, updated_at
		FROM user_payment_methods
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []UserPaymentMethod
	for rows.Next() {
		var pm UserPaymentMethod
		if err := rows.Scan(
			&pm.ID, &pm.UserID, &pm.PaymentMethodCode, &pm.DisplayName,
			&pm.AccountName, &pm.AccountNumber, &pm.BankName,
			&pm.IsActive, &pm.IsDefault, &pm.CreatedAt, &pm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		methods = append(methods, pm)
	}
	return methods, nil
}

// GetUserPaymentMethodsByUserID is an alias for GetUserPaymentMethods.
func (r *userRepository) GetUserPaymentMethodsByUserID(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error) {
	return r.GetUserPaymentMethods(ctx, userID)
}

// UpdateUserPaymentMethod persists changed fields and refreshes UpdatedAt.
func (r *userRepository) UpdateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error {
	return r.db.QueryRowContext(ctx, `
		UPDATE user_payment_methods
		SET display_name = $1, account_name = $2, account_number = $3,
		    bank_name = $4, is_active = $5, is_default = $6,
		    updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at`,
		pm.DisplayName, pm.AccountName, pm.AccountNumber,
		pm.BankName, pm.IsActive, pm.IsDefault,
		pm.ID, pm.UserID,
	).Scan(&pm.UpdatedAt)
}

// DeleteUserPaymentMethod removes a payment method by its UUID.
func (r *userRepository) DeleteUserPaymentMethod(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_payment_methods WHERE id = $1`, id)
	return err
}

// SetDefaultUserPaymentMethod atomically clears the previous default and sets
// a new one for the user, using a transaction.
func (r *userRepository) SetDefaultUserPaymentMethod(ctx context.Context, userID, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err = tx.ExecContext(ctx,
		`UPDATE user_payment_methods SET is_default = false WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx,
		`UPDATE user_payment_methods SET is_default = true WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return err
	}
	return tx.Commit()
}
