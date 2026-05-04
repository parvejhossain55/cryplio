package kyc

import (
	"context"
	"time"

	identity "cryplio/internal/domain/identity"

	"github.com/google/uuid"
)

// KYCStatus represents the status of a KYC record
type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusApproved KYCStatus = "approved"
	KYCStatusRejected KYCStatus = "rejected"
)

// KYCRecord represents a KYC verification submission
type KYCRecord struct {
	KYCID             uuid.UUID         `db:"kyc_id" json:"kyc_id"`
	UserID            uuid.UUID         `db:"user_id" json:"user_id"`
	Level             identity.KYCLevel `db:"level" json:"level"`
	Status            KYCStatus         `db:"status" json:"status"` // pending, approved, rejected
	DocumentType      string            `db:"document_type" json:"document_type"`
	DocumentFrontURL  string            `db:"document_front_url" json:"document_front_url"`
	DocumentBackURL   *string           `db:"document_back_url" json:"document_back_url,omitempty"`
	SelfieURL         string            `db:"selfie_url" json:"selfie_url"`
	Provider          string            `db:"provider" json:"provider"` // sumsub, internal
	ProviderReference *string           `db:"provider_reference" json:"provider_reference,omitempty"`
	RejectionReason   *string           `db:"rejection_reason" json:"rejection_reason,omitempty"`
	ReviewedBy        *uuid.UUID        `db:"reviewed_by" json:"reviewed_by,omitempty"` // Admin ID
	ReviewedAt        *time.Time        `db:"reviewed_at" json:"reviewed_at,omitempty"`
	AMLScreened       bool              `db:"aml_screened" json:"aml_screened"`
	AMLCheckAt        *time.Time        `db:"aml_check_at" json:"aml_check_at,omitempty"`
	AMLResult         interface{}       `db:"aml_result" json:"aml_result,omitempty"`
	SubmittedAt       time.Time         `db:"submitted_at" json:"submitted_at"`
	CreatedAt         time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at" json:"updated_at"`
}

// IsPending checks if the KYC record is pending review
func (k *KYCRecord) IsPending() bool {
	return k.Status == KYCStatusPending
}

// IsApproved checks if the KYC record is approved
func (k *KYCRecord) IsApproved() bool {
	return k.Status == KYCStatusApproved
}

// IsRejected checks if the KYC record is rejected
func (k *KYCRecord) IsRejected() bool {
	return k.Status == KYCStatusRejected
}

// Approve approves the KYC record
func (k *KYCRecord) Approve(adminID uuid.UUID) {
	k.Status = KYCStatusApproved
	now := time.Now()
	k.ReviewedBy = &adminID
	k.ReviewedAt = &now
}

// Reject rejects the KYC record with a reason
func (k *KYCRecord) Reject(adminID uuid.UUID, reason string) {
	k.Status = KYCStatusRejected
	now := time.Now()
	k.ReviewedBy = &adminID
	k.ReviewedAt = &now
	k.RejectionReason = &reason
}

// SetAMLResult sets the AML check result
func (k *KYCRecord) SetAMLResult(result interface{}) {
	k.AMLScreened = true
	now := time.Now()
	k.AMLCheckAt = &now
	k.AMLResult = result
}

// KYCRepository defines the behavior of a KYC record data store
type KYCRepository interface {
	Create(ctx context.Context, record *KYCRecord) error
	GetByID(ctx context.Context, kycID uuid.UUID) (*KYCRecord, error)
	GetLatestByUserID(ctx context.Context, userID uuid.UUID) (*KYCRecord, error)
	Update(ctx context.Context, record *KYCRecord) error
	GetPending(ctx context.Context, limit, offset int) ([]KYCRecord, error)
}
