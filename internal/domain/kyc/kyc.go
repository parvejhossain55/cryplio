package kyc

import (
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
	KYCID           uuid.UUID         `db:"kyc_id" json:"kyc_id"`
	UserID          uuid.UUID         `db:"user_id" json:"user_id"`
	Level           identity.KYCLevel `db:"level" json:"level"`
	Status          KYCStatus         `db:"status" json:"status"` // pending, approved, rejected
	DocumentType    string            `db:"document_type" json:"document_type"`
	DocumentURL     string            `db:"document_url" json:"document_url"`
	SelfieURL       *string           `db:"selfie_url" json:"selfie_url,omitempty"`
	AddressProofURL *string           `db:"address_proof_url" json:"address_proof_url,omitempty"`
	Country         string            `db:"country" json:"country"`
	FullName        string            `db:"full_name" json:"full_name"`
	DateOfBirth     *time.Time        `db:"date_of_birth" json:"date_of_birth,omitempty"`
	Address         *string           `db:"address" json:"address,omitempty"`
	City            *string           `db:"city" json:"city,omitempty"`
	PostalCode      *string           `db:"postal_code" json:"postal_code,omitempty"`
	RejectedReason  *string           `db:"rejected_reason" json:"rejected_reason,omitempty"`
	ReviewedBy      *uuid.UUID        `db:"reviewed_by" json:"reviewed_by,omitempty"` // Admin ID
	ReviewedAt      *time.Time        `db:"reviewed_at" json:"reviewed_at,omitempty"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
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
	k.RejectedReason = &reason
}
