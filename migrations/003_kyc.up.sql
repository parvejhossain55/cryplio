-- ============================================
-- Migration 003: KYC (Know Your Customer) Tables
-- KYC records, document storage, verification status
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS kyc_records (
    kyc_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    level kyc_level NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_front_url VARCHAR(500) NOT NULL,
    document_back_url VARCHAR(500),
    selfie_url VARCHAR(500) NOT NULL,
    provider VARCHAR(30) NOT NULL DEFAULT 'internal',
    provider_reference VARCHAR(255),
    status kyc_level NOT NULL DEFAULT 'level_0',
    rejection_reason TEXT,
    reviewed_by UUID REFERENCES users(user_id),
    reviewed_at TIMESTAMP,
    aml_screened BOOLEAN NOT NULL DEFAULT false,
    aml_check_at TIMESTAMP,
    aml_result JSONB,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kyc_records_user_id ON kyc_records(user_id);
CREATE INDEX idx_kyc_records_status ON kyc_records(status);
CREATE INDEX idx_kyc_records_level ON kyc_records(level);

COMMIT;
