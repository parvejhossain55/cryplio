package dto

import "time"

// KYCSubmitRequest represents the request to submit KYC documents
type KYCSubmitRequest struct {
	Level            string  `json:"level" binding:"required,oneof=level_1 level_2 level_3"`
	DocumentType     string  `json:"document_type" binding:"required"`
	DocumentFrontURL string  `json:"document_front_url" binding:"required,url"`
	DocumentBackURL  *string `json:"document_back_url,omitempty" binding:"omitempty,url"`
	SelfieURL        string  `json:"selfie_url" binding:"required,url"`
}

// KYCRecordResponse represents the public projection of a KYC record
type KYCRecordResponse struct {
	KYCID            string     `json:"kyc_id"`
	UserID           string     `json:"user_id"`
	Level            string     `json:"level"`
	Status           string     `json:"status"`
	DocumentType     string     `json:"document_type"`
	DocumentFrontURL string     `json:"document_front_url"`
	DocumentBackURL  *string    `json:"document_back_url,omitempty"`
	SelfieURL        string     `json:"selfie_url"`
	Provider         string     `json:"provider"`
	RejectionReason  *string    `json:"rejection_reason,omitempty"`
	AMLScreened      bool       `json:"aml_screened"`
	SubmittedAt      time.Time  `json:"submitted_at"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
}

// KYCVerifyRequest represents the admin request to verify a KYC record
type KYCVerifyRequest struct {
	Approved        bool   `json:"approved"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}
