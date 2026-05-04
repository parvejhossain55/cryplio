package kyc

import (
	"context"
	"database/sql"
	"fmt"

	domainkyc "cryplio/internal/domain/kyc"

	"github.com/google/uuid"
)

type kycRepository struct {
	db *sql.DB
}

func NewKYCRepository(db *sql.DB) domainkyc.KYCRepository {
	return &kycRepository{db: db}
}

func (r *kycRepository) Create(ctx context.Context, record *domainkyc.KYCRecord) error {
	query := `
		INSERT INTO kyc_records (
			user_id, level, status, document_type, document_front_url, 
			document_back_url, selfie_url, provider, provider_reference,
			submitted_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()
		) RETURNING kyc_id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		record.UserID, record.Level, record.Status, record.DocumentType,
		record.DocumentFrontURL, record.DocumentBackURL, record.SelfieURL,
		record.Provider, record.ProviderReference, record.SubmittedAt,
	).Scan(&record.KYCID, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		return fmt.Errorf("insert kyc record: %w", err)
	}
	return nil
}

func (r *kycRepository) GetByID(ctx context.Context, kycID uuid.UUID) (*domainkyc.KYCRecord, error) {
	query := `
		SELECT kyc_id, user_id, level, status, document_type, document_front_url,
		       document_back_url, selfie_url, provider, provider_reference,
		       rejection_reason, reviewed_by, reviewed_at, aml_screened,
		       aml_check_at, aml_result, submitted_at, created_at, updated_at
		FROM kyc_records
		WHERE kyc_id = $1
	`
	var k domainkyc.KYCRecord
	var amlResult []byte
	err := r.db.QueryRowContext(ctx, query, kycID).Scan(
		&k.KYCID, &k.UserID, &k.Level, &k.Status, &k.DocumentType, &k.DocumentFrontURL,
		&k.DocumentBackURL, &k.SelfieURL, &k.Provider, &k.ProviderReference,
		&k.RejectionReason, &k.ReviewedBy, &k.ReviewedAt, &k.AMLScreened,
		&k.AMLCheckAt, &amlResult, &k.SubmittedAt, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query kyc record by id: %w", err)
	}
	// Note: aml_result is JSONB, could be unmarshaled if needed
	k.AMLResult = amlResult
	return &k, nil
}

func (r *kycRepository) GetLatestByUserID(ctx context.Context, userID uuid.UUID) (*domainkyc.KYCRecord, error) {
	query := `
		SELECT kyc_id, user_id, level, status, document_type, document_front_url,
		       document_back_url, selfie_url, provider, provider_reference,
		       rejection_reason, reviewed_by, reviewed_at, aml_screened,
		       aml_check_at, aml_result, submitted_at, created_at, updated_at
		FROM kyc_records
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var k domainkyc.KYCRecord
	var amlResult []byte
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&k.KYCID, &k.UserID, &k.Level, &k.Status, &k.DocumentType, &k.DocumentFrontURL,
		&k.DocumentBackURL, &k.SelfieURL, &k.Provider, &k.ProviderReference,
		&k.RejectionReason, &k.ReviewedBy, &k.ReviewedAt, &k.AMLScreened,
		&k.AMLCheckAt, &amlResult, &k.SubmittedAt, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query latest kyc record by user id: %w", err)
	}
	k.AMLResult = amlResult
	return &k, nil
}

func (r *kycRepository) Update(ctx context.Context, k *domainkyc.KYCRecord) error {
	query := `
		UPDATE kyc_records
		SET level = $1, status = $2, rejection_reason = $3, reviewed_by = $4,
		    reviewed_at = $5, aml_screened = $6, aml_check_at = $7, aml_result = $8,
		    updated_at = NOW()
		WHERE kyc_id = $9
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		k.Level, k.Status, k.RejectionReason, k.ReviewedBy,
		k.ReviewedAt, k.AMLScreened, k.AMLCheckAt, k.AMLResult,
		k.KYCID,
	).Scan(&k.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update kyc record: %w", err)
	}
	return nil
}

func (r *kycRepository) GetPending(ctx context.Context, limit, offset int) ([]domainkyc.KYCRecord, error) {
	query := `
		SELECT kyc_id, user_id, level, status, document_type, document_front_url,
		       document_back_url, selfie_url, provider, provider_reference,
		       rejection_reason, reviewed_by, reviewed_at, aml_screened,
		       aml_check_at, aml_result, submitted_at, created_at, updated_at
		FROM kyc_records
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query pending kyc records: %w", err)
	}
	defer rows.Close()

	var records []domainkyc.KYCRecord
	for rows.Next() {
		var k domainkyc.KYCRecord
		var amlResult []byte
		err := rows.Scan(
			&k.KYCID, &k.UserID, &k.Level, &k.Status, &k.DocumentType, &k.DocumentFrontURL,
			&k.DocumentBackURL, &k.SelfieURL, &k.Provider, &k.ProviderReference,
			&k.RejectionReason, &k.ReviewedBy, &k.ReviewedAt, &k.AMLScreened,
			&k.AMLCheckAt, &amlResult, &k.SubmittedAt, &k.CreatedAt, &k.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan kyc record: %w", err)
		}
		k.AMLResult = amlResult
		records = append(records, k)
	}
	return records, nil
}
