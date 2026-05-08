-- Trade status change log for audit trail
CREATE TABLE IF NOT EXISTS trade_status_log (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    from_status VARCHAR(20) NOT NULL,
    to_status VARCHAR(20) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for trade_status_log
CREATE INDEX IF NOT EXISTS idx_trade_status_log_trade_id ON trade_status_log(trade_id);
CREATE INDEX IF NOT EXISTS idx_trade_status_log_created_at ON trade_status_log(created_at);

-- Add email service tables
CREATE TABLE IF NOT EXISTS email_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    variables TEXT[], -- Array of variable names used in template
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_queue (
    email_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    to_email VARCHAR(255) NOT NULL,
    from_email VARCHAR(255) NOT NULL DEFAULT 'noreply@cryplio.com',
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    template_id UUID REFERENCES email_templates(template_id) ON DELETE SET NULL,
    variables JSONB, -- Key-value pairs for template variables
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'cancelled')),
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    last_attempt_at TIMESTAMP,
    sent_at TIMESTAMP,
    error_message TEXT,
    priority INT NOT NULL DEFAULT 5, -- 1=high, 5=normal, 10=low
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for email_queue
CREATE INDEX IF NOT EXISTS idx_email_queue_status ON email_queue(status);
CREATE INDEX IF NOT EXISTS idx_email_queue_priority ON email_queue(priority);
CREATE INDEX IF NOT EXISTS idx_email_queue_created_at ON email_queue(created_at);
CREATE INDEX IF NOT EXISTS idx_email_queue_next_attempt ON email_queue(status, attempts, max_attempts) WHERE status = 'pending';

-- Update triggers
CREATE TRIGGER update_email_templates_updated_at BEFORE UPDATE ON email_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_email_queue_updated_at BEFORE UPDATE ON email_queue FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default email templates
INSERT INTO email_templates (name, subject, body, variables) VALUES
(
    'trade_created',
    'New Trade Created - Cryplio',
    'Hello {{user_name}},

A new trade has been created:

Trade ID: {{trade_id}}
Amount: {{crypto_amount}} {{crypto_symbol}} / {{fiat_amount}} {{fiat_symbol}}
With: {{counterpart_username}}
Payment Method: {{payment_method}}

Please complete the payment within {{payment_window}} minutes.

Best regards,
The Cryplio Team',
    ARRAY['user_name', 'trade_id', 'crypto_amount', 'crypto_symbol', 'fiat_amount', 'fiat_symbol', 'counterpart_username', 'payment_method', 'payment_window']
),
(
    'trade_completed',
    'Trade Completed - Cryplio',
    'Hello {{user_name}},

Your trade has been successfully completed:

Trade ID: {{trade_id}}
Amount: {{crypto_amount}} {{crypto_symbol}}
With: {{counterpart_username}}

Thank you for using Cryplio!

Best regards,
The Cryplio Team',
    ARRAY['user_name', 'trade_id', 'crypto_amount', 'crypto_symbol', 'counterpart_username']
),
(
    'dispute_created',
    'New Dispute Created - Cryplio',
    'Hello {{user_name}},

A dispute has been created for your trade:

Trade ID: {{trade_id}}
Reason: {{dispute_reason}}
Description: {{dispute_description}}

Our team will review this dispute and take appropriate action.

Best regards,
The Cryplio Team',
    ARRAY['user_name', 'trade_id', 'dispute_reason', 'dispute_description']
),
(
    'withdrawal_approved',
    'Withdrawal Approved - Cryplio',
    'Hello {{user_name}},

Your withdrawal request has been approved:

Amount: {{amount}} USDT
To Address: {{to_address}}
Transaction ID: {{tx_hash}}

The funds have been sent to your wallet.

Best regards,
The Cryplio Team',
    ARRAY['user_name', 'amount', 'to_address', 'tx_hash']
),
(
    'security_alert',
    'Security Alert - Cryplio',
    'Hello {{user_name}},

We detected unusual activity on your account:

{{alert_message}}

If this was not you, please secure your account immediately.

Best regards,
The Cryplio Team',
    ARRAY['user_name', 'alert_message']
) ON CONFLICT (name) DO NOTHING;
