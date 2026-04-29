-- ============================================
-- Migration 013: Monitoring & Rate Limiting
-- System monitoring, rate limit tracking
-- ============================================

BEGIN;

-- Rate limiting (per IP or user_id)
CREATE TABLE IF NOT EXISTS rate_limit_counts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) NOT NULL,
    endpoint VARCHAR(200) NOT NULL,
    count INT NOT NULL DEFAULT 1,
    window_start TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '1 minute')
);

CREATE INDEX idx_rate_limit ON rate_limit_counts(identifier, endpoint);
CREATE INDEX idx_rate_limit_expires ON rate_limit_counts(expires_at);

-- Login attempts tracking (for brute force protection, SRS FR-116)
CREATE TABLE IF NOT EXISTS login_attempts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT false,
    reason VARCHAR(50), -- 'wrong_password', 'account_locked', etc.
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_ip ON login_attempts(ip_address);
CREATE INDEX idx_login_attempts_created ON login_attempts(created_at DESC);

-- API request logs (for analytics and debugging)
CREATE TABLE IF NOT EXISTS api_request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint VARCHAR(200) NOT NULL,
    method VARCHAR(10) NOT NULL,
    user_id UUID REFERENCES users(user_id),
    ip_address INET,
    user_agent TEXT,
    status_code INT NOT NULL,
    latency_ms INT NOT NULL,
    request_size INT,
    response_size INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_logs_endpoint ON api_request_logs(endpoint, created_at DESC);
CREATE INDEX idx_api_logs_user ON api_request_logs(user_id);
CREATE INDEX idx_api_logs_created ON api_request_logs(created_at DESC);

COMMIT;
